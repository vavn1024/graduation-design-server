package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"type:varchar(20);not null;unique"`
	Nickname       string `gorm:"type:varchar(20);not null"`
	Password       string `gorm:"size:255;not null"`
	Avatar         string `gorm:"size:255;not null;default:'https://misskey.io/identicon/98ctey7bpu'"`
	Banner         string `gorm:"size:255;not null;default:'http://127.0.0.1:1012/avatar/1676785688.jpg'"`
	Describe       string `gorm:"size:255;not null"`
	FollowingCount int    `gorm:"type:int(11);not null"`
	FollowedCount  int    `gorm:"type:int(11);not null"`
	NotesCount     int    `gorm:"type:int(11);not null"`
}
