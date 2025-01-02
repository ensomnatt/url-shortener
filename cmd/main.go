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
  dbPassword := os.Getenv("DB_PASSWORD")
  secret := []byte(os.Getenv("JWT_SECRET"))
  
  logger.Init(env)
  slog.Debug("logger is running")

  db, err := database.Init(dbPassword)
  if err != nil {
    slog.Error("failed to open db", "error", err)
    panic(err)
  } 
  defer db.Close()

  handlers.Start(db, secret)
}
