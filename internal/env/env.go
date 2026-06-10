package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}

	if os.Getenv("GO_ENV") != "production" {

		err := godotenv.Load("../.env")
		if err == nil {
			log.Printf("Loaded environment from %s", "../.env")
		}

		val, ok = os.LookupEnv(key)
		if ok {
			return val
		}
	}

	return fallback
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if ok {
		valAsInt, err := strconv.Atoi(val)
		if err == nil {
			return valAsInt
		}
	}

	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load("../.env")
		if err == nil {
			log.Printf("Loaded environment from %s", "../.env")
		}

		val, ok = os.LookupEnv(key)
		if ok {
			valAsInt, err := strconv.Atoi(val)
			if err == nil {
				return valAsInt
			}
		}
	}

	return fallback
}
