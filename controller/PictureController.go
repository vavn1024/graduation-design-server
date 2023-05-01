package controller

import (
	"api/common"
	"api/dto"
	"api/model"
)

func saveUrls(urls []string, noteId string) {
	DB := common.DB
	for i := 0; i < len(urls); i++ {
		picture := model.Picture{
			NoteId: noteId,
			Url:    urls[i],
		}
		DB.Create(&picture)
	}
}

func getUrls(noteId string) []dto.PictureDto {
	DB := common.DB
	var pictures []model.Picture
	DB.Where("note_id = ?", noteId).Find(&pictures)
	var picturesDto []dto.PictureDto
	for i := 0; i < len(pictures); i++ {
		picturesDto = append(picturesDto, dto.ToPictureDto(pictures[i]))
	}
	return picturesDto
}
