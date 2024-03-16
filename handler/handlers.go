package handler

import (
	"bubblr/datastore"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	d *datastore.Datastore
}

func NewHandler(d *datastore.Datastore) *Handler {
	return &Handler{
		d: d,
	}
}

var SigningKey = []byte("bubblr_stay_in_your_bubble") // Change this to a secure secret key

func generateJWTToken(email string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 48).Unix(), // Expires in 24 minutes
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
