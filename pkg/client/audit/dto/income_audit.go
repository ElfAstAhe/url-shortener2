package dto

import "time"

type IncomeAuditDto struct {
	TS     int64  `json:"ts"`                // unix timestamp события
	Action string `json:"action"`            // действие: shorten (создание) или follow (прохождение по ссылке)
	UserID string `json:"user_id,omitempty"` // идентификатор пользователя, если есть
	URL    string `json:"url"`               // оригинальный (не сокращенный) URL
}

func NewIncomeAuditDto(date time.Time, action string, userID string, url string) *IncomeAuditDto {
	return &IncomeAuditDto{
		TS:     date.Unix(),
		Action: action,
		UserID: userID,
		URL:    url,
	}
}
