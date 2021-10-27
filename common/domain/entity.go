package domain

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type EntityListType []string

func (l EntityListType) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l EntityListType) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &l)
}

type Entity struct {
	ID        int32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	isDeleted bool
}
