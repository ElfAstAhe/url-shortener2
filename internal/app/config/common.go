package config

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
	FlagConfigFile      string = "c"
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
	EnvConfigFile          string = "CONFIG"
)

// supported databases
const (
	DBKindInMemory string = "IN_MEMORY"
	DBKindPostgres string = "POSTGRES"
)

// supported project stages
const (
	ProjectStageProduction  = "PROD"
	ProjectStageDevelopment = "DEV"
)
