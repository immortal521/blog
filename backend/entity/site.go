package entity

import "time"

type Site struct {
	Name        string
	Logo        *string
	Greeting    *string
	Description *string
	CreatedAt   *time.Time
}
