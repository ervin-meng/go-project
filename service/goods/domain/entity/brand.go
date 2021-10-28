package entity

import "github.com/ervin-meng/go-conch/domain"

type Brand struct {
	domain.Entity
	Name string `gorm:"type:varchar(20);default:'';not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`
}
