package entity

import "gorm.io/gorm"

type LinkStatus int

const (
	LinkNormal LinkStatus = iota + 1
	LinkAbnormal
)

type Link struct {
	*gorm.Model

	Description string     `gorm:"column:description;size:255;comment:链接描述"`
	Enabled     bool       `gorm:"column:enabled;not null;default:false;comment:链接显示状态"`
	Name        string     `gorm:"column:name;size:100;not null;comment:链接名称"`
	SortOrder   int        `gorm:"column:sort_order;not null;default:0;comment:链接排序级别"`
	URL         string     `gorm:"column:url;size:255;unique;not null;comment:链接地址"`
	Avatar      string     `gorm:"column:avatar;size:255;comment:链接头像"`
	Status      LinkStatus `gorm:"column:status;type:smallint;not null;default:1;comment:链接状态"`

	CategoryID *uint         `gorm:"column:category_id;comment:分类ID，可以为空表示默认分类"`
	Category   *LinkCategory `gorm:"foreignkey:CategoryID"`
}

func (Link) TableName() string {
	return "links"
}
