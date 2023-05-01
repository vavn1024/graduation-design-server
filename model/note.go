package model

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	NoteId  string `gorm:"type:varchar(20);not null;unique"`
	Name    string `gorm:"type:varchar(20);not null"`
	Content string `gorm:"size:3000;not null"`
}
