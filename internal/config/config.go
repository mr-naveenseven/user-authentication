package config

import (
	"os"
	"strconv"
	logger "user-authentication/pkg"

	"github.com/joho/godotenv"
)

// PGConfig holds the Postgres database configuration details.
type PGConfig struct {
	User        string
	Password    string
	Host        string
	Port        string
	DBName      string
	SSLMode     string
	ConnTimeout int
}

// ServerConfig holds the overall server configuration including Postgres config.
type ServerConfig struct {
	Host         string
	Port         string
	MigrationDir string
	PGConfig     PGConfig
}

// default configuration path
const ENV_SERVER_CONFIG = "configs/server_config.env"

const (
	// server .evn config keys
	KEY_SERVER_HOST   = "SERVER_HOST"
	KEY_SERVER_PORT   = "SERVER_PORT"
	KEY_MIGRATION_DIR = "MIGRATION_DIR"

	// postgres .env config keys
	KEY_DB_HOST         = "DB_HOST"
	KEY_DB_PORT         = "DB_PORT"
	KEY_DB_USER         = "DB_USER"
	KEY_DB_PASSWORD     = "DB_PASSWORD"
	KEY_DB_NAME         = "DB_NAME"
	KEY_DB_SSLMODE      = "DB_SSLMODE"
	KEY_DB_CONN_TIMEOUT = "DB_CONN_TIMEOUT"

	// jwt .env config keys
	KEY_JWT_SECRET         = "JWT_SECRET"
	KEY_JWT_EXPIRE_MINUTES = "JWT_EXPIRE_MINUTES"
)

// NewServerConfig creates a new instance of ServerConfig with default values.
func NewServerConfig() ServerConfig {

	return ServerConfig{
		Host:         "",
		Port:         "",
		MigrationDir: "",
		PGConfig:     PGConfig{},
	}
}

// loadServerConfig loads server-related configurations from environment variables.
func (sc *ServerConfig) loadServerConfig() {
	sc.Host = os.Getenv(KEY_SERVER_HOST)
	sc.Port = os.Getenv(KEY_SERVER_PORT)
	sc.MigrationDir = os.Getenv(KEY_MIGRATION_DIR)
}

// loadPostgresConfig loads Postgres-related configurations from environment variables.
func (sc *ServerConfig) loadPostgresConfig() {
	sc.PGConfig.Host = os.Getenv(KEY_DB_HOST)
	sc.PGConfig.Port = os.Getenv(KEY_DB_PORT)
	sc.PGConfig.User = os.Getenv(KEY_DB_USER)
	sc.PGConfig.Password = os.Getenv(KEY_DB_PASSWORD)
	sc.PGConfig.DBName = os.Getenv(KEY_DB_NAME)
	sc.PGConfig.SSLMode = os.Getenv(KEY_DB_SSLMODE)
	sc.PGConfig.ConnTimeout, _ = strconv.Atoi(os.Getenv(KEY_DB_CONN_TIMEOUT))
}

// LoadConfigs loads all configurations from the specified .env file.
func (sc *ServerConfig) LoadConfigs() error {
	if err := godotenv.Load(ENV_SERVER_CONFIG); err != nil {
		logger.Error("Loading config from env file failed", "error", err)

		return err
	}

	sc.loadServerConfig()
	sc.loadPostgresConfig()

	logger.Info("Server configuration loaded", "config", sc)

	return nil
}
