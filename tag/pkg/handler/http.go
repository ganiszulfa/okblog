package handler

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"okblog/tag/pkg/config"
	"okblog/tag/pkg/database"
	"okblog/tag/pkg/models"
)

// InitFiberApp initializes the Fiber application and sets up routes.
func InitFiberApp() *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("{\"status\":\"ok\"}")
	})

	api := app.Group("/api")
	api.Get("/tag/:tagName", getPostsByTagHandler)

	return app
}

// getPostsByTagHandler handles requests for posts by a specific tag with pagination.
func getPostsByTagHandler(c *fiber.Ctx) error {
	tagName := strings.ToLower(strings.TrimSpace(c.Params("tagName")))
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

	valkeyClient := database.GetClient()
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
