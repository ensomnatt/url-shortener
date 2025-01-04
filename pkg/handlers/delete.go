package handlers

import (
	"log/slog"
	"net/http"
)

func (h Handler) Delete(w http.ResponseWriter, r *http.Request, alias string) { 
  _, owner, err := h.db.Get(alias)
  if err != nil {
    slog.Error("failed to get links owner", "error", err)
    http.Error(w, "failed to serve your request", http.StatusInternalServerError)
    return
  }

  username, err := h.tokener.ValidateToken(r)
  if err != nil {
    slog.Debug("failed to validate token", "error", err)
    http.Error(w, "invalid token", http.StatusUnauthorized)
    return
  }

  if username != owner {
    slog.Debug("username is not equal owner", "username", username, "owner", owner)
    http.Error(w, "it's not your link", http.StatusUnauthorized)
    return
  }

  err = h.db.Delete(alias)
  if err != nil {
    slog.Error("failed to delete link", "error", err)
    http.Error(w, "failed to delete link", http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("deleted your link"))

  slog.Debug("sent response to user")

  return
}
