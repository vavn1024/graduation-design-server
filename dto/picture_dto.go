package dto

import (
	"api/model"
)

type PictureDto struct {
	Url string `json:"url"`
}

func ToPictureDto(picture model.Picture) PictureDto {
	return PictureDto{
		Url: picture.Url,
	}
}
