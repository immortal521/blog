package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

func ToLink(l *ent.Link) *entity.Link {
	if l == nil {
		return &entity.Link{}
	}
	return &entity.Link{
		ID:          l.ID,
		Name:        l.Name,
		URL:         l.URL,
		Enabled:     l.Enabled,
		Description: &l.Description,
		Avatar:      &l.Avatar,
		Status:      l.Status,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func ToLinks(ls []*ent.Link) []*entity.Link {
	if len(ls) == 0 {
		return nil
	}

	res := make([]*entity.Link, len(ls))
	for i, l := range ls {
		res[i] = ToLink(l)
	}

	return res
}
