package entity

import "github.com/ervin-meng/go-stitch-monster/domain"

type CategoryBrand struct {
	domain.Entity
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category
	BrandID    int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brand      Brand
}

func (CategoryBrand) TableName() string {
	return "category_brand"
}
