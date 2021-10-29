package entity

import "github.com/ervin-meng/go-stitch-monster/domain"

type User struct {
	domain.Entity
	Mobile   string `gorm:"index:id_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	Gender   string `gorm:"default:male;type:varchar(6) comment 'female 女 male 男'"`
	Role     int    `gorm:"default:1;type:int comment '1 普通用户 2 管理员'"`
}
