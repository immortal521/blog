package entity

import (
	"time"
)

type Comment struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
