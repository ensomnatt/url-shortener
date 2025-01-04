package handlers

import (
	"log/slog"
	"net/http"
	"strings"
)

func (h Handler) RouteActions(w http.ResponseWriter, r *http.Request) {
  value := r.URL.Query()
  action := value.Get("action") 
  alias := strings.ReplaceAll(r.URL.Path, "/", "")

  if action == "delete" {
    slog.Debug("route user to delete method")
    h.Delete(w, r, alias)
    return
  } else if action == "update" {
    slog.Debug("route user to update method")
    //h.Update(w, alias)
    return
  }

  slog.Debug("invalid action", "action", action)
  http.Error(w, "invalid action", http.StatusBadRequest)
  return
}
