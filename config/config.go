package config

import (
	"os"
)

type DBConfig struct {
	MongoURI string
	MongoDB  string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),
	}

	return &dbConfig
}
