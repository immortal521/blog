package entity

import "time"

// LinkStatus represents the status of a link
type LinkStatus int

const (
	// LinkNormal indicates the link is accessible and working
	LinkNormal LinkStatus = iota + 1
	// LinkAbnormal indicates the link is broken or inaccessible
	LinkAbnormal
)

// Link represents a friend link entity
type Link struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`

	Description string     `gorm:"column:description;size:255"`                    // Link description
	Enabled     bool       `gorm:"column:enabled;not null;default:false"`          // Whether the link is enabled
	Name        string     `gorm:"column:name;size:100;not null"`                  // Link name/title
	SortOrder   int        `gorm:"column:sort_order;not null;default:0"`           // Display order
	URL         string     `gorm:"column:url;size:255;unique;not null"`            // Link URL (unique)
	Avatar      string     `gorm:"column:avatar;size:255"`                         // Avatar image URL
	Status      LinkStatus `gorm:"column:status;type:smallint;not null;default:1"` // Link status

	CategoryID *uint         `gorm:"column:category_id"`    // Foreign key to LinkCategory (optional)
	Category   *LinkCategory `gorm:"foreignkey:CategoryID"` // Associated category
}

// TableName returns the table name for Link model
func (Link) TableName() string {
	return "links"
}
