package main

import (
	"fmt"
	"log"
	"os"
	"rms/database"
	"rms/server"

	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file" + " " + err.Error())
	}

	// Connect to database and run migrations
	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	); err != nil {
		fmt.Println("Failed to connect to database or run migrations")
		log.Fatal(err)
	}

	//Setup routes
	svc := server.SetupRoutes()

	//server run
	if err := svc.Run(":" + os.Getenv("PORT")); err != nil {
		fmt.Println("Failed to start server")
		log.Fatal(err)
	}

	// Close database
	if err := database.ShutdownDatabase(); err != nil {
		fmt.Println("Failed to close database")
		log.Fatal(err)
	}

}
