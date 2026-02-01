package config

// Flags
const (
	FlagAppName           string = "p"
	FlagProjectStage      string = "stage"
	FlagLogLevel          string = "l"
	FlagBaseURL           string = "b"
	FlagDBKind            string = "k"
	FlagHTTPInterface     string = "a"
	FlagDBInterface       string = "d"
	FlagStoragePath       string = "f"
	FlagStorageUserPath   string = "fu"
	FlagAuditFile         string = "audit-file"
	FlagAuditURL          string = "audit-url"
	FlagEnableHTTPS       string = "s"
	FlagConfigPath        string = "c"
	FlagConfigBigPath     string = "config"
	FlagTrustedSubnetCIDR string = "t"
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
	EnvConfig              string = "CONFIG"
	EnvTrustedSubnetCIDR   string = "TRUSTED_SUBNET"
)

// Database kinds
const (
	DBKindInMemory string = "IN_MEMORY"
	DBKindPostgres string = "POSTGRES"
)

// Project stages
const (
	ProjectStageProduction  string = "PROD"
	ProjectStageDevelopment string = "DEV"
)
