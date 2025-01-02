package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http" 
  "urlshortener/pkg/models")

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
  var user models.User 
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
    slog.Error("failed to get request", "error", err)
    http.Error(w, "invalid credentials", http.StatusBadRequest)
    return
  }

  x, err := h.db.Check("users", "username", user.Username)
  if err != nil {
    slog.Error("failed to check username", "error", err)
    http.Error(w, "failed to check username", http.StatusUnauthorized)
    return
  }

  if !x {
    slog.Debug("user doesn't exist", "username", user.Username)
    http.Error(w, "invalid login or password", http.StatusUnauthorized)
  }

  hash, err := h.db.Get("password", "users", "username", user.Username)
  if err != nil {
    slog.Error("failed to get hash password from db", "error", err)
    http.Error(w, "invalid login or password", http.StatusUnauthorized)
    return
  }
  success := h.crypter.GetPassword([]byte(hash), []byte(user.Password))
  
  if !success {
    slog.Debug("failed login try", "username", user.Username)
    http.Error(w, "invalid login or password", http.StatusUnauthorized)
    return
  } 

  token, err := h.tokener.GenToken(user.Username)
  if err != nil {
    slog.Error("failed to generate token", "error", err, "username", user.Username)
    http.Error(w, "failed to generate token", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  err = json.NewEncoder(w).Encode(map[string]string{
    "token": token,
  })
  if err != nil {
    slog.Error("failed to send response to user", "error", err)
    http.Error(w, "failed to send your token", http.StatusInternalServerError)
    return
  }

  slog.Info("successfully sent token to user", "username", user.Username)
}
