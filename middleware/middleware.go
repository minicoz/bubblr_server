package middleware

import (
	"bubblr/handler"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// enableCORS is a middleware function to enable CORS headers.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set CORS headers
		origin := os.Getenv("CLIENT_ORIGIN")
		if len(origin) == 0 {
			panic("CLIENT_ORIGIN is not set in env!")
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Add("Access-Control-Allow-Origin", fmt.Sprintf("%s", origin))
		w.Header().Set("Access-Control-Request-Method", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// JWTMiddleware is a middleware function to verify JWT tokens.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Exclude the signup endpoint from JWT verification
		if r.URL.Path == "/" || r.URL.Path == "/signup" || r.URL.Path == "/login" || r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract the token from the Authorization header
		tokenString := extractTokenFromHeader(r.Header.Get("Authorization"))

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// TODO: Replace "your-secret-key" with your actual secret key
			return handler.SigningKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// extractTokenFromHeader extracts the token from the Authorization header.
// It expects the header value to be in the format "Bearer <token>".
func extractTokenFromHeader(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
