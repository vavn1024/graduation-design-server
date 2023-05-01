package model

import (
	"gorm.io/gorm"
)

type Picture struct {
	gorm.Model
	NoteId string `gorm:"type:varchar(20);not null"`
	Url    string `gorm:"type:varchar(255);not null"`
}
