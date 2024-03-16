package handler

import (
	"bubblr/convert"
	"bubblr/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Signup with email and password
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {

	var uc UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&uc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	generatedID := uuid.NewString()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(uc.Password), 10)
	if err != nil {
		http.Error(w, "error in generating hashed password", http.StatusConflict)
		return
	}

	dbUser, err := h.d.GetUserByEmail(uc.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	user := convert.ConvertDBUserToUser(dbUser)

	if len(user.UserID) > 0 {
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write([]byte("user already exists, please log in"))
		return
	}

	sanitizedEmail := strings.ToLower(uc.Email)

	signupCreds := &models.User{
		UserID:         generatedID,
		Email:          sanitizedEmail,
		HashedPassword: string(hashedPassword),
	}

	insertedUser, err := h.d.InsertUser(signupCreds)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in inserting user %v", err), http.StatusBadRequest)
		return
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
	}{
		Token:  token,
		UserID: insertedUser,
	}

	userJSON, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userJSON)
	if err != nil {
		http.Error(w, "error on msg", http.StatusBadRequest)
	}
}
