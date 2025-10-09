package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	logger "user-authentication/pkg"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PGClient represents the Postgres client connection.
type PGClient struct {
	DB               *gorm.DB
	Timeout          int
	GormConnURI      string
	MigrationConnURL string
	MigrationDir     string
}

// NewPGClient creates a new instance of PGClient.
func NewPGClient(gormConnURI string, migrationConnURL string, timeout int, migrationDir string) *PGClient {

	return &PGClient{
		DB:               nil,
		GormConnURI:      gormConnURI,
		MigrationConnURL: migrationConnURL,
		Timeout:          timeout,
		MigrationDir:     migrationDir,
	}
}

// Connect establishes a connection to the Postgres database.
func (pg *PGClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pg.Timeout)*time.Second)
	defer cancel()

	db, err := gorm.Open(gormpg.Open(pg.GormConnURI), &gorm.Config{
		// Logger: gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {

		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {

		return fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	// Set connection pool configs
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Test ping
	if err := sqlDB.PingContext(ctx); err != nil {

		return fmt.Errorf("database ping failed: %w", err)
	}

	pg.DB = db
	logger.Info("Connected to postgres database using GORM.....")

	return nil
}

// Initialize initializes the Postgres client by migrating the database tables.
func (pg *PGClient) Initialize() error {

	// Migrate the database tables
	if err := pg.migrateTables(); err != nil {
		log.Println("Error migrating tables:", err)

		return err
	}
	log.Println("Postgres client initialized successfully.")

	return nil
}

// nugrateTables applies database migrations using the golang-migrate library.
func (pg *PGClient) migrateTables() error {

	// Create a new migrate instance
	m, err := migrate.New("file://"+pg.MigrationDir, pg.MigrationConnURL)
	if err != nil {
		log.Println("Error creating migration instance:", err)

		return fmt.Errorf("error creating migration instance: %w", err)
	}
	defer m.Close()

	// Apply all up migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println("Error applying migrations:", err)

		return fmt.Errorf("error applying migrations: %w", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No migrations to apply.")

		return nil
	} else if err == nil {
		log.Println("Migrations applied successfully.")
	}

	return nil

}

// Disconnect closes the connection to the Postgres database.
func (pg *PGClient) Disconnect() {
	if pg.DB != nil {
		sqlDB, _ := pg.DB.DB()
		sqlDB.Close()
		log.Println("Disconnected from Postgres database successfully.")
	} else {
		log.Println("Postgres client is already disconnected or was never connected.")
	}
}
