package handler

import (
	database "backend-crm/internal/database/mongodb"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryManagementActivityTypeCollection struct {
	mock.Mock
}

func (m *MockCategoryManagementActivityTypeCollection) GetAll(ctx context.Context, tenant string) ([]database.CategoryManagementActivityTypeResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.CategoryManagementActivityTypeResponse), args.Error(1)
}

func (m *MockCategoryManagementActivityTypeCollection) Create(ctx context.Context, createParams database.CategoryManagementActivityType) (string, error) {
	args := m.Called(ctx, createParams)
	return args.String(0), args.Error(1)
}

func (m *MockCategoryManagementActivityTypeCollection) GetByID(ctx context.Context, id string, tenant string) (database.CategoryManagementActivityTypeResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.CategoryManagementActivityTypeResponse), args.Error(1)
}

func (m *MockCategoryManagementActivityTypeCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.CategoryManagementActivityTypeResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.CategoryManagementActivityTypeResponse), args.Error(1)
}

func (m *MockCategoryManagementActivityTypeCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.CategoryManagementActivityTypeResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.CategoryManagementActivityTypeResponse), args.Error(1)
}

func (m *MockCategoryManagementActivityTypeCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCategoryManagementActivityTypeCollection) Update(ctx context.Context, id string, updateParams database.CategoryManagementActivityType) error {
	args := m.Called(ctx, id, updateParams)
	return args.Error(0)
}

func TestCategoryManagementActivityTypeHandler_GetAllCategoryManagementActivityTypes(t *testing.T) {
	mockCollection := new(MockCategoryManagementActivityTypeCollection)
	handler := NewActivityTypeHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/category_management/activity_type", handler.GetAllCategoryManagementActivityTypes)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.CategoryManagementActivityTypeResponse{
		{
			ID:           "6666784b67a23198be669d5d",
			Tenant:       "customerA",
			ActivityType: "Oil Transport",
		},
	}, nil)

	req := httptest.NewRequest("GET", "/category_management/activity_type", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Oil Transport")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestCategoryManagementActivityTypeHandler_CreateActivityType(t *testing.T) {
	mockCollection := new(MockCategoryManagementActivityTypeCollection)
	handler := NewActivityTypeHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/category_management/activity_type", handler.CreateActivityType)

	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(activityType database.CategoryManagementActivityType) bool {
		return activityType.ActivityType == "Oil Transport"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/category_management/activity_type", bytes.NewBufferString(`{"activity_type":"Oil Transport"}`))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestCategoryManagementActivityTypeHandler_GetActivityTypeFromId(t *testing.T) {
	mockCollection := new(MockCategoryManagementActivityTypeCollection)
	handler := NewActivityTypeHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/category_management/activity_type/{activity_type_id}", handler.GetActivityTypeFromId)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.CategoryManagementActivityTypeResponse{
		ID:           "6666784b67a23198be669d5d",
		Tenant:       "customerA",
		ActivityType: "Oil Transport",
	}, nil)

	req := httptest.NewRequest("GET", "/category_management/activity_type/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("activity_type_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Oil Transport")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestCategoryManagementActivityTypeHandler_DeleteActivityTypeById(t *testing.T) {
	mockCollection := new(MockCategoryManagementActivityTypeCollection)
	handler := NewActivityTypeHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/category_management/activity_type/{activity_type_id}", handler.DeleteActivityTypeById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/category_management/activity_type/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("activity_type_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
