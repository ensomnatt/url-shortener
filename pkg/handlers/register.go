package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"urlshortener/pkg/models"
)

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
  var user models.User
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
    slog.Error("failed to get user credentials", "error", err)
    http.Error(w, "invalid credentials", http.StatusBadRequest)
  }

  slog.Debug("got request")

  x, err := h.db.Check("users", "username", user.Username) 
  if err != nil {
    slog.Error("failed to check username", "error", err, "username", user.Username)
    http.Error(w, "failed to check username", http.StatusUnauthorized)
    return
  }

  if x {
    slog.Debug("user is already exists", "username", user.Username)
    http.Error(w, "user with this username is already exists", http.StatusConflict)
    return
  }

  hashedPassword, err := h.crypter.Password([]byte(user.Password))
  if err != nil {
    slog.Error("failed to hash password", "error", err)
    http.Error(w, "failed to hash your password", http.StatusInternalServerError)
    return
  }

  err = h.db.Save("users", "username", "password", user.Username, hashedPassword)
  if err != nil {
    slog.Error("failed to save user to the db", "error", err)
    http.Error(w, "registration was failed", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

  err = json.NewEncoder(w).Encode(user)   
  if err != nil {
    slog.Error("failed to send response to user", "error", err)
    http.Error(w, "failed to send response", http.StatusInternalServerError)
    return
  }

  slog.Info("created new user", "username", user.Username)
  
  return
}
