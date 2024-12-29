package handlers

import (
	"log/slog"
	"net/http"
	"urlshortener/pkg/database"
)

type Handler struct {
  db *database.Storage
}

func Start(db *database.Storage) {
  r := http.NewServeMux()
  h := Handler{
    db: db,
  }

  r.HandleFunc("POST /shorten", h.Save)
  r.HandleFunc("GET /{alias}", h.Get)

  slog.Info("server is up and running")
  http.ListenAndServe(":8181", r)
}
