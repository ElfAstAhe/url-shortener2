package config

const (
	DefaultAppName  string = "URL shorter"
	DefaultLogLevel string = "INFO"
	DefaultStage           = ProjectStageDevelopment
	DefaultBaseURL  string = "http://localhost:8080"
)

// server listener
const (
	DefaultHTTPHost string = "localhost"
	DefaultHTTPPort int    = 8080
)

// database
const (
	DefaultDBKind        = DBKindPostgres
	DefaultDBDsn  string = ""
)

// storage
const (
	DefaultStoragePath     = "./shortener.txt"
	DefaultStorageUserPath = "./shortener_user.txt"
)

// income audit

const (
	DefaultAuditIncomePath = "./audit_income.txt"
)
