package handler

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/core"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTerminalCollection struct {
	mock.Mock
}

func (m *MockTerminalCollection) GetAll(ctx context.Context, tenant string) ([]database.TerminalManagementResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.TerminalManagementResponse), args.Error(1)
}

func (m *MockTerminalCollection) Create(ctx context.Context, terminal database.Terminal) (string, error) {
	args := m.Called(ctx, terminal)
	return args.String(0), args.Error(1)
}

func (m *MockTerminalCollection) GetByID(ctx context.Context, id string, tenant string) (database.TerminalManagementResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.TerminalManagementResponse), args.Error(1)
}

func (m *MockTerminalCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.TerminalManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.TerminalManagementResponse), args.Error(1)
}

func (m *MockTerminalCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.TerminalManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.TerminalManagementResponse), args.Error(1)
}

func (m *MockTerminalCollection) Update(ctx context.Context, id string, terminal database.Terminal) error {
	args := m.Called(ctx, id, terminal)
	return args.Error(0)
}

func (m *MockTerminalCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTerminalHandler_GetAllTerminals(t *testing.T) {
	mockCollection := new(MockTerminalCollection)
	handler := NewTerminalHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/terminal_management", handler.GetAllTerminals)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.TerminalManagementResponse{
		{
			Tenant: "customerA",
			TerminalSpecifications: core.TerminalSpecifications{
				Name:    "Terminal Name",
				Address: "Terminal Address",
				Email:   "terminal@example.com",
				Contact: "1234567890",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/terminal_management", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestTerminalHandler_CreateTerminal(t *testing.T) {
	mockCollection := new(MockTerminalCollection)
	handler := NewTerminalHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/terminal_management", handler.CreateTerminal)

	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(terminal database.Terminal) bool {
		return terminal.TerminalSpecifications.Name == "Terminal Name" && terminal.TerminalSpecifications.Address == "Terminal Address"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/terminal_management", bytes.NewBufferString(`{
		"terminal_specifications": {
			"name": "Terminal Name",
			"address": "Terminal Address",
			"email": "terminal@example.com",
			"contact": "1234567890"
		},
		"tenant":"customerA"
	}`))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestTerminalHandler_GetTerminalById(t *testing.T) {
	mockCollection := new(MockTerminalCollection)
	handler := NewTerminalHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/terminal_management/{terminal_id}", handler.GetTerminalById)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.TerminalManagementResponse{
		Tenant: "customerA",
		TerminalSpecifications: core.TerminalSpecifications{
			Name:    "Terminal Name",
			Address: "Terminal Address",
			Email:   "terminal@example.com",
			Contact: "1234567890",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/terminal_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("terminal_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestTerminalHandler_UpdateTerminalById(t *testing.T) {
	mockCollection := new(MockTerminalCollection)
	handler := NewTerminalHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/terminal_management/{terminal_id}", handler.UpdateTerminalById)

	mockCollection.On("Update", mock.Anything, "1", mock.MatchedBy(func(terminal database.Terminal) bool {
		return terminal.Tenant == "customerA" && terminal.TerminalSpecifications.Name == "Updated Terminal"
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/terminal_management/1", bytes.NewBufferString(`{
		"terminal_specifications": {
			"name": "Updated Terminal",
			"address": "Updated Address",
			"email": "updated_terminal@example.com",
			"contact": "0987654321"
		},
		"tenant":"customerA"
	}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("terminal_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestTerminalHandler_DeleteTerminalById(t *testing.T) {
	mockCollection := new(MockTerminalCollection)
	handler := NewTerminalHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/terminal_management/{terminal_id}", handler.DeleteTerminalById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/terminal_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("terminal_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
