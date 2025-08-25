package dto

import "time"

type PostShortResponseDTO struct {
	ID    uint    `json:"id"`
	Title string  `json:"title"`
	Cover *string `json:"cover"`
}

type PostResponseDTO struct {
	ID              uint       `json:"id"`
	Title           string     `json:"title"`
	Summary         *string    `json:"summary"`
	Content         string     `json:"content"`
	Cover           *string    `json:"cover"`
	ReadTimeMinutes int64      `json:"readTimeMinutes"`
	ViewCount       int64      `json:"viewCount"`
	PublishedAt     *time.Time `json:"publishedAt"`
}
