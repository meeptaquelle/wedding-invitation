package guest

import (
	"context"
	"net/http"
	"time"

	"backend/internal/httputil"
	"backend/internal/repo"
	"backend/internal/sheets"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	slug := httputil.ExtractSlug(r)
	if slug == "" {
		httputil.WriteError(w, http.StatusBadRequest, "slug required")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	rp, err := repo.New(ctx)
	if err != nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "backend unavailable")
		return
	}

	guests, err := rp.Client.CachedGuests(ctx, rp.Cfg.GuestListTab)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to load guests")
		return
	}

	guest := sheets.FindBySlug(guests, slug)
	if guest == nil {
		httputil.WriteError(w, http.StatusNotFound, "guest not found")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, guest)
}
