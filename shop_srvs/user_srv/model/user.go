package model

import "time"
import "gorm.io/gorm"

// 作为公共类
type BaseModel struct {
	ID        int32     `gorm:"primary_key;"`
	CreatedAt time.Time `gorm:"Column:add_time"`
	UpdatedAt time.Time `gorm:"Column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool `gorm:"Column:is_deleted"`
}

// 用户密码用md5加密；密码找回，给他一个链接；
type User struct {
	BaseModel            //继承
	Mobile    string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null "`                       //手机号
	Password  string     `gorm:"type:varchar(100);not null"`                                               //密码
	NickName  string     `gorm:"type:varchar(20)"`                                                         //昵称
	Birthday  *time.Time `gorm:"type:datetime"`                                                            //生日
	Gender    string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female 表示女 male 表示男'"` //性别
	Role      int32      `gorm:"column:role;default:1;type:int comment '1 表示普通用户 2 表示管理员'"`                //角色
}
