package handlers

import (
	"log/slog"
	"net/http"
	"urlshortener/pkg/crypt"
	"urlshortener/pkg/database"
	"urlshortener/pkg/tokener"
)

type Handler struct {
  db *database.Storage
  crypter *crypt.Crypt
  tokener *tokener.Tokener
}

func Start(db *database.Storage, secret []byte) {
  r := http.NewServeMux()
  crypter := crypt.Create()
  tokener := tokener.Create(secret)
  h := Handler{
    db: db,
    crypter: crypter,
    tokener: tokener,
  }

  r.HandleFunc("POST /shorten", h.Save)
  r.HandleFunc("GET /{alias}", h.Get)
  r.HandleFunc("POST /register", h.Register)
  r.HandleFunc("POST /login", h.Login)

  slog.Info("server is up and running")
  http.ListenAndServe(":8181", r)
}
