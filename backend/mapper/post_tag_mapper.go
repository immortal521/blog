package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

func ToPostTags(tags []*ent.PostTag) []entity.PostTag {
	if len(tags) == 0 {
		return nil
	}

	res := make([]entity.PostTag, len(tags))

	for i, t := range tags {
		res[i] = entity.PostTag{
			ID:   t.ID,
			Name: t.Name,
		}
	}

	return res
}
