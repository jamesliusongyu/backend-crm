package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var mySigningKey = []byte("jameslsy") // Replace with your secret key

// Claims struct
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Custom context key types to avoid collisions
type contextKey string

const (
	ContextKeyEmail   contextKey = "email"
	ContextKeySession contextKey = "session"
)

// func JWTMiddleware(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
// 			return
// 		}

// 		bearerToken := strings.Split(authHeader, " ")
// 		if len(bearerToken) != 2 {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, http.ErrAbortHandler
// 			}
// 			return mySigningKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		if !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Log or use the email information from the claims
// 		// Here we store the email information in the request context
// 		ctx := context.WithValue(r.Context(), ContextKeyEmail, claims.Email)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// Function to get the email from the token context
func GetEmailFromToken(ctx context.Context) string {
	email, _ := ctx.Value(ContextKeyEmail).(string)
	return email
}

// Function to get the email from the token context
func GetSessionInfoFromToken(ctx context.Context) string {
	session, _ := ctx.Value(ContextKeySession).(string)
	return session
}

// GenerateJWT generates a JWT token with a given email and an expiration time
func GenerateJWT(email string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
