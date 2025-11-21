package model

import "time"

type TechData struct {
	CreateUser string    `db:"create_user" json:"create_user,omitempty"`
	Created    time.Time `db:"created" json:"created,omitempty"`
	UpdateUser string    `db:"update_user" json:"update_user,omitempty"`
	Updated    time.Time `db:"updated" json:"updated,omitempty"`
}
