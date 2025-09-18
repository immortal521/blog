package dto

import "time"

type PostListRes struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Cover           *string    `json:"cover"`
	Summary         *string    `json:"summary"`
	ReadTimeMinutes uint       `json:"readTimeMinutes"`
	ViewCount       uint       `json:"viewCount"`
	PublishedAt     *time.Time `json:"publishedAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	Author          string     `json:"author"`
}

type PostMetaRes struct {
	ID        uint      `json:"id"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostRes struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Summary         *string    `json:"summary"`
	Content         string     `json:"content"`
	Cover           *string    `json:"cover"`
	ReadTimeMinutes uint       `json:"readTimeMinutes"`
	ViewCount       uint       `json:"viewCount"`
	PublishedAt     *time.Time `json:"publishedAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	Author          string     `json:"author"`
}
