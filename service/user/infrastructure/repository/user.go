package repository

import (
	"fmt"
	"go-project/service/user/domain/entity"
	"go-project/service/user/global"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAll() []entity.User {
	var users []entity.User
	_ = global.DB_PROJECT.Find(&users)
	return users
}

func (r *UserRepository) GetById(Id int) (entity.User, error) {
	var user entity.User
	result := global.DB_PROJECT.First(&user, Id)
	if result.RowsAffected == 0 {
		return user, fmt.Errorf("not found")
	}
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UserRepository) Paginate(page, pageSize int) []entity.User {
	var users []entity.User
	global.DB_PROJECT.Scopes(Paginate(page, pageSize)).Find(&users)
	return users
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
