package model

import (
	"database/sql/driver"
	"fmt"
	"net/url"
)

// CorrelationUrls - URL для преобразования, где:
//
// key - correlationID
//
// value - URL
type CorrelationUrls map[string]string

// CorrelationShorts - URL после преобразования, где:
//
// key - correlationID
//
// value - short URL
type CorrelationShorts map[string]string

// UserShorts - список пользовательских сокращённых URL, где:
//
// key - original URL
//
// value - short URL
type UserShorts map[string]string

// UserBatchDeletes - список keys для удаления short URL
type UserBatchDeletes []string

// CustomURL - URL для серияализации
type CustomURL struct {
	URL *url.URL
}

// Valuer

func (cu *CustomURL) Value() (driver.Value, error) {
	if cu.URL == nil || cu.URL.String() == "" {
		return nil, nil
	}

	return cu.URL.String(), nil
}

// ==============

// Scanner

func (cu *CustomURL) Scan(src any) error {
	if src == nil {
		cu.URL = nil

		return nil
	}

	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("cannot scan custom url from %T", src)
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return err
	}

	cu.URL = parsed

	return nil
}

// ==================
