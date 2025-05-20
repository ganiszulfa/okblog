package handler

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"okblog/tag/pkg/config"
	"okblog/tag/pkg/models"
)

// Global variables to store original functionality
var originalGetRedisClient func() *redis.Client

// setupRedis creates a miniredis instance for testing
func setupRedis(t *testing.T) (*miniredis.Miniredis, func()) {
	// Start a miniredis server
	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Create a redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Store original function
	originalGetRedisClient = GetRedisClient

	// Override the global function to return our test client
	GetRedisClient = func() *redis.Client {
		return client
	}

	// Return cleanup function
	return mr, func() {
		client.Close()
		mr.Close()
		GetRedisClient = originalGetRedisClient
	}
}

func TestGetPostsByTagHandler(t *testing.T) {
	// Skip if short testing mode
	// Uncomment to run only when not in short mode
	// if testing.Short() {
	//    t.Skip("skipping integration test in short mode")
	// }

	// Setup
	mr, cleanup := setupRedis(t)
	defer cleanup()

	// Create fiber app
	app := fiber.New()
	app.Get("/api/tag/:tagName", getPostsByTagHandler)

	// Setup test data
	tagName := "golang"
	setName := config.TagPostsPrefix + tagName
	postID1 := "post-id-1"
	postID2 := "post-id-2"
	detailsKey1 := config.PostDetailsPrefix + postID1
	detailsKey2 := config.PostDetailsPrefix + postID2

	now := time.Now()
	postDetails1 := models.CachedPostDetails{
		Title:       "Test Post 1",
		PublishedAt: now.Add(-24 * time.Hour),
		Tags:        []string{tagName},
		Slug:        "test-post-1",
		ViewCount:   100,
	}
	postDetails2 := models.CachedPostDetails{
		Title:       "Test Post 2",
		PublishedAt: now.Add(-48 * time.Hour),
		Tags:        []string{tagName},
		Slug:        "test-post-2",
		ViewCount:   200,
	}

	// Populate Redis with test data
	postJSON1, err := json.Marshal(postDetails1)
	require.NoError(t, err)
	postJSON2, err := json.Marshal(postDetails2)
	require.NoError(t, err)

	// First, add the posts to the sorted set
	mr.ZAdd(setName, float64(now.Unix()), postID1)
	mr.ZAdd(setName, float64(now.Add(-24*time.Hour).Unix()), postID2)

	// Then, set the post details
	mr.Set(detailsKey1, string(postJSON1))
	mr.Set(detailsKey2, string(postJSON2))

	// Make request
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/tag/%s", tagName), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Assertions
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result models.PaginatedPostsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, 1, result.Pagination.CurrentPage)
	assert.Equal(t, 10, result.Pagination.PerPage)
	assert.Equal(t, 1, result.Pagination.TotalPages)
	assert.Equal(t, 2, result.Pagination.TotalItems)

	// Verify post data
	assert.Equal(t, "Test Post 1", result.Data[0].Title)
	assert.Equal(t, "test-post-1", result.Data[0].Slug)
	assert.Equal(t, 100, result.Data[0].ViewCount)

	assert.Equal(t, "Test Post 2", result.Data[1].Title)
	assert.Equal(t, "test-post-2", result.Data[1].Slug)
	assert.Equal(t, 200, result.Data[1].ViewCount)
}

func TestGetPostsByTagHandler_NoResults(t *testing.T) {
	// Setup
	_, cleanup := setupRedis(t)
	defer cleanup()

	// Create fiber app
	app := fiber.New()
	app.Get("/api/tag/:tagName", getPostsByTagHandler)

	// Define test data - tag with no posts
	tagName := "nonexistenttag"

	// Make request
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/tag/%s", tagName), nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Assertions
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result models.PaginatedPostsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Verify response structure for empty results
	assert.Equal(t, 0, len(result.Data))
	assert.Equal(t, 1, result.Pagination.CurrentPage)
	assert.Equal(t, 10, result.Pagination.PerPage)
	assert.Equal(t, 0, result.Pagination.TotalPages)
	assert.Equal(t, 0, result.Pagination.TotalItems)
}
