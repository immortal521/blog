package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

// ToPostTag converts an ent.PostTag to entity.PostTag
func ToPostTag(t *ent.PostTag) entity.PostTag {
	if t == nil {
		return entity.PostTag{}
	}
	return entity.PostTag{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// ToPostTags converts ent.PostTags to entity.PostTags.
func ToPostTags(tags []*ent.PostTag) []entity.PostTag {
	if len(tags) == 0 {
		return nil
	}

	res := make([]entity.PostTag, len(tags))

	for i, t := range tags {
		res[i] = ToPostTag(t)
	}

	return res
}
