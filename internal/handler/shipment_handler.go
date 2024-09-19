package handler

import (
	"backend-crm/internal/clients"
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
	"backend-crm/pkg/enum"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShipmentHandler struct {
	ShipmentCollection database.Collection[database.Shipment, database.ShipmentResponse]
}

func NewShipmentHandler(
	shipmentCollection database.Collection[database.Shipment, database.ShipmentResponse],
) *ShipmentHandler {
	return &ShipmentHandler{
		ShipmentCollection: shipmentCollection,
	}
}

type getAllShipmentsResponse struct {
	Shipments []database.ShipmentResponse `json:"shipments"`
}

func (h *ShipmentHandler) GetAllShipment(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	var shipmentsList []database.ShipmentResponse

	shipmentsList, err := h.ShipmentCollection.GetAll(r.Context(), tenant)
	log.Println(shipmentsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var shipments []database.ShipmentResponse
	for _, shipment := range shipmentsList {
		shipments = append(shipments, database.ShipmentResponse{
			ID:                   shipment.ID,
			Tenant:               shipment.Tenant,
			MasterEmail:          shipment.MasterEmail,
			InitialETA:           shipment.InitialETA,
			CurrentETA:           shipment.CurrentETA,
			VoyageNumber:         shipment.VoyageNumber,
			CurrentStatus:        shipment.CurrentStatus,
			ShipmentType:         shipment.ShipmentType,
			VesselSpecifications: shipment.VesselSpecifications,
			ShipmentDetails:      shipment.ShipmentDetails,
			CreatedAt:            shipment.CreatedAt,
			UpdatedAt:            shipment.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllShipmentsResponse{Shipments: shipments}
	render.JSON(w, r, response)
}

func (h *ShipmentHandler) GetAllShipmentForETAReminder(w http.ResponseWriter, r *http.Request) ([]database.ShipmentResponse, error) {
	// email := auth.GetEmailFromToken(r.Context())
	// tenant := config.MakeMapping(email)
	var shipmentsList []database.ShipmentResponse
	// get all tenants in this case
	shipmentsList, err := h.ShipmentCollection.GetAll(r.Context(), "")
	log.Println(shipmentsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return nil, err
	}

	var shipments []database.ShipmentResponse
	for _, shipment := range shipmentsList {
		shipments = append(shipments, database.ShipmentResponse{
			ID:                   shipment.ID,
			Tenant:               shipment.Tenant,
			MasterEmail:          shipment.MasterEmail,
			InitialETA:           shipment.InitialETA,
			CurrentETA:           shipment.CurrentETA,
			VoyageNumber:         shipment.VoyageNumber,
			CurrentStatus:        shipment.CurrentStatus,
			ShipmentType:         shipment.ShipmentType,
			VesselSpecifications: shipment.VesselSpecifications,
			ShipmentDetails:      shipment.ShipmentDetails,
			CreatedAt:            shipment.CreatedAt,
			UpdatedAt:            shipment.UpdatedAt,
		})
	}

	return shipments, nil
}

func (h *ShipmentHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createShipmentParams database.Shipment
	log.Println(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&createShipmentParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createShipmentParams.Tenant = tenant

	log.Println(createShipmentParams)
	shipment_id, err := h.ShipmentCollection.Create(r.Context(), createShipmentParams)

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

	response := &SuccessResponse{
		HTTPStatusCode: SuccessCreated.HTTPStatusCode,
		Message:        SuccessCreated.Message,
		Data:           shipment_id,
	}
	render.Render(w, r, response)

	// Send WhatsApp message

	// Convert to local time
	// Load the location for GMT+8
	// location, err := time.LoadLocation("Asia/Singapore")
	// if err != nil {
	// 	log.Println("Error loading location:", err)
	// 	return
	// }
	systemTimeNow := time.Now()
	systemTimeNowLocation := systemTimeNow.Location()
	// Check if the location is not "Local"
	if systemTimeNowLocation.String() != "Local" {
		// Define a fixed time zone with +8 hours offset from UTC
		systemTimeNowLocation = time.FixedZone("UTC+8", 8*60*60)

	}
	log.Println("Current Time:", systemTimeNow)
	log.Println("Time Zone Location:", systemTimeNowLocation)

	localETA := createShipmentParams.CurrentETA.In(systemTimeNowLocation)
	log.Println("this is local ETA")

	log.Println(localETA, "this is local ETA")
	_, invalid := clients.WhatsAppClient.CreateShipmentWhatsAppMessage(
		createShipmentParams.ShipmentDetails.Agent.Name,
		createShipmentParams.VesselSpecifications.VesselName,
		strconv.FormatInt(createShipmentParams.VesselSpecifications.ImoNumber, 10), // Convert int64 to string
		localETA.Format("02-Jan-2006"),
		localETA.Format("15:04"),
		createShipmentParams.ShipmentDetails.Agent.Contact,
	)
	if invalid != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
	}
}

func (h *ShipmentHandler) DeleteShipmentById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.ShipmentCollection.Delete(r.Context(), _id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)

	render.Render(w, r, SuccessOK)
}

func (h *ShipmentHandler) UpdateShipmentById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval
	log.Println(r)
	var updateShipmentParams database.Shipment
	if err := json.NewDecoder(r.Body).Decode(&updateShipmentParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	// log.Println(updateShipmentParams)
	err := h.ShipmentCollection.Update(r.Context(), _id, updateShipmentParams)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)

	render.Render(w, r, SuccessOK)
}

func (h *ShipmentHandler) GetShipmentFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	shipment, err := h.ShipmentCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved shipment info from DB: %v", shipment) // Confirm shipment retrieval

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.ShipmentResponse{
		ID:                   shipment.ID,
		Tenant:               shipment.Tenant,
		MasterEmail:          shipment.MasterEmail,
		InitialETA:           shipment.InitialETA,
		CurrentETA:           shipment.CurrentETA,
		VoyageNumber:         shipment.VoyageNumber,
		CurrentStatus:        shipment.CurrentStatus,
		ShipmentType:         shipment.ShipmentType,
		VesselSpecifications: shipment.VesselSpecifications,
		ShipmentDetails:      shipment.ShipmentDetails,
		CreatedAt:            shipment.CreatedAt,
		UpdatedAt:            shipment.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *ShipmentHandler) FilterShipment(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	shipment, err := h.ShipmentCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.ShipmentResponse{
		ID:                   shipment.ID,
		Tenant:               shipment.Tenant,
		MasterEmail:          shipment.MasterEmail,
		InitialETA:           shipment.InitialETA,
		CurrentETA:           shipment.CurrentETA,
		VoyageNumber:         shipment.VoyageNumber,
		CurrentStatus:        shipment.CurrentStatus,
		ShipmentType:         shipment.ShipmentType,
		VesselSpecifications: shipment.VesselSpecifications,
		ShipmentDetails:      shipment.ShipmentDetails,
		CreatedAt:            shipment.CreatedAt,
		UpdatedAt:            shipment.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *ShipmentHandler) GetAllShipmentStatuses(w http.ResponseWriter, r *http.Request) {
	response := struct {
		ShipmentStatuses []string `json:"shipment_statuses"`
	}{
		ShipmentStatuses: enum.GetShipmentStatuses(),
	}
	render.JSON(w, r, response)
}

func (h *ShipmentHandler) GetAllShipmentStatusesWithColours(w http.ResponseWriter, r *http.Request) {
	response := struct {
		ShipmentStatuses map[string]string `json:"shipment_statuses"`
	}{
		ShipmentStatuses: enum.GetShipmentStatusesWithColours(),
	}
	render.JSON(w, r, response)
}

func (h *ShipmentHandler) GetAllAnchorageLocations(w http.ResponseWriter, r *http.Request) {
	response := struct {
		AnchorageLocations []string `json:"anchorage_locations"`
	}{
		AnchorageLocations: enum.GetAnchorageLocations(),
	}
	render.JSON(w, r, response)
}
