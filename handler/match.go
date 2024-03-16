package handler

import (
	"bubblr/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var ErrUnique = errors.New("pq: duplicate key value violates unique constraint \"unique_like_pair\"")

func (h *Handler) AddMatch(w http.ResponseWriter, r *http.Request) {
	var m models.Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	matched, err := h.d.CheckForPrevMatched(m.UserID, m.MatchedUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matched, err = h.d.AddMatch(m.UserID, m.MatchedUserID, matched)
	if err == ErrUnique {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]bool{"matched": matched}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RemoveMatch(w http.ResponseWriter, r *http.Request) {
	var m models.Match
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.d.RemoveMatch(m.UserID, m.MatchedUserID); err != nil {
		http.Error(w, "Error adding match!", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetMatches(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userID := queryParams.Get("user_id")

	dbMatches, err := h.d.GetMatches(userID)
	if err != nil {
		http.Error(w, "Error adding match!", http.StatusInternalServerError)
		return
	}

	var matchesIDS []string
	var likedAt []time.Time
	for i := range dbMatches {
		matchesIDS = append(matchesIDS, dbMatches[i].Matches)
		likedAt = append(likedAt, dbMatches[i].LikedAt)
	}

	resp := models.Matches{
		UserID:  userID,
		Matches: matchesIDS,
		LikedAt: likedAt,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "error with user", http.StatusBadRequest)
	}

}
