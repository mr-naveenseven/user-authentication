package main

import (
	"fmt"
	"log"
	"user-authentication/internal/config"
	"user-authentication/internal/postgres"
	"user-authentication/internal/router"
	logger "user-authentication/pkg"
)

func main() {
	// Initialize the logger (modes: "json" or "text")
	logger.Init(logger.LOGGER_MODE_TEXT)

	// Load server configurations
	config := config.NewServerConfig()
	err := config.LoadConfigs()
	if err != nil {
		log.Println("Failed to load server configurations:", err)

		return
	}

	log.Println("Starting the application...")

	// Postgres connection string
	pgConnAddress := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		config.PGConfig.User,
		config.PGConfig.Password,
		config.PGConfig.Host,
		config.PGConfig.Port,
		config.PGConfig.DBName,
		config.PGConfig.SSLMode,
	)

	// Migration connection string
	pgConnMigAddress := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.PGConfig.User,
		config.PGConfig.Password,
		config.PGConfig.Host,
		config.PGConfig.Port,
		config.PGConfig.DBName,
		config.PGConfig.SSLMode,
	)

	// postgresClient represents the Postgres client connection.
	postgresClient := postgres.NewPGClient(pgConnAddress, pgConnMigAddress, config.PGConfig.ConnTimeout, config.MigrationDir)
	err = postgresClient.Connect()
	if err != nil {
		println("Failed to connect to Postgres database:", err)

		return
	}

	// Initialize the Postgres client.
	if err := postgresClient.Initialize(); err != nil {
		println("Failed to initialize Postgres client:", err)

		return
	}
	// defer postgres
	defer postgresClient.Disconnect()

	// Set up the router and start the server
	rEngine := router.SetupRouter()
	router.InitRoutes(rEngine)
	rEngine.Run(":" + config.Port)
}
