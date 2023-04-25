package config

import (
	"os"
)

var(
	PORT			= os.Getenv("PORT")
	DB_HOST			= os.Getenv("DB_HOST")
	DB_PORT			= os.Getenv("DB_PORT")
	DB_USER			= os.Getenv("DB_USER")
	DB_PASSWORD		= os.Getenv("DB_PASSWORD")
	DB_NAME			= os.Getenv("DB_NAME")
	AUTH_USERNAME	= os.Getenv("AUTH_USERNAME")
	AUTH_PASSWORD	= os.Getenv("AUTH_PASSWORD")
)