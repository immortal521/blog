package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

func extractAuthor(p *ent.Post) string {
	if p.Edges.Author == nil {
		return ""
	}
	return p.Edges.Author.Username
}

// ToPost converts an ent.Post to entity.Post.
func ToPost(p *ent.Post) *entity.Post {
	if p == nil {
		return &entity.Post{}
	}

	return &entity.Post{
		ID:      p.ID,
		Title:   p.Title,
		Summary: p.Summary,
		Cover:   p.Cover,

		User:   extractAuthor(p),
		UserID: p.UserID,

		Content:         p.Content,
		ReadTimeMinutes: p.ReadTimeMinutes,
		ViewCount:       p.ViewCount,

		Status: p.Status,

		PublishedAt: p.PublishedAt,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,

		Tags:       ToPostTags(p.Edges.Tags),
		Categories: ToPostCategories(p.Edges.Categories),
	}
}

func ToPosts(ps []*ent.Post) []*entity.Post {
	if len(ps) == 0 {
		return nil
	}

	res := make([]*entity.Post, len(ps))
	for i, p := range ps {
		res[i] = ToPost(p)
	}

	return res
}
