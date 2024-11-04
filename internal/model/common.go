package model

import "gorm.io/gorm"

type Common struct {
	gorm.Model
}

func (m *Common) TableName() string {
    return "common"
}
