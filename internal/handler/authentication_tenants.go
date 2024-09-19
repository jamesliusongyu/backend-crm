package handler

import (
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
	"context"

	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("jameslsy")
var cookieMaxAge = 86400

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type LoginHandler struct {
	LoginCollection   database.Collection[database.TenantUser, database.TenantUserResponse]
	SessionCollection database.Collection[database.Session, database.SessionResponse]
}

func NewLoginHandler(
	loginCollection database.Collection[database.TenantUser, database.TenantUserResponse],
	sessionCollection database.Collection[database.Session, database.SessionResponse],
) *LoginHandler {
	return &LoginHandler{
		LoginCollection:   loginCollection,
		SessionCollection: sessionCollection, // Add this line
	}
}

// StatusCheck handles the status check request
func (h *LoginHandler) StatusCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]int{"status_check": 200}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the status code
	w.WriteHeader(http.StatusOK)
	// Write the JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *LoginHandler) CreateTenantUser(w http.ResponseWriter, r *http.Request) {
	var createTenantUserParams database.TenantUser
	log.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&createTenantUserParams)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	tenant := config.MakeMapping(createTenantUserParams.Email)

	// Hash the password
	hashedPassword, err := auth.HashPassword(createTenantUserParams.Password)
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}
	createTenantUserParams.Password = hashedPassword
	createTenantUserParams.Tenant = tenant
	log.Println(createTenantUserParams, "as")
	_, err = h.LoginCollection.Create(r.Context(), createTenantUserParams)

	if err != nil {
		if errors.Is(err, database.ErrDuplicateKey) {
			render.Render(w, r, ErrDuplicate(err))
			return
		}
		render.Render(w, r, ErrInternalServerError)
		return
	}

	w.WriteHeader(201)
	w.Write(nil)

	render.Render(w, r, SuccessCreated)
}

// func (h *LoginHandler) ValidateCookie(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Retrieve the cookie from the request
// 		cookie, err := r.Cookie("jwtToken")
// 		if err != nil {
// 			http.Error(w, "Unauthorized: No JWT token found in the request", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := cookie.Value

// 		// Parse the JWT token
// 		claims := &Claims{}
// 		token, invalid := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, http.ErrAbortHandler
// 			}
// 			return jwtKey, nil
// 		})
// 		log.Println(invalid, "invalid")
// 		log.Println(token, "token")
// 		log.Println(token.Valid, "token")

// 		if invalid != nil || !token.Valid {
// 			http.Error(w, "Unauthorized: Invalid or expired JWT token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Extract the email from the claims
// 		email := claims.Email
// 		if email == "" {
// 			http.Error(w, "Unauthorized: Invalid JWT token, email missing", http.StatusUnauthorized)
// 			return
// 		}

// 		// Map the tenant based on the email
// 		tenant := config.MakeMapping(email)

// 		// Check if the session exists in the MongoDB session collection
// 		sessionResponse, err := h.SessionCollection.GetByKeyValue(r.Context(), "jwt_token", tokenString, tenant)
// 		log.Println(sessionResponse, "sessreponse")
// 		if err != nil {
// 			if err == mongo.ErrNoDocuments {
// 				http.Error(w, "Unauthorized: Session not found", http.StatusUnauthorized)
// 				return
// 			}
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		// Check if the session has expired
// 		if time.Now().After(sessionResponse.ExpiresAt) {
// 			http.Error(w, "Unauthorized: Session expired", http.StatusUnauthorized)
// 			return
// 		}

// 		if invalid != nil || !token.Valid {
// 			http.Error(w, "Unauthorized: Invalid or expired JWT token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Check if the token is close to expiry (e.g., within 5 minutes)

// 		if time.Until(claims.ExpiresAt.Time) < 8*time.Minute {
// 			log.Println("fresh token")
// 			// Generate a new JWT token with extended expiration time
// 			newExpirationTime := time.Now().Add(10 * time.Minute)
// 			claims.ExpiresAt = jwt.NewNumericDate(newExpirationTime)
// 			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 			newTokenString, err := newToken.SignedString(jwtKey)
// 			if err != nil {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}

// 			// Set the new JWT token in the cookie
// 			http.SetCookie(w, &http.Cookie{
// 				Name:     "jwtToken",
// 				Value:    newTokenString,
// 				Path:     "/",
// 				SameSite: http.SameSiteStrictMode,
// 				MaxAge:   86400, // Expiration time of the cookie
// 				HttpOnly: true,
// 				Secure:   true,
// 			})

// 			// Create and update the new session with the new JWT token
// 			newSession := database.Session{
// 				SessionID: sessionResponse.SessionID,
// 				Email:     sessionResponse.Email,
// 				JWTToken:  newTokenString,
// 				Tenant:    sessionResponse.Tenant,
// 				CreatedAt: time.Now(),
// 				ExpiresAt: newExpirationTime,
// 			}

// 			err = h.SessionCollection.Update(r.Context(), sessionResponse.ID, newSession)
// 			log.Println("here")
// 			if err != nil {
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}
// 			ctx := context.WithValue(r.Context(), auth.ContextKeySession, newSession.JWTToken)
// 			ctx = context.WithValue(ctx, auth.ContextKeyEmail, email)
// 			next.ServeHTTP(w, r.WithContext(ctx))

// 		} else {
// 			ctx := context.WithValue(r.Context(), auth.ContextKeySession, sessionResponse.JWTToken)
// 			ctx = context.WithValue(ctx, auth.ContextKeyEmail, email)
// 			next.ServeHTTP(w, r.WithContext(ctx))

// 		}
// 	})
// }

func (h *LoginHandler) LoginAndAuth(w http.ResponseWriter, r *http.Request) {

	var tenantUserCredentials database.TenantUser
	err := json.NewDecoder(r.Body).Decode(&tenantUserCredentials)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	log.Println(r.Body, "sddd")
	tenant := config.MakeMapping(tenantUserCredentials.Email)

	tenantUserResponse, err := h.LoginCollection.GetByKeyValue(r.Context(), "email", tenantUserCredentials.Email, tenant)
	log.Println(tenantUserResponse, "hi")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !auth.CheckPasswordHash(tenantUserCredentials.Password, tenantUserResponse.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Step 2: Generate JWT token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: tenantUserCredentials.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",  // The name of the cookie
		Value:    tokenString, // The JWT token as the cookie's value
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		MaxAge:   cookieMaxAge, // Expiration time of the cookie
		HttpOnly: true,         // Prevents client-side scripts from accessing the cookie
		Secure:   true,         // Ensures the cookie is only sent over HTTPS
	})

	// Step 3: Generate sessionID
	sessionID := uuid.NewString()

	// Store the sessionID in MongoDB session collection
	session := database.Session{
		SessionID: sessionID,
		Email:     tenantUserCredentials.Email,
		JWTToken:  tokenString,
		Tenant:    tenant,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Set session expiration time
	}
	_, err = h.SessionCollection.Create(r.Context(), session)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Step 4: Return both JWT token and sessionID
	response := map[string]string{
		"jwtToken":  tokenString,
		"sessionID": sessionID,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *LoginHandler) ValidateCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the cookie from the request
		cookie, err := r.Cookie("jwtToken")
		if err != nil {
			http.Error(w, "Unauthorized: No JWT token found in the request", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Parse the JWT token
		claims := &Claims{}
		_, invalid := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return jwtKey, nil
		})

		email := claims.Email
		if email == "" {
			http.Error(w, "Unauthorized: Invalid JWT token, email missing", http.StatusUnauthorized)
			return
		}

		// Map the tenant based on the email
		tenant := config.MakeMapping(email)

		// Check if the session exists in the MongoDB session collection
		sessionResponse, err := h.SessionCollection.GetByKeyValue(r.Context(), "jwt_token", tokenString, tenant)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Unauthorized: Session not found", http.StatusUnauthorized)
				log.Println("unauthorized, session not found")
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Internal server error")
			return
		}

		if time.Now().After(sessionResponse.ExpiresAt) {
			http.Error(w, "Unauthorized: Session expired", http.StatusUnauthorized)
			log.Println("unauthorized, session expired, deleting old session")
			log.Println(sessionResponse)
			deleteSuccess := h.SessionCollection.Delete(r.Context(), sessionResponse.ID)
			if deleteSuccess != nil {
				log.Println("Error deleting session")
			}
			return
		}

		// Handle different cases based on the parsing result
		switch {
		case invalid != nil && errors.Is(invalid, jwt.ErrTokenExpired):
			// If the token is expired, refresh the JWT token and session
			log.Println("Token expired.. Refreshing session")

			h.refreshJWTAndSession(w, r, claims, sessionResponse, next)

			return
		case invalid != nil:
			// If there's any other error, return an internal server error
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Println("Token expired.. Internal server error")

			return
		}

		// If the JWT token is close to expiry, refresh it
		if time.Until(claims.ExpiresAt.Time) < 5*time.Minute {
			// Refresh JWT token and extend session duration
			h.refreshJWTAndSession(w, r, claims, sessionResponse, next)
		} else {
			// JWT token and session are valid, proceed with the request
			ctx := context.WithValue(r.Context(), auth.ContextKeySession, sessionResponse.JWTToken)
			ctx = context.WithValue(ctx, auth.ContextKeyEmail, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (h *LoginHandler) refreshJWTAndSession(w http.ResponseWriter, r *http.Request, claims *Claims, sessionResponse database.SessionResponse, next http.Handler) {
	// Generate a new JWT token with extended expiration time
	newExpirationTime := time.Now().Add(30 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(newExpirationTime)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the new JWT token in the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    newTokenString,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		MaxAge:   cookieMaxAge, // Expiration time of the cookie
		HttpOnly: true,
		Secure:   true,
	})

	// Create and update the new session with the new JWT token
	newSession := database.Session{
		SessionID: sessionResponse.SessionID,
		Email:     sessionResponse.Email,
		JWTToken:  newTokenString,
		Tenant:    sessionResponse.Tenant,
		CreatedAt: time.Now(),
		ExpiresAt: sessionResponse.ExpiresAt,
	}

	err = h.SessionCollection.Update(r.Context(), sessionResponse.ID, newSession)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("refreshed token")

	// Proceed with the request
	ctx := context.WithValue(r.Context(), auth.ContextKeySession, newSession.JWTToken)
	ctx = context.WithValue(ctx, auth.ContextKeyEmail, claims.Email)
	next.ServeHTTP(w, r.WithContext(ctx))
}
