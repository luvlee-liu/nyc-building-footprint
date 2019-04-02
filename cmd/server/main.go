package main

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	server := App{}

	var (
		dbUser     = getEnv("DB_USER", "postgres")
		dbPassword = getEnv("DB_PASSWORD", "")
		dbName     = getEnv("DB_NAME", "nycbuilding")
		dbHost     = getEnv("DB_HOST", "localhost")
		dbPort     = getEnv("DB_PORT", "5432")
		dbSSLMode  = getEnv("DB_SSL_MODE", "disable")
		serverPort = getEnv("PORT", "8080")
	)

	if dbPassword == "" {
		fmt.Println("env variable DB_PASSWORD is required, cannot be empty")
		os.Exit(1)
	}

	// set db connection from envvar
	server.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	fmt.Println("server listening on port " + serverPort)
	server.Run(":" + serverPort)
}
