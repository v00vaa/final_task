package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/v00vaa/final_task/pkg/db"
	"github.com/v00vaa/final_task/pkg/server"
)

type Config struct {
	Port     string
	DBFile   string
	Password string
}

func loadConfig() Config {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	return Config{
		Port:     port,
		DBFile:   dbFile,
		Password: os.Getenv("TODO_PASSWORD"),
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}
	cfg := loadConfig()

	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(cfg.Port, cfg.Password); err != nil {
		log.Println(err)
	}
}
