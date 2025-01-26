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
	DBName        string
	ServerAddress string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load("examples/.env")
	return Config {
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		ServerAddress:        os.Getenv("SERVER_ADDR"),
	}
}

func GetDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?checkConnLiveness=false&parseTime=true", Envs.DBUser, Envs.DBPassword, Envs.DBHost, Envs.DBPort, Envs.DBName)
}