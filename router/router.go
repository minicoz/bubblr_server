package router

import (
	"bubblr/handler"
	"bubblr/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK healthy")
}

// Router is exported and used in main.go
func Router(h *handler.Handler) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.EnableCORS)
	r.Use(middleware.JWTMiddleware)
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/", healthCheckHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/login", h.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", h.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/user", h.GetUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/user", h.CompleteUser).Methods("PUT", "OPTIONS")
	r.HandleFunc("/users", h.GetUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/unverified", h.GetUnverifiedUsers).Methods("GET", "OPTIONS")
	r.HandleFunc("/schools", h.GetSchools).Methods("GET", "OPTIONS")
	r.HandleFunc("/schools/tier", h.GetSchoolsWithTier).Methods("GET", "OPTIONS")
	r.HandleFunc("/gendered-users", h.GetGenderedUser).Methods("GET", "OPTIONS")
	r.HandleFunc("/add-match", h.AddMatch).Methods("POST", "OPTIONS")
	r.HandleFunc("/rm-match", h.RemoveMatch).Methods("POST", "OPTIONS")
	r.HandleFunc("/get-matches", h.AddMatch).Methods("GET", "OPTIONS")
	r.HandleFunc("/messages", h.GetMessages).Methods("GET", "OPTIONS")
	r.HandleFunc("/latest-messages", h.GetLatestMessages).Methods("GET", "OPTIONS")
	r.HandleFunc("/messages", h.AddMessage).Methods("POST", "OPTIONS")

	r.HandleFunc("/vote", h.GetVote).Methods("GET", "OPTIONS")
	r.HandleFunc("/vote", h.AddVote).Methods("POST", "OPTIONS")
	r.HandleFunc("/rm-vote", h.RmVote).Methods("POST", "OPTIONS")
	// r.HandleFunc("/rm-vote", h.RmVote).Methods("POST", "OPTIONS")

	return r
}
