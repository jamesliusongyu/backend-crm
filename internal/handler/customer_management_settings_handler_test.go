package handler

import (
	database "backend-crm/internal/database/mongodb"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/go-chi/chi/v5"
)

type MockCustomerCollection struct {
	mock.Mock
}

func (m *MockCustomerCollection) GetAll(ctx context.Context, tenant string) ([]database.CustomerResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.CustomerResponse), args.Error(1)
}

func (m *MockCustomerCollection) Create(ctx context.Context, createCustomerParams database.Customer) (string, error) {
	args := m.Called(ctx, createCustomerParams)
	return args.String(0), args.Error(1)
}

func (m *MockCustomerCollection) GetByID(ctx context.Context, id string, tenant string) (database.CustomerResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.CustomerResponse), args.Error(1)
}

func (m *MockCustomerCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.CustomerResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.CustomerResponse), args.Error(1)
}

func (m *MockCustomerCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.CustomerResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.CustomerResponse), args.Error(1)
}

func (m *MockCustomerCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerCollection) Update(ctx context.Context, id string, updateCustomerParams database.Customer) error {
	args := m.Called(ctx, id, updateCustomerParams)
	return args.Error(0)
}

func TestCustomerHandler_GetAllCustomers(t *testing.T) {
	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/customer_management", handler.GetAllCustomers)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.CustomerResponse{
		{
			ID:        "66403d4660a455ef0b8f88f6",
			Tenant:    "customerA",
			Customer:  "Vitol (Crude Oil)",
			Company:   "Vitol (Crude Oil)",
			Email:     "abc@gmail.com",
			Contact:   "+6512345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/customer_management", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Vitol (Crude Oil)")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestCustomerHandler_CreateCustomer(t *testing.T) {

	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/customer_management", handler.CreateCustomer)

	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(customer database.Customer) bool {
		return customer.Customer == "Vitol (Crude Oil)" && customer.Email == "vitol@customera.com"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/customer_management", bytes.NewBufferString(`{"customer":"Vitol (Crude Oil)", "tenant":"customerA", "company":"Vitol (Crude Oil)","email":"vitol@customera.com", "contact":"+6512345678"}`))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestCustomerHandler_GetCustomerFromId(t *testing.T) {
	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/customer_management/{customer_id}", handler.GetCustomerFromId)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.CustomerResponse{
		ID:        "66403d4660a455ef0b8f88f6",
		Tenant:    "customerA",
		Customer:  "Vitol (Crude Oil)",
		Company:   "Vitol (Crude Oil)",
		Email:     "vitol@gmail.com",
		Contact:   "+6512345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/customer_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("customer_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Vitol (Crude Oil)")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestCustomerHandler_GetCustomerFromName(t *testing.T) {
	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/customer_management/{customer_name}", handler.GetCustomerFromName)

	mockCollection.On("GetByKeyValue", mock.Anything, "customer", "Vitol", mock.Anything).Return(database.CustomerResponse{
		ID:        "66403d4660a455ef0b8f88f6",
		Tenant:    "customerA",
		Customer:  "Vitol",
		Company:   "Vitol",
		Email:     "vitol@gmail.com",
		Contact:   "+6512345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/customer_management/Vitol", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("customer_name", "Vitol")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Vitol")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestCustomerHandler_UpdateCustomerById(t *testing.T) {
	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/customer_management/{customer_id}", handler.UpdateCustomerById)

	mockCollection.On("Update", mock.Anything, "1", database.Customer{
		Tenant:   "customerA",
		Customer: "Updated Customer",
		Company:  "Updated Company",
		Email:    "Updated Email",
		Contact:  "+6512345678",
	}).Return(nil)

	req := httptest.NewRequest("PUT", "/customer_management/1", bytes.NewBufferString(`{"customer":"Updated Customer", "tenant":"customerA", "company":"Updated Company","email":"Updated Email", "contact":"+6512345678"}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("customer_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestCustomerHandler_DeleteCustomerById(t *testing.T) {
	mockCollection := new(MockCustomerCollection)
	handler := NewCustomerHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/customer_management/{customer_id}", handler.DeleteCustomerById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/customer_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("customer_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
