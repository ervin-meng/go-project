package entity

import "github.com/ervin-meng/go-conch/domain"

type Inventory struct {
	domain.Entity
	GoodsId int32 `gorm:"type:int"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"`
}
