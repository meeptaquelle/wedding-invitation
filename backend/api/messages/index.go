package messages

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePostAnonymous(w, r)
	default:
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	rp, err := repo.New(ctx)
	if err != nil {
		httputil.WriteError(w, http.StatusServiceUnavailable, "backend unavailable")
		return
	}

	msgs, err := rp.Client.LoadMessages(ctx, rp.Cfg.MessagesTab)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to load messages")
		return
	}

	// Newest first
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	limit := rp.Cfg.PageSize
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	start := 0
	if cursor := r.URL.Query().Get("cursor"); cursor != "" {
		for i, m := range msgs {
			if m.MessageID == cursor {
				start = i + 1
				break
			}
		}
	}
	end := start + limit
	if end > len(msgs) {
		end = len(msgs)
	}

	page := msgs[start:end]
	var next *string
	if end < len(msgs) && len(page) > 0 {
		c := page[len(page)-1].MessageID
		next = &c
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"messages":    page,
		"next_cursor": next,
	})
}

func handlePostAnonymous(w http.ResponseWriter, r *http.Request) {
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

	ipHash := sheets.HashIP(httputil.ClientIP(r))
	window := time.Duration(rp.Cfg.RateLimitSeconds) * time.Second
	if recent, err := rp.Client.RecentIPSubmission(ctx, rp.Cfg.MessagesTab, ipHash, window); err == nil && recent {
		httputil.WriteError(w, http.StatusTooManyRequests, "please wait before submitting again")
		return
	}

	record := sheets.Message{
		MessageID:   sheets.NewMessageID(),
		GuestID:     "",
		DisplayName: "Anonymous",
		Message:     msg,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := rp.Client.AppendMessage(ctx, rp.Cfg.MessagesTab, record, ipHash); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, "failed to save message")
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, record)
}
