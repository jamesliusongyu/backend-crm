package handler

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
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

type MockAgentCollection struct {
	mock.Mock
}

func (m *MockAgentCollection) GetAll(ctx context.Context, tenant string) ([]database.AgentManagementResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.AgentManagementResponse), args.Error(1)
}

func (m *MockAgentCollection) Create(ctx context.Context, agent database.Agent) (string, error) {
	args := m.Called(ctx, agent)
	return args.String(0), args.Error(1)
}

func (m *MockAgentCollection) GetByID(ctx context.Context, id string, tenant string) (database.AgentManagementResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.AgentManagementResponse), args.Error(1)
}

func (m *MockAgentCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.AgentManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.AgentManagementResponse), args.Error(1)
}

func (m *MockAgentCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.AgentManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.AgentManagementResponse), args.Error(1)
}

func (m *MockAgentCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAgentCollection) Update(ctx context.Context, id string, agent database.Agent) error {
	args := m.Called(ctx, id, agent)
	return args.Error(0)
}

func TestAgentHandler_CreateAgent(t *testing.T) {
	mockCollection := new(MockAgentCollection)
	handler := NewAgentHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/agent_management", handler.CreateAgent)

	// Mock the Create method
	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(agent database.Agent) bool {
		return agent.Name == "John Doe" && agent.Email == "john.doe@example.com" && agent.Contact == "123-456-7890"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/agent_management", bytes.NewBufferString(`{"name":"John Doe", "email":"john.doe@example.com", "contact":"123-456-7890"}`))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))
	token, _ := auth.GenerateJWT("james@customera.com", time.Now().Add(time.Hour*1))
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestAgentHandler_GetAllAgents(t *testing.T) {
	mockCollection := new(MockAgentCollection)
	handler := NewAgentHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/agent_management", handler.GetAllAgents)

	// Mock the GetAll method
	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.AgentManagementResponse{
		{
			ID:        "1",
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			Contact:   "123-456-7890",
			Tenant:    "customerA",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/agent_management", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{}))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	mockCollection.AssertExpectations(t)
}

func TestAgentHandler_DeleteAgentById(t *testing.T) {
	mockCollection := new(MockAgentCollection)
	handler := NewAgentHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/agent_management/{agent_id}", handler.DeleteAgentById)

	// Mock the Delete method
	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/agent_management/1", nil)
	// This creates a new instance of chi.RouteContext, which is used by the chi router to hold route-specific information such as URL parameters (route variables).
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("agent_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestAgentHandler_UpdateAgentById(t *testing.T) {
	mockCollection := new(MockAgentCollection)
	handler := NewAgentHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/agent_management/{agent_id}", handler.UpdateAgentById)

	// Mock the Update method
	mockCollection.On("Update", mock.Anything, "1", mock.MatchedBy(func(agent database.Agent) bool {
		return agent.Name == "John Updated" && agent.Email == "john.updated@example.com" && agent.Contact == "987-654-3210"
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/agent_management/1", bytes.NewBufferString(`{"name":"John Updated","email":"john.updated@example.com","contact":"987-654-3210"}`))
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("agent_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestAgentHandler_GetAgentFromId(t *testing.T) {
	mockCollection := new(MockAgentCollection)
	handler := NewAgentHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/agent_management/{agent_id}", handler.GetAgentFromId)

	// Mock the GetByID method
	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.AgentManagementResponse{
		ID:        "1",
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Contact:   "123-456-7890",
		Tenant:    "customerA",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/agent_management/1", nil)
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("agent_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	mockCollection.AssertExpectations(t)
}
