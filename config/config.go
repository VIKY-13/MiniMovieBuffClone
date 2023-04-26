package config

import (
	"os"
	"github.com/joho/godotenv"
)

func init() {
	// Load the environment variables from the .env file
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	HOST_PORT		= os.Getenv("HOST_PORT")
	DB_HOST			= os.Getenv("DB_HOST")
	DB_PORT			= os.Getenv("DB_PORT")
	DB_USER			= os.Getenv("DB_USER")
	DB_PASSWORD		= os.Getenv("DB_PASSWORD")
	DB_NAME			= os.Getenv("DB_NAME")
	AUTH_USERNAME	= os.Getenv("AUTH_USERNAME")
	AUTH_PASSWORD	= os.Getenv("AUTH_PASSWORD")
}

var(
	HOST_PORT,
	DB_HOST,
	DB_PORT,
	DB_USER,
	DB_PASSWORD,
	DB_NAME,
	AUTH_USERNAME,
	AUTH_PASSWORD string
)