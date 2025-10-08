package main

import (
	"log"
	"user-authentication/internal/config"
	"user-authentication/internal/postgres"
	"user-authentication/internal/router"
	logger "user-authentication/pkg"
)

func main() {
	// Initialize the logger (modes: "json" or "text")
	logger.Init("text")

	// Load server configurations
	config := config.NewServerConfig()
	config.LoadConfigs()

	log.Println("Starting the application...")

	// postgresClient represents the Postgres client connection.
	postgresClient := postgres.NewPGClient(config.PGConfig.ConnTimeout, config.MigrationDir)
	pgConnAddress := config.PGConfig.User + ":" + config.PGConfig.Password + "@" + config.PGConfig.Host + ":" + config.PGConfig.Port +
		"/" + config.PGConfig.DBName + "?sslmode=" + config.PGConfig.SSLMode

	err := postgresClient.Connect("postgres://" + pgConnAddress)
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
