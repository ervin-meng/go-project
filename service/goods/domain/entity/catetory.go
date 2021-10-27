package entity

import "go-project/common/domain"

type Category struct {
	domain.Entity
	ParentID int32
	Parent   *Category
	Name     string `gorm:"type:varchar(20);not null"`
	Level    int32  `gorm:"type:int;not null;default:1"`
	IsTab    bool   `gorm:"default:false;not null"`
}
