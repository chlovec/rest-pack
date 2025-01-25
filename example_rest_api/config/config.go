package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword 	  string
	DBAddress     string
	DBName        string
	ServerAddress string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load("example_rest_api/.env")
	return Config {
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBAddress:            os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		ServerAddress:        os.Getenv("SERVER_ADDR"),
	}
}

func GetDataSourceName() string {
	// username:password@tcp(host:port)/database?parseTime=true&AllowNativePasswords=true
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&AllowNativePasswords=true", Envs.DBUser, Envs.DBPassword, Envs.DBHost, Envs.DBPort, Envs.DBName)
}