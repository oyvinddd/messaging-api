package app

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	listeningPort 	string
	mailerAPIKey 	string
	dbURI 			string
	migrationPath	string
	hmacSecret		string
	serviceKey 		string
}

func NewDefaultConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load environment: %v\n", err)
	}

	cfg := Config{}
	cfg.listeningPort 	= varFromEnv("LISTENING_PORT")
	cfg.mailerAPIKey	= varFromEnv("POSTMARK_API_KEY")
	cfg.dbURI			= varFromEnv("POSTGRES_URI")
	cfg.migrationPath	= varFromEnv("MIGRATION_PATH")
	cfg.hmacSecret		= varFromEnv("HMAC_SECRET")
	cfg.serviceKey		= varFromEnv("SERVICE_KEY")
	return &cfg
}

func varFromEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing environment var for key: %s\n", key)
	}
	return value
}
