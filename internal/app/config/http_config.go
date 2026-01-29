package config

import (
	"fmt"
	"strconv"
	"strings"

	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
)

type HTTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewHTTPConfig(host string, port int) *HTTPConfig {
	return &HTTPConfig{
		Host: host,
		Port: port,
	}
}

func DefaultHTTPConfig() *HTTPConfig {
	return NewHTTPConfig(DefaultHTTPHost, DefaultHTTPPort)
}

func (hc *HTTPConfig) GetListenerAddr() string {
	return hc.Host + ":" + strconv.Itoa(hc.Port)
}

func (hc *HTTPConfig) Copy() *HTTPConfig {
	return NewHTTPConfig(hc.Host, hc.Port)
}

func (hc *HTTPConfig) fromString(s string) error {
	params := strings.Split(s, ":")
	if len(params) < 2 {
		return errs.NewAppInvalidConfigError("http config", s, "invalid http config format, example: localhost:8080", nil)
	}

	hc.Host = params[0]
	hc.Port, _ = strconv.Atoi(params[1])

	return nil
}

// flag.Value ==================================

func (hc *HTTPConfig) String() string {
	return fmt.Sprintf("%s:%v", hc.Host, hc.Port)
}

func (hc *HTTPConfig) Set(s string) error {
	return hc.fromString(s)
}

// =============================================

// json.Marshaler ==============================

func (hc *HTTPConfig) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%s:%d", hc.Host, hc.Port)), nil
}

// json.Unmrshaler ==============================

func (hc *HTTPConfig) UnmarshalJSON(b []byte) error {
	return hc.fromString(string(b))
}

// =============================================
