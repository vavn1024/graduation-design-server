package dto

import (
	"api/model"
)

type NoteDto struct {
	NoteId    string       `json:"note_id"`
	Content   string       `json:"content"`
	Name      string       `json:"name"`
	CreatedAt int64        `json:"created_at"`
	User      UserDto      `json:"user"`
	Pictures  []PictureDto `json:"pictures"`
	IsLiked   int          `json:"isLiked"`
}

func ToNoteDto(note model.Note, user UserDto, pictures []PictureDto, isLiked int) NoteDto {
	return NoteDto{
		NoteId:    note.NoteId,
		Content:   note.Content,
		Name:      note.Name,
		CreatedAt: note.CreatedAt.Unix(),
		User:      user,
		Pictures:  pictures,
		IsLiked:   isLiked,
	}
}

func ToNoteDto2(note model.Note, user UserDto) NoteDto {
	return NoteDto{
		NoteId:    note.NoteId,
		Content:   note.Content,
		Name:      note.Name,
		CreatedAt: note.CreatedAt.Unix(),
		User:      user,
	}
}
