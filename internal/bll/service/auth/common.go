package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

const (
	GuestAdmin  bool   = false
	GuestUserID string = "guest_user_ID"
	GuestUser   string = "guest"
)
const (
	UnknownAdmin  bool   = false
	UnknownUserID string = "unknown_user_ID"
	UnknownUser   string = "unknown"
)
const CookieName string = "Authorization"

type ContextUserInfoType string

const ContextUserInfo ContextUserInfoType = "UserInfo"

type Roles []string

type UserInfo struct {
	Admin  bool   `json:"admin"`
	UserID string `json:"user_id,omitempty"`
	User   string `json:"user,omitempty"`
	Roles  Roles  `json:"roles,omitempty"`
}

func NewUserInfo(admin bool, userID string, user string, roles Roles) *UserInfo {
	return &UserInfo{
		Admin:  admin,
		UserID: userID,
		User:   user,
		Roles:  roles,
	}
}

func UserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	if res, ok := ctx.Value(ContextUserInfo).(*UserInfo); ok {
		return res, nil
	}

	return nil, errs.NewAuthInfoAbsentError("user info not found in context", nil)
}

func HasUserInfoInContext(ctx context.Context) bool {
	userInfo, err := UserInfoFromContext(ctx)
	if err != nil {
		return false
	}

	return userInfo != nil
}

func HasUserInfoInRequest(r *http.Request) bool {
	if r == nil {
		return false
	}

	return HasUserInfoInContext(r.Context())
}

func BuildRandomUserInfo() *UserInfo {
	randUUID, err := uuid.NewRandom()
	if err != nil {
		return BuildUnknownUserInfo()
	}

	return NewUserInfo(false, fmt.Sprintf("user_id-%s", randUUID.String()), fmt.Sprintf("user-%s", randUUID.String()), nil)
}

func BuildGuestUserInfo() *UserInfo {
	return NewUserInfo(GuestAdmin, GuestUserID, GuestUser, nil)
}

func BuildUnknownUserInfo() *UserInfo {
	return NewUserInfo(UnknownAdmin, UnknownUserID, UnknownUser, nil)
}
