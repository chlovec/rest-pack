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
	BaseUrl       string
	PathPrefix    string
}

var Envs Config

func InitConfig() {
	// Note: Best practice .env file should be in the root folder.
	// It's put in the examples folder as the root directory these examples.
	godotenv.Load()
	Envs = Config {
		DBHost:               os.Getenv("DB_HOST"),
		DBPort:               os.Getenv("DB_PORT"),
		DBUser:               os.Getenv("DB_USER"),
		DBPassword:           os.Getenv("DB_PASSWORD"),
		DBName:               os.Getenv("DB_NAME"),
		ServerAddress:        os.Getenv("SERVER_ADDR"),
		BaseUrl:              os.Getenv("BASE_URL"),
		PathPrefix:           os.Getenv("PATH_PREFIX"),
	}
}

func GetDataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?checkConnLiveness=false&parseTime=true", Envs.DBUser, Envs.DBPassword, Envs.DBHost, Envs.DBPort, Envs.DBName)
}