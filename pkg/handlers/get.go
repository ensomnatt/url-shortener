package handlers

import (
	"log/slog"
	"net/http"
	"strings"
)

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
  alias := strings.ReplaceAll(r.URL.Path, "/", "")
  link, err := h.db.Get("link", "urls", "alias", alias)
  if err != nil {
    slog.Error("failed to get link", "error", err, "alias", alias)
    http.Error(w, "failed to get link", http.StatusInternalServerError)
    return
  }

  slog.Info("redirected user", "alias", alias, "link", link)
  http.Redirect(w, r, link, http.StatusMovedPermanently)
}
