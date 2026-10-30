package guestmessages

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"backend/internal/httputil"
	"backend/internal/repo"
	"backend/internal/sheets"
)

type submitRequest struct {
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	slug := httputil.ExtractSlug(r)
	if slug == "" {
		httputil.WriteError(w, http.StatusBadRequest, "slug required")
		return
	}

	var req submitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid payload")
		return
	}

	msg := strings.TrimSpace(req.Message)
	if msg == "" {
		httputil.WriteError(w, http.StatusBadRequest, "message cannot be empty")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	rp, err := repo.New(ctx)
	if err != nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "backend unavailable")
		return
	}

	if len(msg) > rp.Cfg.MaxMessageLength {
		httputil.WriteError(w, http.StatusBadRequest, "message too long")
		return
	}

	// Resolve guest identity from the URL slug — never trust client-supplied names.
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

	ipHash := sheets.HashIP(httputil.ClientIP(r))
	window := time.Duration(rp.Cfg.RateLimitSeconds) * time.Second
	if recent, err := rp.Client.RecentIPSubmission(ctx, rp.Cfg.MessagesTab, ipHash, window); err == nil && recent {
		httputil.WriteError(w, http.StatusTooManyRequests, "please wait before submitting again")
		return
	}

	record := sheets.Message{
		MessageID:   sheets.NewMessageID(),
		GuestID:     guest.GuestID,
		DisplayName: guest.GuestName,
		Message:     msg,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := rp.Client.AppendMessage(ctx, rp.Cfg.MessagesTab, record, ipHash); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to save message")
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, record)
}
