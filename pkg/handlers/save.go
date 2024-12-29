package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"urlshortener/pkg/database"
)

type Request struct {
  Alias string `json:"alias"`
  Link string `json:"link"`
}

type Response struct {
  Link string `json:"link"`
}

func (h Handler) check(link string) (bool, error) {
  resp, err := http.Get(link)
  if err != nil {
    return false, err
  }

  if resp.StatusCode >= 400 {
    return false, nil
  }

  slog.Debug("checked link")
  return true, nil
}

func (h Handler) Save(w http.ResponseWriter, r *http.Request) {
  //get request
  var req Request
  err := json.NewDecoder(r.Body).Decode(&req)
  if err != nil {
    slog.Error("failed to get user request", "error", err)
    http.Error(w, "failed to get your request", http.StatusInternalServerError)
    return
  }
  slog.Debug("get request from user", "alias", req.Alias, "link", req.Link)

  //check link
  x, err := h.check(req.Link)
  if !x {
    if err != nil {
      slog.Error("failed to check link", "error", err, "link", req.Link)
      http.Error(w, "failed to check link", http.StatusInternalServerError)
      return
    } else {
      slog.Debug("link is doesn't work or exists", "link", req.Link)
      http.Error(w, "link is doesn't work or exists", http.StatusBadRequest)
      return
    }
  }

  //save link
  err = h.db.Save(req.Alias, req.Link)
  if err != nil {
    if err == database.AliasExists {
      slog.Debug("user's alias is already exists", "alias", req.Alias)
      http.Error(w, "alias is already exists", http.StatusConflict)
      return
    } else {
      slog.Error("failed to save link", "error", err, "alias", req.Alias, "link", req.Link)
      http.Error(w, "failed to save link", http.StatusInternalServerError)
      return
    }
  }

  //send response
  resp := Response{
    Link: "150.241.82.204:8181/" + req.Alias,
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)

  err = json.NewEncoder(w).Encode(resp)
  if err != nil {
    slog.Error("failed to send response", "error", err)
    http.Error(w, "failed to send response to you", http.StatusInternalServerError)
    return
  }
  slog.Debug("sent response to user")

  return
}
