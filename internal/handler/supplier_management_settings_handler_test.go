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

type MockSupplierCollection struct {
	mock.Mock
}

func (m *MockSupplierCollection) GetAll(ctx context.Context, tenant string) ([]database.SupplierManagementResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.SupplierManagementResponse), args.Error(1)
}

func (m *MockSupplierCollection) Create(ctx context.Context, supplier database.Supplier) (string, error) {
	args := m.Called(ctx, supplier)
	return args.String(0), args.Error(1)
}

func (m *MockSupplierCollection) GetByID(ctx context.Context, id string, tenant string) (database.SupplierManagementResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.SupplierManagementResponse), args.Error(1)
}

func (m *MockSupplierCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.SupplierManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.SupplierManagementResponse), args.Error(1)
}

func (m *MockSupplierCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.SupplierManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.SupplierManagementResponse), args.Error(1)
}

func (m *MockSupplierCollection) Update(ctx context.Context, id string, supplier database.Supplier) error {
	args := m.Called(ctx, id, supplier)
	return args.Error(0)
}

func (m *MockSupplierCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSupplierHandler_GetAllSuppliers(t *testing.T) {
	mockCollection := new(MockSupplierCollection)
	handler := NewSupplierHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/supplier_management", handler.GetAllSuppliers)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.SupplierManagementResponse{
		{
			Tenant: "customerA",
			SupplierSpecifications: core.SupplierSpecifications{
				Name:    "Supplier Name",
				Vessels: []string{"Supplier Vessel"},
				Email:   "supplier@example.com",
				Contact: "1234567890",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/supplier_management", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestSupplierHandler_CreateSupplier(t *testing.T) {
	mockCollection := new(MockSupplierCollection)
	handler := NewSupplierHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/supplier_management", handler.CreateSupplier)

	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(supplier database.Supplier) bool {
		return supplier.SupplierSpecifications.Name == "Supplier Name" && supplier.SupplierSpecifications.Vessels[0] == "Supplier Vessel"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/supplier_management", bytes.NewBufferString(`{
		"supplier_specifications": {
			"name": "Supplier Name",
			"vessels": ["Supplier Vessel"],
			"email": "supplier@example.com",
			"contact": "1234567890"
		},
		"tenant":"customerA"
	}`))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestSupplierHandler_GetSupplierById(t *testing.T) {
	mockCollection := new(MockSupplierCollection)
	handler := NewSupplierHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/supplier_management/{supplier_id}", handler.GetSupplierById)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.SupplierManagementResponse{
		Tenant: "customerA",
		SupplierSpecifications: core.SupplierSpecifications{
			Name:    "Supplier Name",
			Vessels: []string{"Supplier Vessel"},
			Email:   "supplier@example.com",
			Contact: "1234567890",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/supplier_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("supplier_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestSupplierHandler_UpdateSupplierById(t *testing.T) {
	mockCollection := new(MockSupplierCollection)
	handler := NewSupplierHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/supplier_management/{supplier_id}", handler.UpdateSupplierById)

	mockCollection.On("Update", mock.Anything, "1", mock.MatchedBy(func(supplier database.Supplier) bool {
		return supplier.Tenant == "customerA" && supplier.SupplierSpecifications.Name == "Updated Supplier"
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/supplier_management/1", bytes.NewBufferString(`{
		"supplier_specifications": {
			"name": "Updated Supplier",
			"vessel": "Updated Vessel",
			"email": "updated_supplier@example.com",
			"contact": "0987654321"
		},
		"tenant":"customerA"
	}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("supplier_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestSupplierHandler_DeleteSupplierById(t *testing.T) {
	mockCollection := new(MockSupplierCollection)
	handler := NewSupplierHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/supplier_management/{supplier_id}", handler.DeleteSupplierById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/supplier_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("supplier_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
