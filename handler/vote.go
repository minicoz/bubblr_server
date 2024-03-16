package handler

import (
	"bubblr/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// This package should be used to change the prospective user to verified when the voting happens.
// High level, this endpoint will be shown to a verified user to vote on whether they want the
// prospective user in.

func (h *Handler) GetVote(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK healthy")
}

func (h *Handler) AddVote(w http.ResponseWriter, r *http.Request) {
	var v *models.Vote
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload, %v", err), http.StatusBadRequest)
		return
	}

	if err := h.d.AddVote(v.UserID, v.MatchedUserID); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK healthy")
}

func (h *Handler) RmVote(w http.ResponseWriter, r *http.Request) {
	var v *models.Vote
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload, %v", err), http.StatusBadRequest)
		return
	}

	if err := h.d.RmVote(v.UserID, v.MatchedUserID); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK healthy")
}
