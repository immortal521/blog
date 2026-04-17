package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

func ToPostCategories(cats []*ent.PostCategory) []entity.PostCategory {
	if len(cats) == 0 {
		return nil
	}

	res := make([]entity.PostCategory, len(cats))

	for i, c := range cats {
		res[i] = entity.PostCategory{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	return res
}
