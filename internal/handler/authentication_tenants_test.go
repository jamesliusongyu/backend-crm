package handler

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/chi/v5"
)

type MockLoginCollection struct {
	mock.Mock
}

func (m *MockLoginCollection) GetAll(ctx context.Context, tenant string) ([]database.TenantUserResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.TenantUserResponse), args.Error(1)
}

func (m *MockLoginCollection) Create(ctx context.Context, createTenantUserParams database.TenantUser) (string, error) {
	args := m.Called(ctx, createTenantUserParams)
	return args.String(0), args.Error(1)
}

func (m *MockLoginCollection) GetByID(ctx context.Context, id string, tenant string) (database.TenantUserResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.TenantUserResponse), args.Error(1)
}

func (m *MockLoginCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.TenantUserResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.TenantUserResponse), args.Error(1)
}

func (m *MockLoginCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.TenantUserResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.TenantUserResponse), args.Error(1)
}

func (m *MockLoginCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLoginCollection) Update(ctx context.Context, id string, updateTenantUserParams database.TenantUser) error {
	args := m.Called(ctx, id, updateTenantUserParams)
	return args.Error(0)
}

type MockSessionCollection struct {
	mock.Mock
}

func (m *MockSessionCollection) Create(ctx context.Context, createSessionParams database.Session) (string, error) {
	args := m.Called(ctx, createSessionParams)
	return args.String(0), args.Error(1)
}

func (m *MockSessionCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.SessionResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.SessionResponse), args.Error(1)
}

func (m *MockSessionCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionCollection) GetAll(ctx context.Context, tenant string) ([]database.SessionResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.SessionResponse), args.Error(1)
}

func (m *MockSessionCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.SessionResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.SessionResponse), args.Error(1)
}

func (m *MockSessionCollection) GetByID(ctx context.Context, id string, tenant string) (database.SessionResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.SessionResponse), args.Error(1)
}

func (m *MockSessionCollection) Update(ctx context.Context, id string, createSessionParams database.Session) error {
	args := m.Called(ctx, id, createSessionParams)
	return args.Error(0)
}

func TestHandler_CreateTenantUser(t *testing.T) {
	mockLoginCollection := new(MockLoginCollection)
	mockSessionCollection := new(MockSessionCollection)

	handler := NewLoginHandler(mockLoginCollection, mockSessionCollection)

	r := chi.NewRouter()
	r.Post("/tenant_user", handler.CreateTenantUser)

	// Mock the Create method
	mockLoginCollection.On("Create", mock.Anything, mock.MatchedBy(func(user database.TenantUser) bool {
		hashedPassword, _ := auth.HashPassword("password123")
		return user.Email == "test@example.com" && auth.CheckPasswordHash("password123", hashedPassword)
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/tenant_user", bytes.NewBufferString(`{"email":"test@example.com", "password":"password123"}`))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockLoginCollection.AssertExpectations(t)
}

func TestHandler_LoginAndAuth(t *testing.T) {
	mockLoginCollection := new(MockLoginCollection)
	mockSessionCollection := new(MockSessionCollection)

	handler := NewLoginHandler(mockLoginCollection, mockSessionCollection)

	r := chi.NewRouter()
	r.Post("/login", handler.LoginAndAuth)

	// Mock the GetByKeyValue method
	hashedPassword, _ := auth.HashPassword("password123")
	mockLoginCollection.On("GetByKeyValue", mock.Anything, "email", "test@example.com", mock.Anything).Return(database.TenantUserResponse{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}, nil)

	mockSessionCollection.On("Create", mock.Anything, mock.Anything).Return("id123", nil)

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"email":"test@example.com", "password":"password123"}`))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// Check if the response contains "jwtToken" and "sessionID"
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "jwtToken")
	assert.Contains(t, response, "sessionID")

	// Check if the cookie contains the JWT token
	cookie := w.Result().Cookies()
	assert.Len(t, cookie, 1)
	assert.Equal(t, "jwtToken", cookie[0].Name)
	assert.Equal(t, response["jwtToken"], cookie[0].Value)

	// Ensure mocks were called
	mockLoginCollection.AssertExpectations(t)
	mockSessionCollection.AssertExpectations(t)
}

func TestHandler_ValidateCookie(t *testing.T) {
	mockSessionCollection := new(MockSessionCollection)

	handler := &LoginHandler{
		SessionCollection: mockSessionCollection,
	}

	// Set up the router with the middleware and a protected route
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(handler.ValidateCookie)
		r.Get("/shipments", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Shipments content"))
		})
	})

	// Step 1: Prepare a valid JWT token
	email := "test@example.com"
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	// Step 2: Mock the session collection to return a valid session
	sessionResponse := database.SessionResponse{
		SessionID: "some-session-id",
		Email:     email,
		JWTToken:  tokenString,
		Tenant:    "tenant123",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	mockSessionCollection.On("GetByKeyValue", mock.Anything, "jwt_token", tokenString, mock.Anything).Return(sessionResponse, nil)

	// Step 3: Create a request with the JWT token cookie to a protected route
	req := httptest.NewRequest("GET", "/shipments", nil)
	req.AddCookie(&http.Cookie{
		Name:     "jwtToken",
		Value:    tokenString,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAge, // Expiration time of the cookie
		HttpOnly: true,         // Prevents client-side scripts from accessing the cookie
		Secure:   true,         // Ensures the cookie is only sent over HTTPS
	})

	// Step 4: Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Shipments content", w.Body.String())

	mockSessionCollection.AssertExpectations(t)
}

func TestHandler_ValidateCookie_ExpiredToken(t *testing.T) {
	mockSessionCollection := new(MockSessionCollection)

	handler := &LoginHandler{
		SessionCollection: mockSessionCollection,
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(handler.ValidateCookie)
		r.Get("/shipments", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Shipments content"))
		})
	})

	// Prepare an expired JWT token
	email := "test@example.com"
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
		},
	}

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, _ := expiredToken.SignedString(jwtKey)

	// Mock the session collection to return a valid session
	sessionResponse := database.SessionResponse{
		SessionID: "some-session-id",
		Email:     email,
		JWTToken:  expiredTokenString,
		Tenant:    "tenant123",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Mock the refreshed session
	mockSessionCollection.On("GetByKeyValue", mock.Anything, "jwt_token", expiredTokenString, mock.Anything).Return(sessionResponse, nil)

	mockSessionCollection.On("Update", mock.Anything, sessionResponse.ID, mock.MatchedBy(func(session database.Session) bool {
		return session.SessionID == sessionResponse.SessionID &&
			session.Email == sessionResponse.Email &&
			session.JWTToken != expiredTokenString &&
			session.Tenant == sessionResponse.Tenant
	})).Return(nil)

	// Create a request with the expired JWT token cookie
	req := httptest.NewRequest("GET", "/shipments", nil)
	req.AddCookie(&http.Cookie{
		Name:     "jwtToken",
		Value:    expiredTokenString,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieMaxAge, // Expiration time of the cookie
		HttpOnly: true,         // Prevents client-side scripts from accessing the cookie
		Secure:   true,         // Ensures the cookie is only sent over HTTPS
	})

	// Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Shipments content", w.Body.String())

	mockSessionCollection.AssertExpectations(t)
}

func TestHandler_ValidateCookie_SessionNotFound(t *testing.T) {
	mockSessionCollection := new(MockSessionCollection)

	handler := &LoginHandler{
		SessionCollection: mockSessionCollection,
	}

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(handler.ValidateCookie)
		r.Get("/shipments", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Shipments content"))
		})
	})

	// Step 1: Prepare a valid JWT token
	email := "test@example.com"
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	// Step 2: Mock the session collection to return a session not found error
	mockSessionCollection.On("GetByKeyValue", mock.Anything, "jwt_token", tokenString, mock.Anything).Return(database.SessionResponse{}, mongo.ErrNoDocuments)

	// Step 3: Create a request with the JWT token cookie
	req := httptest.NewRequest("GET", "/shipments", nil)
	req.AddCookie(&http.Cookie{
		Name:  "jwtToken",
		Value: tokenString,
	})

	// Step 4: Make the request and check the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized: Session not found")

	mockSessionCollection.AssertExpectations(t)
}
