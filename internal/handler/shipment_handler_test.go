package handler

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/core"
	"backend-crm/pkg/enum"
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

type MockWhatsAppClient struct {
	mock.Mock
}

func (m *MockWhatsAppClient) CreateShipmentWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime, agentContact string) (string, error) {
	args := m.Called(agentName, vesselName, imoNumber, etaDate, etaTime, agentContact)
	return args.String(0), args.Error(1)
}

type MockShipmentCollection struct {
	mock.Mock
}

func (m *MockShipmentCollection) GetAll(ctx context.Context, tenant string) ([]database.ShipmentResponse, error) {
	args := m.Called(ctx, tenant)
	return args.Get(0).([]database.ShipmentResponse), args.Error(1)
}

func (m *MockShipmentCollection) Create(ctx context.Context, createShipmentParams database.Shipment) (string, error) {
	args := m.Called(ctx, createShipmentParams)
	return args.String(0), args.Error(1)
}

func (m *MockShipmentCollection) GetByID(ctx context.Context, id string, tenant string) (database.ShipmentResponse, error) {
	args := m.Called(ctx, id, tenant)
	return args.Get(0).(database.ShipmentResponse), args.Error(1)
}

func (m *MockShipmentCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (database.ShipmentResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).(database.ShipmentResponse), args.Error(1)
}

func (m *MockShipmentCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]database.ShipmentResponse, error) {
	args := m.Called(ctx, key, value, tenant)
	return args.Get(0).([]database.ShipmentResponse), args.Error(1)
}

func (m *MockShipmentCollection) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockShipmentCollection) Update(ctx context.Context, id string, updateShipmentParams database.Shipment) error {
	args := m.Called(ctx, id, updateShipmentParams)
	return args.Error(0)
}

func TestShipmentHandler_GetAllShipment(t *testing.T) {
	mockCollection := new(MockShipmentCollection)
	handler := NewShipmentHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/shipments", handler.GetAllShipment)

	mockCollection.On("GetAll", mock.Anything, mock.Anything).Return([]database.ShipmentResponse{
		{
			ID:          "66403d4660a455ef0b8f88f6",
			Tenant:      "customerA",
			MasterEmail: "james@gmail.com",
			// ETA:           time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
			InitialETA: time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
			CurrentETA: time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),

			VoyageNumber:  "abc",
			CurrentStatus: enum.SHIPMENT_STATUS_AT_ANCHORAGE,
			ShipmentType: database.ShipmentType{
				CargoOperations: database.CargoOperations{
					CargoOperations: true,
					CargoOperationsActivity: []*core.CargoOperationsActivity{
						{
							ActivityType:      "Loading",
							AnchorageLocation: "Anchorage A",
							TerminalName:      "Terminal A",
							ShipmentProduct: []*core.ShipmentProduct{{
								SubProductType: "Crude",
								Quantity:       1000,
								QuantityCode:   "MT",
								Percentage:     80,
							}},
							Readiness: time.Date(2024, time.May, 22, 10, 0, 0, 0, time.UTC),
							ETB:       time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
							ETD:       time.Date(2024, time.May, 23, 18, 0, 0, 0, time.UTC),
							ArrivalDepartureInformation: core.ArrivalDepartureInformation{
								ArrivalDisplacement:   40000,
								DepartureDisplacement: 42000,
								ArrivalDraft:          10.5,
								DepartureDraft:        11.0,
								ArrivalMastHeight:     30.0,
								DepartureMastHeight:   31.0,
							},
						},
					},
				},
				Bunkering: database.Bunkering{
					Bunkering: true,
					BunkeringActivity: []*core.BunkeringActivity{
						{
							CustomerName:      "customerA",
							Supplier:          "supplier1",
							SupplierContact:   "supplier1",
							AppointedSurveyor: "supplier1",
							Docking:           "starboard",
							SupplierVessel:    "supplier1",
							BunkerIntakeSpecifications: []*core.BunkerIntakeSpecifications{
								{
									SubProductType:        "hh",
									MaximumQuantityIntake: 2,
									MaximumHoseSize:       3,
								},
							},
							ShipmentProduct: []*core.ShipmentProduct{{
								SubProductType: "hh",
								Quantity:       1000,
								QuantityCode:   "MT",
								Percentage:     80,
							}},
							Freeboard: 2,
							Readiness: time.Date(2024, time.June, 21, 3, 16, 1, 0, time.UTC),
							ETB:       time.Date(2024, time.June, 21, 3, 16, 3, 0, time.UTC),
							ETD:       time.Date(2024, time.June, 21, 3, 16, 7, 0, time.UTC),
						},
					},
				},
				OwnerMatters: database.OwnerMatters{
					OwnerMatters: false,
					Activity:     []*core.Activity{},
				},
			},
			VesselSpecifications: database.VesselSpecifications{
				ImoNumber:  1234567,
				VesselName: "Oceanic Voyager",
				CallSign:   "OV1234",
				SDWT:       50000,
				NRT:        20000,
				Flag:       "Panama",
				GRT:        30000,
				LOA:        250.0,
			},
			ShipmentDetails: database.ShipmentDetails{
				Agent: database.Agent{
					Name:    "John Doe",
					Contact: "123-456-7890",
					Email:   "john.doe@example.com",
					Tenant:  "customerA",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil)

	req := httptest.NewRequest("GET", "/shipments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "66403d4660a455ef0b8f88f6") // check if the response contains certain data
	assert.Contains(t, w.Body.String(), "customerA")                // check if the response contains certain data
	mockCollection.AssertExpectations(t)

}

// func TestShipmentHandler_CreateShipment(t *testing.T) {
// 	mockCollection := new(MockShipmentCollection)
// 	handler := NewShipmentHandler(mockCollection)
// 	mockWhatsAppClient := new(MockWhatsAppClient)

// 	r := chi.NewRouter()
// 	r.Post("/shipments", handler.CreateShipment)

// 	mockCollection.On("Create", mock.Anything, mock.MatchedBy(func(shipment database.Shipment) bool {
// 		return shipment.MasterEmail == "james@gmail.com" && shipment.ShipmentType.CargoOperations.CargoOperationsActivity[0].TerminalName == "Terminal A"
// 	})).Return("id123", nil)

// 	mockWhatsAppClient.On("CreateShipmentWhatsAppMessage", "John Doe", "Oceanic Voyager", "1234567", "22-May-2024", "12:00", "123-456-7890").Return("", nil)

// 	req := httptest.NewRequest("POST", "/shipments", bytes.NewBufferString(`{
// 		"shipment_type": {
// 			"cargo_operations": {
// 				"cargo_operations": true,
// 				"cargo_operations_activity": [
// 					{
// 						"activity_type": "Loading",
// 						"anchorage_location": "Anchorage A",
// 						"terminal_name": "Terminal A",
// 						"shipment_product":[
// 							{
// 								"sub_product_type": "Crude",
// 								"quantity": 1000,
// 								"quantity_code": "MT",
// 								"percentage": 80
// 							}
// 						],
// 						"readiness": "2024-05-22T10:00:00Z",
// 						"etb": "2024-05-22T12:00:00Z",
// 						"etd": "2024-05-23T18:00:00Z",
// 						"arrival_departure_information": {
// 							"arrival_displacement": 40000,
// 							"departure_displacement": 42000,
// 							"arrival_draft": 10.5,
// 							"departure_draft": 11.0,
// 							"arrival_mast_height": 30.0,
// 							"departure_mast_height": 31.0
// 						}
// 					}
// 				]
// 			},
// 			"bunkering": {
// 				"bunkering": true,
// 				"bunkering_activity": [
// 					{
// 						"supplier": "supplier1",
// 						"supplier_contact": "supplier1",
// 						"appointed_surveyor": "supplier1",
// 						"docking": "starboard",
// 						"supplier_vessel": "supplier1",
// 						"bunker_intake_specifications": [
// 							{
// 								"sub_product_type": "hh",
// 								"maximum_quantity_intake": 2,
// 								"maximum_hose_size": 3
// 							}
// 						],
// 						"freeboard": 2,
// 						"readiness": "2024-06-21T03:16:01Z",
// 						"etb": "2024-06-21T03:16:03Z",
// 						"etd": "2024-06-21T03:16:07Z"
// 					}
// 				]
// 			},
// 			"owner_matters": {
// 				"owner_matters": false,
// 				"activity": []
// 			}
// 		},
// 		"tenant":"customerA",
// 		"master_email":"james@gmail.com",
// 		"ETA": "2024-05-22T12:00:00Z",
// 		"voyage_number":"abc",
// 		"current_status": "Not Started",
// 		"vessel_specifications": {
// 			"imo_number": 1234567,
// 			"vessel_name": "Oceanic Voyager",
// 			"call_sign": "OV1234",
// 			"sdwt": 50000,
// 			"nrt": 20000,
// 			"flag": "Panama",
// 			"grt": 30000,
// 			"loa": 250.0
// 		},
// 		"shipment_details": {
// 			"agent_details": {
// 				"name": "John Doe",
// 				"contact": "123-456-7890",
// 				"email": "john.doe@example.com"
// 			}
// 		}
// 	}`))

// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusCreated, w.Code)
// 	mockCollection.AssertExpectations(t)
// 	mockWhatsAppClient.AssertExpectations(t)

// }

func TestShipmentHandler_GetShipmentFromId(t *testing.T) {
	mockCollection := new(MockShipmentCollection)
	handler := NewShipmentHandler(mockCollection)

	r := chi.NewRouter()
	r.Get("/shipments/{shipment_id}", handler.GetShipmentFromId)

	mockCollection.On("GetByID", mock.Anything, "1", mock.Anything).Return(database.ShipmentResponse{
		ID:          "66403d4660a455ef0b8f88f6",
		Tenant:      "customerA",
		MasterEmail: "james@gmail.com",
		// ETA:           time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
		InitialETA: time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
		CurrentETA: time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),

		VoyageNumber:  "abc",
		CurrentStatus: enum.SHIPMENT_STATUS_NOT_STARTED,
		ShipmentType: database.ShipmentType{
			CargoOperations: database.CargoOperations{
				CargoOperations: true,
				CargoOperationsActivity: []*core.CargoOperationsActivity{
					{
						ActivityType:      "Loading",
						AnchorageLocation: "Anchorage A",
						TerminalName:      "Terminal A",
						ShipmentProduct: []*core.ShipmentProduct{{
							SubProductType: "Crude",
							Quantity:       1000,
							QuantityCode:   "MT",
							Percentage:     80,
						}},
						Readiness: time.Date(2024, time.May, 22, 10, 0, 0, 0, time.UTC),
						ETB:       time.Date(2024, time.May, 22, 12, 0, 0, 0, time.UTC),
						ETD:       time.Date(2024, time.May, 23, 18, 0, 0, 0, time.UTC),
						ArrivalDepartureInformation: core.ArrivalDepartureInformation{
							ArrivalDisplacement:   40000,
							DepartureDisplacement: 42000,
							ArrivalDraft:          10.5,
							DepartureDraft:        11.0,
							ArrivalMastHeight:     30.0,
							DepartureMastHeight:   31.0,
						},
					},
				},
			},
			Bunkering: database.Bunkering{
				Bunkering: true,
				BunkeringActivity: []*core.BunkeringActivity{
					{
						CustomerName:      "customerA",
						Supplier:          "supplier1",
						SupplierContact:   "supplier1",
						AppointedSurveyor: "supplier1",
						Docking:           "starboard",
						SupplierVessel:    "supplier1",
						BunkerIntakeSpecifications: []*core.BunkerIntakeSpecifications{
							{
								SubProductType:        "hh",
								MaximumQuantityIntake: 2,
								MaximumHoseSize:       3,
							},
						},
						ShipmentProduct: []*core.ShipmentProduct{{
							SubProductType: "hh",
							Quantity:       1000,
							QuantityCode:   "MT",
							Percentage:     80,
						}},
						Freeboard: 2,
						Readiness: time.Date(2024, time.June, 21, 3, 16, 1, 0, time.UTC),
						ETB:       time.Date(2024, time.June, 21, 3, 16, 3, 0, time.UTC),
						ETD:       time.Date(2024, time.June, 21, 3, 16, 7, 0, time.UTC),
					},
				},
			},
			OwnerMatters: database.OwnerMatters{
				OwnerMatters: true,
				Activity: []*core.Activity{
					{
						ActivityType:      "Inspection",
						AnchorageLocation: "Anchorage B",
						TerminalName:      "Terminal B",
						ShipmentProduct: []*core.ShipmentProduct{{
							SubProductType: "Acids",
							Quantity:       500,
							QuantityCode:   "L",
							Percentage:     50,
						}},
						Readiness: time.Date(2024, time.May, 23, 10, 0, 0, 0, time.UTC),
						ETB:       time.Date(2024, time.May, 23, 12, 0, 0, 0, time.UTC),
						ETD:       time.Date(2024, time.May, 24, 18, 0, 0, 0, time.UTC),
						ArrivalDepartureInformation: core.ArrivalDepartureInformation{
							ArrivalDisplacement:   20000,
							DepartureDisplacement: 21000,
							ArrivalDraft:          7.5,
							DepartureDraft:        8.0,
							ArrivalMastHeight:     20.0,
							DepartureMastHeight:   21.0,
						},
					},
				},
			},
		},
		VesselSpecifications: database.VesselSpecifications{
			ImoNumber:  1234567,
			VesselName: "Oceanic Voyager",
			CallSign:   "OV1234",
			SDWT:       50000,
			NRT:        20000,
			Flag:       "Panama",
			GRT:        30000,
			LOA:        250.0,
		},
		ShipmentDetails: database.ShipmentDetails{
			Agent: database.Agent{
				Name:    "John Doe",
				Contact: "123-456-7890",
				Email:   "john.doe@example.com",
				Tenant:  "customerA",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	req := httptest.NewRequest("GET", "/shipments/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shipment_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "66403d4660a455ef0b8f88f6")
	assert.Contains(t, w.Body.String(), "customerA")
	mockCollection.AssertExpectations(t)
}

func TestShipmentHandler_UpdateShipmentById(t *testing.T) {
	mockCollection := new(MockShipmentCollection)
	handler := NewShipmentHandler(mockCollection)

	r := chi.NewRouter()
	r.Put("/shipments/{shipment_id}", handler.UpdateShipmentById)

	mockCollection.On("Update", mock.Anything, "1", mock.MatchedBy(func(shipment database.Shipment) bool {
		return shipment.MasterEmail == "james@gmail.com" && shipment.Tenant == "customerA" && shipment.ShipmentType.CargoOperations.CargoOperationsActivity[0].TerminalName == "Terminal A"
	})).Return(nil)

	req := httptest.NewRequest("PUT", "/shipments/1", bytes.NewBufferString(`{
		"shipment_type": {
			"cargo_operations": {
				"cargo_operations": true,
				"cargo_operations_activity": [
					{
						"activity_type": "Loading",
						"anchorage_location": "Anchorage A",
						"terminal_name": "Terminal A",
						"shipment_product": [
							{
								"sub_product_type": "Crude",
								"quantity": 1000,
								"quantity_code": "MT",
								"percentage": 80
							}
						],
						"readiness": "2024-05-22T10:00:00Z",
						"etb": "2024-05-22T12:00:00Z",
						"etd": "2024-05-23T18:00:00Z",
						"arrival_departure_information": {
							"arrival_displacement": 40000,
							"departure_displacement": 42000,
							"arrival_draft": 10.5,
							"departure_draft": 11.0,
							"arrival_mast_height": 30.0,
							"departure_mast_height": 31.0
						}
					}
				]
			},
			"bunkering": {
				"bunkering": true,
				"bunkering_activity": [
					{
						"consumer_name": "customerA",
						"supplier": "supplier1",
						"supplier_contact": "supplier1",
						"appointed_surveyor": "supplier1",
						"docking": "starboard",
						"supplier_vessel": "supplier1",
						"bunker_intake_specifications": [
							{
								"sub_product_type": "hh",
								"maximum_quantity_intake": 2,
								"maximum_hose_size": 3
							}
						],
						"freeboard": 2,
						"readiness": "2024-06-21T03:16:01Z",
						"etb": "2024-06-21T03:16:03Z",
						"etd": "2024-06-21T03:16:07Z"
					}
				]
			},
			"owner_matters": {
				"owner_matters": false,
				"activity": []
			}
		},
		"tenant":"customerA",
		"master_email":"james@gmail.com",
		"ETA": "2024-05-22T12:00:00Z",
		"voyage_number":"abc",
		"current_status": "Not Started",
		"vessel_specifications": {
			"imo_number": 1234567,
			"vessel_name": "Oceanic Voyager",
			"call_sign": "OV1234",
			"sdwt": 50000,
			"nrt": 20000,
			"flag": "Panama",
			"grt": 30000,
			"loa": 250.0
		},
		"shipment_details": {
			"agent_details": {
				"name": "John Doe",
				"contact": "123-456-7890",
				"email": "john.doe@example.com"
			}
		}
	}`))
	req.Header.Set("Content-Type", "application/json")
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shipment_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}

func TestShipmentHandler_DeleteShipmentById(t *testing.T) {
	mockCollection := new(MockShipmentCollection)
	handler := NewShipmentHandler(mockCollection)

	r := chi.NewRouter()
	r.Delete("/shipments/{shipment_id}", handler.DeleteShipmentById)

	mockCollection.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest("DELETE", "/shipments/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("shipment_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockCollection.AssertExpectations(t)
}
