// Package config
/*
  Iteration 5

  Configuration params priority :

  1 - ENV vars
  2 - CLI params
  3 - Default values
*/
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Config struct {
	AppName           string      `json:"app_name,omitempty"`
	ProjectStage      string      `json:"project_stage,omitempty"`
	LogLevel          string      `json:"log_level,omitempty"`
	BaseURL           string      `json:"base_url,omitempty" env:"BASE_URL"`
	HTTP              *HTTPConfig `json:"http,omitempty"`
	DBKind            string      `json:"db_kind,omitempty"`
	DBDsn             string      `json:"db_dsn,omitempty" env:"DATABASE_DSN"`
	StoragePath       string      `json:"storage_path,omitempty" env:"FILE_STORAGE_PATH"`
	StorageUserPath   string      `json:"storage_user_path,omitempty" env:"FILE_STORAGE_USER_PATH"`
	AuditFile         string      `json:"audit_file,omitempty" env:"AUDIT_FILE"`
	AuditIncomeLocal  bool
	AuditURL          string `json:"audit_url,omitempty" env:"AUDIT_URL"`
	AuditIncomeRemote bool
	EnableHTTPS       bool `json:"enable_https,omitempty"`
}

// Flags
const (
	FlagAppName         string = "p"
	FlagProjectStage    string = "stage"
	FlagLogLevel        string = "l"
	FlagBaseURL         string = "b"
	FlagDBKind          string = "k"
	FlagHTTPInterface   string = "a"
	FlagDBInterface     string = "d"
	FlagStoragePath     string = "f"
	FlagStorageUserPath string = "fu"
	FlagAuditFile       string = "audit-file"
	FlagAuditURL        string = "audit-url"
	FlagEnableHTTPS     string = "s"
)

// Environment variables
const (
	EnvBaseURL             string = "BASE_URL"
	EnvHTTPInterface       string = "SERVER_ADDR"
	EnvStorageFilename     string = "FILE_STORAGE_PATH"
	EnvStorageUserFilename string = "FILE_STORAGE_USER_PATH"
	EnvDatabaseDSN         string = "DATABASE_DSN"
	EnvAuditFile           string = "AUDIT_FILE"
	EnvAuditURL            string = "AUDIT_URL"
	EnvEnableHTTPS         string = "ENABLE_HTTPS"
)

func NewConfig() *Config {
	var cfg = defaultConfig()

	cfg.initFlags()

	return cfg
}

func newConfig(appName string, projectStage string, logLevel string, baseURL string, HTTP *HTTPConfig, DBKind string, DBDsn string, storagePath string) *Config {
	return &Config{
		AppName:           appName,
		ProjectStage:      projectStage,
		LogLevel:          logLevel,
		BaseURL:           baseURL,
		HTTP:              HTTP,
		DBKind:            DBKind,
		DBDsn:             DBDsn,
		StoragePath:       storagePath,
		AuditIncomeLocal:  false,
		AuditIncomeRemote: false,
		EnableHTTPS:       false,
	}
}

func defaultConfig() *Config {
	return newConfig(DefaultAppName, DefaultStage, DefaultLogLevel, DefaultBaseURL, DefaultHTTPConfig(), DefaultDBKind, DefaultDBDsn, DefaultStoragePath)
}

func (c *Config) LoadConfig() error {
	fmt.Println("Parse cli params")
	var err = c.loadCli()
	if err != nil {
		return err
	}

	fmt.Println("Parse env params")
	err = c.loadEnv()
	if err != nil {
		return err
	}

	if strings.TrimSpace(c.DBDsn) == "" {
		c.DBKind = DBKindInMemory
	}

	if !c.AuditIncomeLocal {
		c.AuditFile = DefaultAuditIncomePath
	}

	fmt.Printf("Config FINAL: [%+v]\r\n", c)

	return nil
}

func (c *Config) loadCli() error {
	flag.Parse()

	c.AuditIncomeLocal = strings.TrimSpace(c.AuditFile) != ""
	c.AuditIncomeRemote = strings.TrimSpace(c.AuditURL) != ""

	c.EnableHTTPS = c.cliFlagExists(FlagEnableHTTPS)

	fmt.Printf("Config after CLI: [%+v]\r\n", c)

	return nil
}

func (c *Config) cliFlagExists(name string) bool {
	res := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			res = true
		}
	})

	return res
}

func (c *Config) loadEnv() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}

	err = parseFlag(EnvHTTPInterface, c.HTTP)
	if err != nil {
		return err
	}

	c.AuditIncomeLocal = strings.TrimSpace(c.AuditFile) != ""
	c.AuditIncomeRemote = strings.TrimSpace(c.AuditURL) != ""

	_, c.EnableHTTPS = os.LookupEnv(EnvEnableHTTPS)

	fmt.Printf("Config after ENV: [%+v]\r\n", c)

	return nil
}

func parseFlag(env string, value flag.Value) error {
	var envVar = os.Getenv(env)
	if envVar == "" {
		return nil
	}

	fmt.Printf("[DEBUG] Config: ENV [%s] VALUE [%+v]\r\n", env, value)

	err := value.Set(envVar)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) initFlags() {
	flag.StringVar(&c.AppName, FlagAppName, DefaultAppName, "application name")
	flag.StringVar(&c.ProjectStage, FlagProjectStage, ProjectStageDevelopment, "project stage")
	flag.StringVar(&c.LogLevel, FlagLogLevel, zap.InfoLevel.CapitalString(), "log level")
	flag.StringVar(&c.BaseURL, FlagBaseURL, DefaultBaseURL, "base url")
	flag.StringVar(&c.DBKind, FlagDBKind, DefaultDBKind, "db kind")
	flag.Var(c.HTTP, FlagHTTPInterface, "http interface")
	flag.StringVar(&c.DBDsn, FlagDBInterface, DefaultDBDsn, "database dsn")
	flag.StringVar(&c.StoragePath, FlagStoragePath, DefaultStoragePath, "storage path")
	flag.StringVar(&c.StorageUserPath, FlagStorageUserPath, DefaultStorageUserPath, "storage user path")
	flag.StringVar(&c.AuditFile, FlagAuditFile, "", "audit file path")
	flag.StringVar(&c.AuditURL, FlagAuditURL, "", "audit url")
	flag.BoolVar(&c.EnableHTTPS, FlagEnableHTTPS, false, "enable https")
}
