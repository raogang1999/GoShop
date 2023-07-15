package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// 作为公共类
type BaseModel struct {
	ID        int32          `gorm:"primary_key;type:int"` //为什么使用int32，bigint，外键
	CreatedAt time.Time      `gorm:"Column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"Column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool           `gorm:"Column:is_deleted" json:"-"`
}

type GormList []string

// 实现Valuer接口
func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 实现Scanner接口
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), g)
}
