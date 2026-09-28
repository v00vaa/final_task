package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/v00vaa/final_task/pkg/db"
	"github.com/v00vaa/final_task/pkg/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err = db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server.Run()
}
