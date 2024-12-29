package main

import (
	"log/slog"
	"os"
	"urlshortener/pkg/database"
	"urlshortener/pkg/handlers"
	"urlshortener/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
  godotenv.Load(".env") 
  env := os.Getenv("ENV")
  connStr := os.Getenv("DB_CONNECTION_STRING")
  
  logger.Init(env)
  slog.Debug("logger is running")

  db, err := database.Init(connStr)
  if err != nil {
    slog.Error("failed to open db", "error", err)
    panic(err)
  } 
  defer db.Close()

  handlers.Start(db)
}
