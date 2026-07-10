package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

// ToPostCategory converts an ent.PostCategory to entity.PostCategory.
func ToPostCategory(c *ent.PostCategory) entity.PostCategory {
	if c == nil {
		return entity.PostCategory{}
	}
	return entity.PostCategory{
		ID:        c.ID,
		Name:      c.Name,
		Slug:      c.Slug,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToPostCategories converts ent.PostCategories to entity.PostCategories.
func ToPostCategories(cats []*ent.PostCategory) []entity.PostCategory {
	if len(cats) == 0 {
		return nil
	}

	res := make([]entity.PostCategory, len(cats))
	for i, c := range cats {
		res[i] = ToPostCategory(c)
	}

	return res
}
