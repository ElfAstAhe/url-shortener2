package model

import (
	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/google/uuid"
)

type ShortURIUser struct {
	ID         string `db:"id" json:"id"`
	ShortURIID string `db:"short_uri_id" json:"short_uri_id"`
	UserID     string `db:"user_id" json:"user_id"`
	Deleted    bool   `db:"deleted" json:"deleted"`
}

func NewShortURIUser(shortURIID string, userID string) (*ShortURIUser, error) {
	newID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &ShortURIUser{
		ID:         newID.String(),
		ShortURIID: shortURIID,
		UserID:     userID,
	}, nil
}

func ValidateShortURIUser(entity *ShortURIUser) error {
	if entity == nil {
		return errs.NewAppInvalidArgumentError("entity", "empty")
	}
	if entity.ShortURIID == "" {
		return errs.NewAppInvalidArgumentError("entity.ShortURIID", "empty")
	}
	if entity.UserID == "" {
		return errs.NewAppInvalidArgumentError("entity.UserID", "empty")
	}

	return nil
}
