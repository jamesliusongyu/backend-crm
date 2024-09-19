package handler

import (
	database "backend-crm/internal/database/mongodb"
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

type MockVesselCollection struct {
	mock.Mock
}

func (m *MockVesselCollection) GetAll(ctx context.Context, tenant string) ([]database.VesselManagementResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.VesselManagementResponse), args.Error(1)
}

func (m *MockVesselCollection) Create(ctx context.Context, vessel database.Vessel) (string, error) {
	args := m.Called(ctx, vessel)
	return args.String(0), args.Error(1)
}

func (m *MockVesselCollection) GetByID(ctx context.Context, id string, tenant string) (database.VesselManagementResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.VesselManagementResponse), args.Error(1)
}

func (m *MockVesselCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.VesselManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.VesselManagementResponse), args.Error(1)
}

func (m *MockVesselCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.VesselManagementResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.VesselManagementResponse), args.Error(1)
}

func (m *MockVesselCollection) Update(ctx context.Context, id string, vessel database.Vessel) error {
	args := m.Called(ctx, id, vessel)
	return args.Error(0)
}

func (m *MockVesselCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestVesselHandler_GetAllVessels(t *testing.T) {
	mockCollection := new(MockVesselCollection)
	handler := NewVesselHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/vessel_management", handler.GetAllVessels)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.VesselManagementResponse{
		{
			Tenant: "customerA",
			VesselSpecifications: database.VesselSpecifications{
				ImoNumber:  123456789,
				VesselName: "Vessel Name",
				CallSign:   "CALLSIGN",
				SDWT:       50000,
				NRT:        30000,
				Flag:       "Flag",
				GRT:        60000,
				LOA:        200.0,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/vessel_management", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestVesselHandler_CreateVessel(t *testing.T) {
	mockCollection := new(MockVesselCollection)
	handler := NewVesselHandler(mockCollection)

	r := chi.NewRouter()
	r.Post("/vessel_management", handler.CreateVessel)

	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(vessel database.Vessel) bool {
		return vessel.VesselSpecifications.Flag == "Flag" && vessel.VesselSpecifications.VesselName == "Vessel Name"
	})).Return("id123", nil)

	req := httptest.NewRequest("POST", "/vessel_management", bytes.NewBufferString(`{
		"vessel_specifications": {
			"imo_number": 123456789,
			"vessel_name": "Vessel Name",
			"call_sign": "CALLSIGN",
			"sdwt": 50000,
			"nrt": 30000,
			"flag": "Flag",
			"grt": 60000,
			"loa": 200.3
		},
		"tenant":"customerA"
	}`))
	// token, _ := auth.GenerateJWT("james@customera.com", time.Now().Add(time.Hour*1))
	// req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestVesselHandler_GetVesselById(t *testing.T) {
	mockCollection := new(MockVesselCollection)
	handler := NewVesselHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/vessel_management/{vessel_id}", handler.GetVesselById)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.VesselManagementResponse{
		Tenant: "customerA",
		VesselSpecifications: database.VesselSpecifications{
			ImoNumber:  123456789,
			VesselName: "Vessel Name",
			CallSign:   "CALLSIGN",
			SDWT:       50000,
			NRT:        30000,
			Flag:       "Flag",
			GRT:        60000,
			LOA:        200.7,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/vessel_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("vessel_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestVesselHandler_UpdateVesselById(t *testing.T) {
	mockCollection := new(MockVesselCollection)
	handler := NewVesselHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/vessel_management/{vessel_id}", handler.UpdateVesselById)

	mockCollection.On("Update", mock.Anything, "1", mock.MatchedBy(func(vessel database.Vessel) bool {
		return vessel.Tenant == "customerA" && vessel.VesselSpecifications.VesselName == "Updated Vessel"
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/vessel_management/1", bytes.NewBufferString(`{
		"vessel_specifications": {
			"imo_number": 123456789,
			"vessel_name": "Updated Vessel",
			"call_sign": "UPDATEDCALLSIGN",
			"sdwt": 60000,
			"nrt": 40000,
			"flag": "Updated Flag",
			"grt": 70000,
			"loa": 250.5
		},
		"tenant":"customerA"
	}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("vessel_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestVesselHandler_DeleteVesselById(t *testing.T) {
	mockCollection := new(MockVesselCollection)
	handler := NewVesselHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/vessel_management/{vessel_id}", handler.DeleteVesselById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/vessel_management/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("vessel_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
