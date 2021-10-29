package entity

import "github.com/ervin-meng/go-stitch-monster/domain"

type Goods struct {
	domain.Entity
	No           string `gorm:"type:varchar(50);not null"`
	Name         string `gorm:"type:varchar(50);not null"`
	CategoryId   int32  `gorm:"type:int;not null"`
	Category     Category
	BrandID      int32 `gorm:"typse:int;not null"`
	Brand        Brand
	IsShow       bool                  `gorm:"default:false;not null"`
	IsNew        bool                  `gorm:"default:false;not null"`
	IsHot        bool                  `gorm:"default:false;not null"`
	MarketPrice  float32               `gorm:"not null"`
	Price        float32               `gorm:"not null"`
	Brief        string                `gorm:"type:varchar(100);not null comment '商品描述'"`
	Images       domain.EntityListType `gorm:"type:varchar(1000);not null"`
	DetailImages domain.EntityListType `gorm:"type:varchar(1000);not null"`
	CoverImage   string                `gorm:"type:varchar(200);not null comment '商品封面图'"`
}
