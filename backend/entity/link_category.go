package entity

import "time"

type LinkCategory struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	SortOrder int
}
