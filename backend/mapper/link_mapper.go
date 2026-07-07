package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

// ToLink converts an ent.Link to entity.Link.
func ToLink(l *ent.Link) *entity.Link {
	if l == nil {
		return &entity.Link{}
	}
	link := &entity.Link{
		ID:          l.ID,
		Name:        l.Name,
		URL:         l.URL,
		Enabled:     l.Enabled,
		Description: &l.Description,
		Avatar:      &l.Avatar,
		Status:      l.Status,
		SortOrder:   l.SortOrder,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}

	// CategoryID is 0 when not set in the database (ent uses plain uint).
	// Convert to *uint only when non-zero to preserve nil semantics.
	if l.CategoryID != 0 {
		id := l.CategoryID
		link.CategoryID = &id
	}

	return link
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
