package dto

import (
	"api/model"
)

type UserDto struct {
	Name           string `json:"name"`
	NickName       string `json:"nick_name"`
	Avatar         string `json:"avatar"`
	Banner         string `json:"banner"`
	Describe       string `json:"describe"`
	NotesCount     int    `json:"notes_count"`
	FollowingCount int    `json:"following_count"`
	FollowedCount  int    `json:"followed_count"`
	CreatedAt      int64  `json:"created_at"`
	IsFollowing    int    `json:"is_following"`
}

func ToUserDto(user model.User, options ...int) UserDto {
	var IsFollowing int
	if options != nil {
		IsFollowing = options[0]
	}
	return UserDto{
		Name:           user.Name,
		NickName:       user.Nickname,
		Avatar:         user.Avatar,
		Banner:         user.Banner,
		Describe:       user.Describe,
		NotesCount:     user.NotesCount,
		FollowingCount: user.FollowingCount,
		FollowedCount:  user.FollowedCount,
		CreatedAt:      user.CreatedAt.Unix(),
		IsFollowing:    IsFollowing,
	}
}
