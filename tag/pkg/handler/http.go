package handler

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"okblog/tag/pkg/config"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/models"
)

// for consistent storage and retrieval
func normalizeTagName(tagName string) string {
	decoded, err := url.QueryUnescape(tagName)
	if err != nil {
		decoded = tagName
	}
	return strings.ToLower(strings.TrimSpace(decoded))
}

// InitFiberApp initializes the Fiber application and sets up routes.
func InitFiberApp() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("{\"status\":\"ok\"}")
	})

	api := app.Group("/api")
	api.Get("/tag/:tagName", getPostsByTagHandler)
	api.Post("/tag/collect", collectTagsHandler)

	return app
}

// Initialize GetRedisClient to use the database client by default
var GetRedisClient = database.GetClient

// getPostsByTagHandler handles requests for posts by a specific tag with pagination.
func getPostsByTagHandler(c *fiber.Ctx) error {
	tagName := normalizeTagName(c.Params("tagName"))
	if tagName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Tag name cannot be empty"})
	}

	pageQuery := c.Query("page", strconv.Itoa(config.DefaultPage))
	perPageQuery := c.Query("per_page", strconv.Itoa(config.DefaultPerPage))

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		page = config.DefaultPage
	}
	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil || perPage < 1 {
		perPage = config.DefaultPerPage
	}

	valkeyClient := GetRedisClient()
	setName := config.TagPostsPrefix + tagName

	// Calculate start and stop for ZREVRANGE (0-indexed)
	start := int64((page - 1) * perPage)
	stop := start + int64(perPage) - 1

	ctx := context.Background()

	// Get total items in the sorted set for pagination
	totalItems, err := valkeyClient.ZCard(ctx, setName).Result()
	if err != nil {
		log.Printf("Error getting ZCard for %s: %v", setName, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve tag count"})
	}

	log.Printf("Total items in %s: %d", setName, totalItems)

	if totalItems == 0 {
		return c.Status(fiber.StatusOK).JSON(models.PaginatedPostsResponse{
			Data: []models.APIPostResponse{},
			Pagination: models.PaginationDetails{
				CurrentPage: page,
				PerPage:     perPage,
				TotalPages:  0,
				TotalItems:  0,
			},
		})
	}

	// Fetch post IDs from the sorted set (scores are publishedAt timestamps, so ZREVRANGE is correct for DESC)
	postIDs, err := valkeyClient.ZRevRange(ctx, setName, start, stop).Result()
	if err != nil {
		log.Printf("Error ZRevRange for %s: %v", setName, err)
		// If the error is because the key doesn't exist, it's like having 0 items.
		// redis.Nil is typically returned in such cases.
		if err == redis.Nil {
			return c.Status(fiber.StatusOK).JSON(models.PaginatedPostsResponse{
				Data:       []models.APIPostResponse{},
				Pagination: models.PaginationDetails{CurrentPage: page, PerPage: perPage, TotalPages: 0, TotalItems: 0},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve posts for tag"})
	}

	if len(postIDs) == 0 {
		return c.Status(fiber.StatusOK).JSON(models.PaginatedPostsResponse{
			Data: []models.APIPostResponse{},
			Pagination: models.PaginationDetails{
				CurrentPage: page,
				PerPage:     perPage,
				TotalPages:  int(math.Ceil(float64(totalItems) / float64(perPage))),
				TotalItems:  int(totalItems),
			},
		})
	}

	// Fetch full post details for each ID
	posts := make([]models.APIPostResponse, 0, len(postIDs))
	for _, postID := range postIDs {
		detailsKey := config.PostDetailsPrefix + postID
		postJSON, err := valkeyClient.Get(ctx, detailsKey).Result()
		if err == redis.Nil {
			log.Printf("Post details not found in Valkey for ID %s (key: %s)", postID, detailsKey)
			continue // Skip this post if details are missing
		}
		if err != nil {
			log.Printf("Error getting post details for ID %s from Valkey: %v", postID, err)
			// Decide if we should fail the request or skip this post
			continue
		}

		var cachedData models.CachedPostDetails
		if err := json.Unmarshal([]byte(postJSON), &cachedData); err != nil {
			log.Printf("Error unmarshalling cached post details for ID %s: %v. JSON: %s", postID, err, postJSON)
			continue
		}
		posts = append(posts, models.ToAPIPostResponseFromCached(postID, cachedData))
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	var nextPage *int
	if page < totalPages {
		p := page + 1
		nextPage = &p
	}
	var prevPage *int
	if page > 1 {
		p := page - 1
		prevPage = &p
	}

	return c.Status(fiber.StatusOK).JSON(models.PaginatedPostsResponse{
		Data: posts,
		Pagination: models.PaginationDetails{
			CurrentPage: page,
			PerPage:     perPage,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			NextPage:    nextPage,
			PrevPage:    prevPage,
		},
	})
}

// collectTagsHandler collects tags from published posts and creates sorted set caches
func collectTagsHandler(c *fiber.Ctx) error {
	ctx := context.Background()

	// Connect to MySQL database
	db := database.GetMySQLDB()

	// Get Valkey client
	valkeyClient := database.GetClient()

	// Query to get all distinct tags from post_tags table joined with published posts
	rows, err := db.Query("SELECT DISTINCT tag FROM post_tags pt JOIN posts p ON pt.post_id = p.id WHERE p.is_published = TRUE AND pt.tag IS NOT NULL AND pt.tag != ''")
	if err != nil {
		log.Printf("Error querying tags from database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query tags from database",
		})
	}
	defer rows.Close()

	// Collect all tags - map normalized tag name to original tag name
	allTags := make(map[string]string) // normalized -> original
	for rows.Next() {
		var originalTag string
		if err := rows.Scan(&originalTag); err != nil {
			log.Printf("Error scanning tag row: %v", err)
			continue
		}

		normalizedTag := normalizeTagName(originalTag)
		if normalizedTag != "" {
			allTags[normalizedTag] = originalTag
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating tag rows: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error processing tags",
		})
	}

	// Check which tags need cache creation (up to 10)
	tagsToProcess := make([]string, 0, 10)
	for normalizedTag := range allTags {
		// Check if sorted set exists
		setName := config.TagPostsPrefix + normalizedTag
		exists, err := valkeyClient.Exists(ctx, setName).Result()
		if err != nil {
			log.Printf("Error checking if sorted set exists for tag %s: %v", normalizedTag, err)
			continue
		}

		if exists == 0 {
			tagsToProcess = append(tagsToProcess, normalizedTag)
			// Break if we have 10 tags to process
			if len(tagsToProcess) >= 10 {
				break
			}
		}
	}

	if len(tagsToProcess) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":        "No new tags to process",
			"processed_tags": []string{},
		})
	}

	// Create a response map for processed tags and their counts
	responseData := make(map[string]int)

	// Process each tag that needs a cache
	for _, normalizedTag := range tagsToProcess {
		// Get the original tag name for database query
		originalTag := allTags[normalizedTag]

		// Query posts for this tag using a JOIN between posts and post_tags
		query := "SELECT BIN_TO_UUID(p.id) as id, p.title, p.published_at, p.slug, p.view_count " +
			"FROM posts p " +
			"JOIN post_tags pt ON p.id = pt.post_id " +
			"WHERE p.is_published = TRUE AND pt.tag = ? " +
			"ORDER BY p.published_at DESC LIMIT 100"

		postRows, err := db.Query(query, originalTag)
		if err != nil {
			log.Printf("Error querying posts for tag %s (original: %s): %v", normalizedTag, originalTag, err)
			continue
		}

		setName := config.TagPostsPrefix + normalizedTag
		postsAdded := 0

		// Process posts for this tag
		for postRows.Next() {
			var (
				id          string
				title       string
				publishedAt time.Time
				slug        string
				viewCount   int
			)

			if err := postRows.Scan(&id, &title, &publishedAt, &slug, &viewCount); err != nil {
				log.Printf("Error scanning post row: %v", err)
				continue
			}

			// Create CachedPostDetails object
			postDetails := models.CachedPostDetails{
				Title:       title,
				PublishedAt: publishedAt,
				// Tags:        []string{normalizedTag}, // TODO: Add tags to the cache
				Slug:      slug,
				ViewCount: viewCount,
			}

			// Convert to JSON
			postJSON, err := json.Marshal(postDetails)
			if err != nil {
				log.Printf("Error marshalling post details: %v", err)
				continue
			}

			// Store post details in cache
			detailsKey := config.PostDetailsPrefix + id
			if err := valkeyClient.Set(ctx, detailsKey, string(postJSON), 0).Err(); err != nil {
				log.Printf("Error setting post details in cache: %v", err)
				continue
			}

			// Add to sorted set with published_at as score
			score := float64(publishedAt.Unix())
			if err := valkeyClient.ZAdd(ctx, setName, redis.Z{Score: score, Member: id}).Err(); err != nil {
				log.Printf("Error adding post to sorted set: %v", err)
				continue
			}

			postsAdded++
		}

		postRows.Close()

		// Record the number of posts added for this tag
		responseData[normalizedTag] = postsAdded
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":        "Tags processed successfully",
		"processed_tags": responseData,
	})
}

// StartHTTPServer starts the HTTP server using the provided Fiber app.
func StartHTTPServer(app *fiber.App) {
	fiberPort := os.Getenv("FIBER_PORT")
	if fiberPort == "" {
		fiberPort = "3001" // Default port for tag service
	}
	log.Printf("Starting Fiber server on port %s", fiberPort)
	if err := app.Listen(":" + fiberPort); err != nil {
		log.Fatalf("Fiber app Listen error: %v", err)
	}
}
