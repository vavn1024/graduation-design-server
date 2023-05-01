package dto

import (
	"api/model"
)

type EventDto struct {
	Content        string  `json:"content"`
	OriginatorName string  `json:"originator_name"`
	Type           string  `json:"type"`
	TargetId       string  `json:"target_id"`
	CreatedAt      int64   `json:"created_at"`
	User           UserDto `json:"user"`
}

func ToEventDto(event model.Event, user UserDto) EventDto {
	return EventDto{
		Content:        event.Content,
		OriginatorName: event.OriginatorName,
		Type:           event.Type,
		TargetId:       event.TargetId,
		CreatedAt:      event.CreatedAt.Unix(),
		User:           user,
	}
}
