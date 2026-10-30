// Local development server.
// NOT deployed to Vercel — production uses api/*/index.go directly.
//
// Run with:
//
//	cd backend && go run ./cmd/dev
package main

import (
	"log"
	"net/http"

	guest "backend/api/guest"
	guestmessages "backend/api/guest_messages"
	"backend/api/health"
	"backend/api/messages"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load("../.env") // repo-root .env
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", health.Handler)
	mux.HandleFunc("/api/messages", messages.Handler)
	mux.HandleFunc("/api/guests/", guestOrGuestMessages)

	addr := ":3000"
	log.Printf("Go API listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

// guestOrGuestMessages dispatches /api/guests/:slug and
// /api/guests/:slug/messages to the right handler based on path shape.
func guestOrGuestMessages(w http.ResponseWriter, r *http.Request) {
	// Path ends with /messages → guest message submission
	if len(r.URL.Path) >= len("/messages") &&
		r.URL.Path[len(r.URL.Path)-len("/messages"):] == "/messages" {
		guestmessages.Handler(w, r)
		return
	}
	guest.Handler(w, r)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
