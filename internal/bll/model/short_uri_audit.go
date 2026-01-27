package model

import "time"

const (
	OperationCreate string = "create"
	OperationUpdate string = "update"
	OperationDelete string = "delete"
)

type ShortURIAudit struct {
	ID         string
	ShortURIID string
	UserID     string
	User       string
	Date       time.Time
	Operation  string
}
