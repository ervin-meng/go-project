package entity

import "go-project/common/domain"

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
