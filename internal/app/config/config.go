// Package config
/*
  Iteration 5

  Configuration params priority :

  1 - ENV vars
  2 - CLI params
  3 - config file
  4 - Default values
*/
package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	errs "github.com/ElfAstAhe/url-shortener2/pkg/error"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type AppConf struct {
	AppName           string      `json:"app_name,omitempty"`
	ProjectStage      string      `json:"project_stage,omitempty"`
	LogLevel          string      `json:"log_level,omitempty"`
	BaseURL           string      `json:"base_url,omitempty" env:"BASE_URL"`
	HTTP              *HTTPConfig `json:"server_address,omitempty"`
	DBKind            string      `json:"database_kind,omitempty"`
	DBDsn             string      `json:"database_dsn,omitempty" env:"DATABASE_DSN"`
	StoragePath       string      `json:"file_storage_path,omitempty" env:"FILE_STORAGE_PATH"`
	StorageUserPath   string      `json:"storage_user_path,omitempty" env:"FILE_STORAGE_USER_PATH"`
	AuditFile         string      `json:"audit_file,omitempty" env:"AUDIT_FILE"`
	AuditIncomeLocal  bool
	AuditURL          string `json:"audit_url,omitempty" env:"AUDIT_URL"`
	AuditIncomeRemote bool
	EnableHTTPS       bool   `json:"enable_https,omitempty"`
	configPath        string `env:"CONFIG"`
}

func newConfig(appName string, projectStage string, logLevel string, baseURL string, HTTP *HTTPConfig, DBKind string, DBDsn string, storagePath string) *AppConf {
	return &AppConf{
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

func DefaultConfig() *AppConf {
	return newConfig(DefaultAppName, DefaultStage, DefaultLogLevel, DefaultBaseURL, DefaultHTTPConfig(), DefaultDBKind, DefaultDBDsn, DefaultStoragePath)
}

func GetAppConfig() (*AppConf, error) {
	fmt.Println("Parse cli params")
	cliConf := intermediateConfig.loadCli()

	fmt.Println("Parse env params")
	envConf, err := intermediateConfig.loadEnv()
	if err != nil {
		return err
	}

	fmt.Printf("Parse conf file params")
	err = intermediateConfig.loadConf()
	if err != nil {
		return err
	}

	if strings.TrimSpace(c.DBDsn) == "" {
		c.DBKind = DBKindInMemory
	}

	if !c.AuditIncomeLocal {
		c.AuditFile = DefaultAuditIncomePath
	}

	fmt.Printf("AppConf FINAL: [%+v]\r\n", c)

	return nil
}

func (c *AppConf) loadConf() error {
	// ToDo: implement

	return nil
}

func (c *AppConf) loadCli() (*AppConf, err error) {
	flag.Parse()
	defer func() {
		if r := recover(); r != nil {
			err = errs.NewAppGeneralInvalidConfigError(fmt.Sprintf("load cli params, panic [%v]", r), nil)
		}
	}()

	if strings.TrimSpace(c.DBDsn) == "" {
		c.DBKind = DBKindInMemory
	}

	if !c.AuditIncomeLocal {
		c.AuditFile = DefaultAuditIncomePath
	}

	c.AuditIncomeLocal = strings.TrimSpace(c.AuditFile) != ""
	c.AuditIncomeRemote = strings.TrimSpace(c.AuditURL) != ""

	if c.cliFlagExists(FlagEnableHTTPS) {
		c.EnableHTTPS = true
	}

	fmt.Printf("AppConf after CLI: [%+v]\r\n", c)

	return c.copy(), err
}

func (c *AppConf) cliFlagExists(name string) bool {
	res := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			res = true
		}
	})

	return res
}

func (c *AppConf) envVarExists(name string) bool {
	_, exists := os.LookupEnv(name)

	return exists
}

func (c *AppConf) loadEnv() (*AppConf, error) {
	err := env.Parse(c)
	if err != nil {
		return nil, err
	}

	err = parseFlag(EnvHTTPInterface, c.HTTP)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(c.DBDsn) == "" {
		c.DBKind = DBKindInMemory
	}

	if !c.AuditIncomeLocal {
		c.AuditFile = DefaultAuditIncomePath
	}

	c.AuditIncomeLocal = strings.TrimSpace(c.AuditFile) != ""
	c.AuditIncomeRemote = strings.TrimSpace(c.AuditURL) != ""

	if _, ok := os.LookupEnv(EnvHTTPInterface); ok {
		c.EnableHTTPS = true
	}

	fmt.Printf("AppConf after ENV: [%+v]\r\n", c)

	return c.copy(), nil
}

func parseFlag(env string, value flag.Value) error {
	var envVar = os.Getenv(env)
	if envVar == "" {
		return nil
	}

	fmt.Printf("[DEBUG] AppConf: ENV [%s] VALUE [%+v]\r\n", env, value)

	err := value.Set(envVar)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConf) initFlags() {
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
	flag.StringVar(&c.configPath, FlagConfigFile, "", "config file path")
}

func (c *AppConf) copy() *AppConf {
	return &AppConf{
		AppName:           c.AppName,
		ProjectStage:      c.ProjectStage,
		LogLevel:          c.LogLevel,
		BaseURL:           c.BaseURL,
		HTTP:              c.HTTP.Copy(),
		DBKind:            c.DBKind,
		DBDsn:             c.DBDsn,
		StoragePath:       c.StoragePath,
		StorageUserPath:   c.StorageUserPath,
		AuditIncomeLocal:  c.AuditIncomeLocal,
		AuditIncomeRemote: c.AuditIncomeRemote,
		EnableHTTPS:       c.EnableHTTPS,
		configPath:        c.configPath,
	}
}

func (c *AppConf) mergeWithPriority(conf *AppConf) {
	if conf == nil {
		return
	}

	c.AppName = conf.AppName
	c.ProjectStage = conf.ProjectStage
	c.LogLevel = conf.LogLevel
	c.BaseURL = conf.BaseURL
	c.HTTP = conf.HTTP.Copy()
	c.DBKind = conf.DBKind
	c.DBDsn = conf.DBDsn
	c.StoragePath = conf.StoragePath
	c.StorageUserPath = conf.StorageUserPath
	c.AuditIncomeLocal = conf.AuditIncomeLocal
	c.AuditIncomeRemote = conf.AuditIncomeRemote
	c.EnableHTTPS = conf.EnableHTTPS
	c.configPath = conf.configPath
}

func init() {
	flag.CommandLine.Init("App cmd line", flag.PanicOnError)

	intermediateConfig = defaultConfig()

	intermediateConfig.initFlags()
}
