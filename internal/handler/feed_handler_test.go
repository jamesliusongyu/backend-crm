package handler

import (
	database "backend-crm/internal/database/mongodb"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFeedMessageCollection struct {
	mock.Mock
}

func (m *MockFeedMessageCollection) GetAll(ctx context.Context, tenant string) ([]database.FeedEmailResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.FeedEmailResponse), args.Error(1)
}

func (m *MockFeedMessageCollection) Create(ctx context.Context, email database.FeedEmail) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *MockFeedMessageCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.FeedEmailResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.FeedEmailResponse), args.Error(1)
}

func (m *MockFeedMessageCollection) GetByID(ctx context.Context, id string, tenant string) (database.FeedEmailResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.FeedEmailResponse), args.Error(1)
}

func (m *MockFeedMessageCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.FeedEmailResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.FeedEmailResponse), args.Error(1)
}

func (m *MockFeedMessageCollection) Update(ctx context.Context, id string, email database.FeedEmail) error {
	args := m.Called(ctx, id, email)
	return args.Error(0)
}

func (m *MockFeedMessageCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFeedHandler_GetAllFeedMessages(t *testing.T) {
	mockFeedCollection := new(MockFeedMessageCollection)
	mockShipmentCollection := new(MockShipmentCollection)
	mockChecklistCollection := new(MockChecklistCollection)
	handler := NewFeedHandler(mockFeedCollection, mockShipmentCollection, mockChecklistCollection)

	r := chi.NewRouter()
	r.Get("/feed", handler.GetAllFeedEmails)

	mockFeedCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.FeedEmailResponse{
		{
			ID:               "1234",
			Tenant:           "customerA",
			MasterEmail:      "yoke@captain.com",
			ReceivedDateTime: time.Now(),
			ToEmailAddress:   "yoke@customera.com",
			Subject:          "Hello world",
			BodyContent:      "This is a test.",
			ShipmentId:       "abc123",
		},
	}, nil)

	req := httptest.NewRequest("GET", "/feed", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "customerA")
	mockFeedCollection.AssertExpectations(t)
}

// func TestFeedHandler_CreateFeedMessage(t *testing.T) {
// 	mockFeedCollection := new(MockFeedMessageCollection)
// 	mockShipmentCollection := new(MockShipmentCollection)
// 	handler := NewFeedHandler(mockFeedCollection, mockShipmentCollection)

// 	r := chi.NewRouter()
// 	r.Post("/master_email_messages/", handler.CreateFeedMessage)

// 	// Setup mock for the Create method on mockFeedCollection
// 	mockFeedCollection.On("Create", mock.Anything, mock.MatchedBy(func(email database.FeedEmail) bool {
// 		return email.Subject == "TEST"
// 	})).Return("id123", nil)

// 	// Setup mock for GetAll method on mockShipmentCollection, fake a shipment id so it doesnt skip creation
// 	mockShipmentCollection.On("GetAll", mock.Anything, "customerA").Return([]database.ShipmentResponse{}, nil)

// 	mockShipmentCollection.On("GetByID", mock.Anything, "abc123", "customerA").Return(database.ShipmentResponse{}, nil)

// 	req := httptest.NewRequest("POST", "/master_email_messages/", bytes.NewBufferString(`{
// 		"received_date_time": "2022-08-01T10:24:17.000Z",
// 		"subject": "TEST",
// 		"body_content": "TESTING",
// 		"from_email_address": "yoke@captain.com",
// 		"to_email_address": "yoke@customera.com",
// 		"shipment_id": "abc123"
// 	}`))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)

// 	// Check the response status code
// 	assert.Equal(t, http.StatusCreated, w.Code)

// 	// Verify all expectations
// 	mockFeedCollection.AssertExpectations(t)
// 	mockShipmentCollection.AssertExpectations(t)
// }

func TestFeedHandler_GetAllFeedMessageByKeyValue(t *testing.T) {
	mockFeedCollection := new(MockFeedMessageCollection)
	mockShipmentCollection := new(MockShipmentCollection)
	handler := NewFeedHandler(mockFeedCollection, mockShipmentCollection)

	r := chi.NewRouter()
	r.Get("/feed/{shipment_id}", handler.GetFeedEmailsByShipmentId)

	mockFeedCollection.On("GetAllByKeyValue", mock.Anything, "shipmentid", "abc123", mock.Anything).Return([]database.FeedEmailResponse{
		{
			ID:               "1234",
			Tenant:           "customerA",
			MasterEmail:      "yoke@captain.com",
			ReceivedDateTime: time.Now(),
			ToEmailAddress:   "yoke@customera.com",
			Subject:          "Hello world",
			BodyContent:      "This is a test.",
			ShipmentId:       "abc123",
		},
	}, nil)

	req := httptest.NewRequest("GET", "/feed/abc123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shipment_id", "abc123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// assert.Contains(t, w.Body.String(), "customerA")
	mockFeedCollection.AssertExpectations(t)
}
