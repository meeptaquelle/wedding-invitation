package health

import (
	"net/http"
	"time"

	"backend/internal/httputil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
