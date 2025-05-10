package models

import (
	"encoding/json"
	"time"
)

// DebeziumMessage matches the overall structure of the Kafka message from Debezium.
type DebeziumMessage struct {
	Schema  interface{} `json:"schema"` // Keeping it flexible
	Payload Payload     `json:"payload"`
}

// Payload matches the 'payload' field in the Debezium message.
type Payload struct {
	Before *PostPayloadData       `json:"before"`
	After  *PostPayloadData       `json:"after"`
	Source map[string]interface{} `json:"source"` // Or a more specific struct if needed
	Op     string                 `json:"op"`
	TsMs   int64                  `json:"ts_ms"`
	// Transaction can be added if needed: Transaction *struct { ... } `json:"transaction"`
}

// PostPayloadData matches the structure of the 'after' or 'before' field in the Debezium payload.
// JSON tags match the fields from the Debezium message.
type PostPayloadData struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Title       string `json:"title"`
	PublishedAt *int64 `json:"published_at"` // pointer for nullable, epoch microseconds
	IsPublished bool   `json:"is_published"`
	Slug        string `json:"slug"`
	ViewCount   int32  `json:"view_count"`
}

// Post represents the structure of a message from Kafka.
// Fields are mapped from com.okblog.post.model.Post but not all fields are included.
type Post struct {
	ID          string    `json:"id"`
	Type        string    `json:"type,omitempty"`
	Title       string    `json:"title,omitempty"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
	IsPublished bool      `json:"isPublished"`
	Slug        string    `json:"slug,omitempty"`
	ViewCount   int       `json:"viewCount,omitempty"`
}

// CachedPostDetails is a slimmed-down version of a Post for caching.
type CachedPostDetails struct {
	Title       string    `json:"title,omitempty"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
	Slug        string    `json:"slug,omitempty"`
	ViewCount   int       `json:"viewCount,omitempty"`
}

// APIPostResponse is used for the list items in the API, omits content.
// It also renames IsPublished to published for the API output.
type APIPostResponse struct {
	Title       string    `json:"title,omitempty"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
	Slug        string    `json:"slug,omitempty"`
	ViewCount   int       `json:"viewCount,omitempty"`
}

// PaginationDetails holds metadata for paginated responses.
type PaginationDetails struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalPages  int  `json:"total_pages"`
	TotalItems  int  `json:"total_items"`
	NextPage    *int `json:"next_page"` // Pointer to allow null
	PrevPage    *int `json:"prev_page"` // Pointer to allow null
}

// PaginatedPostsResponse is the structure for the GET /api/tag/:tagName endpoint.
type PaginatedPostsResponse struct {
	Data       []APIPostResponse `json:"data"`
	Pagination PaginationDetails `json:"pagination"`
}

// ToAPIPostResponse converts Post to APIPostResponse
func ToAPIPostResponse(p Post) APIPostResponse {
	return APIPostResponse{
		Title:       p.Title,
		PublishedAt: p.PublishedAt,
		Tags:        p.Tags,
		Slug:        p.Slug,
		ViewCount:   p.ViewCount,
	}
}

// ToAPIPostResponseFromCached converts a CachedPostDetails object and its ID to an APIPostResponse.
func ToAPIPostResponseFromCached(id string, cached CachedPostDetails) APIPostResponse {
	return APIPostResponse{
		Title:       cached.Title,
		PublishedAt: cached.PublishedAt,
		Tags:        cached.Tags,
		Slug:        cached.Slug,
		ViewCount:   cached.ViewCount,
	}
}

// Post tags models

// PostTagData represents a relationship between a post and a tag
type PostTagData struct {
	PostID string `json:"post_id"`
	Tag    string `json:"tag"`
}

// PostTagPayload matches the payload field in Debezium messages for post tags
type PostTagPayload struct {
	Before *PostTagData `json:"before"`
	After  *PostTagData `json:"after"`
	Op     string       `json:"op"`
	TsMs   int64        `json:"ts_ms"`
}

// PostTagDebeziumMessage is the full Debezium message for post tags
type PostTagDebeziumMessage struct {
	Schema  json.RawMessage `json:"schema,omitempty"`
	Payload PostTagPayload  `json:"payload"`
}
