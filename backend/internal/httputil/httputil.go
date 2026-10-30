package httputil

import (
	"encoding/json"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// ExtractSlug supports both Vercel's query-param injection and direct path parsing.
func ExtractSlug(r *http.Request) string {
	if s := r.URL.Query().Get("slug"); s != "" {
		return s
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, p := range parts {
		if p == "guests" && i+1 < len(parts) {
			next := parts[i+1]
			if next != "messages" {
				return next
			}
		}
	}
	return ""
}

func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return r.RemoteAddr
}

func HandleOptions(w http.ResponseWriter) bool {
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	return true
}
