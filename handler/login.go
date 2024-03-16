package handler

import (
	"bubblr/convert"

	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login with email and password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	var uc UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&uc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	dbUser, err := h.d.GetUserByEmail(uc.Email)
	if err != nil || len(dbUser.Email) == 0 {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	user := convert.ConvertDBUserToUser(dbUser)

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(uc.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// If not verified return 401
	if !user.Verified {
		result := h.checkVerified(w, r, user.UserID)
		if !result {
			return
		}
	}

	// Generate JWT token
	token, err := generateJWTToken(uc.Email)
	if err != nil {
		http.Error(w, "Error in JWT token", http.StatusBadRequest)
		return
	}

	// user.TokenExpiration = time.Now().Add(time.Minute * 24).Unix()
	response := struct {
		Token  string `json:"token"`
		UserID string `json:"userId"`
		IsMale bool   `json:"isMale"`
	}{
		Token:  token,
		UserID: user.UserID,
		IsMale: user.IsMale,
	}

	userJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, "error on msg", http.StatusInternalServerError)
	}
}

func (h *Handler) checkVerified(w http.ResponseWriter, r *http.Request, user_id string) bool {
	// if not enough votes but not rejected return 206
	vote_count, err := h.d.GetVoteCount(user_id)
	if err != nil {
		http.Error(w, "Error in get vote count", http.StatusInternalServerError)
		return false
	}
	rejected_count, err := h.d.GetRejectionVoteCount(user_id)
	if err != nil {
		http.Error(w, "Error in get vote count", http.StatusInternalServerError)
		return false
	}

	if rejected_count >= 5 {
		// You are rejected 205
		w.WriteHeader(http.StatusResetContent)
		return false
	}

	if vote_count >= 5 {
		// You are valid! verify and continue
		if err := h.d.VerifyUser(user_id); err != nil {
			http.Error(w, "error verifying user!", http.StatusInternalServerError)
			return false
		}
		return true
	}

	response := struct {
		VoteCount int    `json:"voteCount"`
		UserID    string `json:"userId"`
	}{
		VoteCount: vote_count,
		UserID:    user_id,
	}

	userJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPartialContent)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, "error on msg", http.StatusBadRequest)
	}
	return false
}
