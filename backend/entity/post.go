package entity

import "time"

type PostStatus string

const (
	PostStatusDraft    PostStatus = "draft"
	PostStatusPublish  PostStatus = "published"
	PostStatusArchived PostStatus = "archived"
)

func (PostStatus) Values() []string {
	return []string{
		string(PostStatusDraft),
		string(PostStatusPublish),
		string(PostStatusArchived),
	}
}

type Post struct {
	ID      uint
	Title   string
	Summary *string
	Cover   *string

	Author string

	Content         string
	ReadTimeMinutes uint
	ViewCount       uint

	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Tags       []PostTag
	Categories []PostCategory
}
