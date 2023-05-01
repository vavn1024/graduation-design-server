package model

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	TargetId       string `gorm:"type:varchar(20);not null"`
	Type           string `gorm:"type:varchar(20);not null"`
	OriginatorName string `gorm:"type:varchar(255);not null"`
	TargetName     string `gorm:"type:varchar(255);not null"`
	Content        string `gorm:"type:varchar(255);"`
}
