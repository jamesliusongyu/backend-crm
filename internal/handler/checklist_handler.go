package handler

import (
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ChecklistHandler struct {
	ChecklistCollection database.Collection[database.Checklist, database.ChecklistResponse]
}

func NewChecklistHandler(
	ChecklistCollection database.Collection[database.Checklist, database.ChecklistResponse],
) *ChecklistHandler {
	return &ChecklistHandler{
		ChecklistCollection: ChecklistCollection,
	}
}

type getAllChecklistResponse struct {
	Checklist []database.ChecklistResponse `json:"checklist"`
}

func (h *ChecklistHandler) GetAllChecklist(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var ChecklistsResult []database.ChecklistResponse

	ChecklistsResult, err := h.ChecklistCollection.GetAll(r.Context(), tenant)
	log.Println(ChecklistsResult)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var Checklist []database.ChecklistResponse
	for _, cl := range ChecklistsResult {
		Checklist = append(Checklist, database.ChecklistResponse{
			ID:     cl.ID,
			Tenant: cl.Tenant,

			PortDues:         cl.PortDues,
			Pilotage:         cl.Pilotage,
			ServiceLaunch:    cl.ServiceLaunch,
			Logistics:        cl.Logistics,
			HotelCharges:     cl.HotelCharges,
			AirTickets:       cl.AirTickets,
			TransportCharges: cl.TransportCharges,
			MedicineSupplies: cl.MedicineSupplies,
			FreshWaterSupply: cl.FreshWaterSupply,
			MarineAdvisory:   cl.MarineAdvisory,
			CourierServices:  cl.CourierServices,
			CrossHarbourFees: cl.CrossHarbourFees,
			SupplyBoat:       cl.SupplyBoat,
			Repairs:          cl.Repairs,
			Extras:           cl.Extras,
			CrewChange:       cl.CrewChange,

			ShipmentID: cl.ShipmentID,

			CreatedAt: cl.CreatedAt,
			UpdatedAt: cl.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllChecklistResponse{Checklist: Checklist}
	render.JSON(w, r, response)
}

func (h *ChecklistHandler) CreateChecklist(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.Checklist
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.ChecklistCollection.Create(r.Context(), createParams)

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

	render.Render(w, r, SuccessCreated)
}

func (h *ChecklistHandler) DeleteChecklistById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "shipping_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.ChecklistCollection.Delete(r.Context(), _id)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
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

func (h *ChecklistHandler) UpdateChecklistById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Checklist
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.ChecklistCollection.Update(r.Context(), _id, updateParams)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
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

func (h *ChecklistHandler) GetChecklistById(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	// Checklist, err := h.ChecklistCollection.GetByID(r.Context(), _id, ten ant)
	Checklist, err := h.ChecklistCollection.GetByKeyValue(r.Context(), "shipmentid", _id, tenant)
	log.Printf("Retrieved Checklist info from DB: %v", Checklist) // Confirm Checklist retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.ChecklistResponse{
		ID:               Checklist.ID,
		Tenant:           Checklist.Tenant,
		PortDues:         Checklist.PortDues,
		Pilotage:         Checklist.Pilotage,
		ServiceLaunch:    Checklist.ServiceLaunch,
		Logistics:        Checklist.Logistics,
		HotelCharges:     Checklist.HotelCharges,
		AirTickets:       Checklist.AirTickets,
		TransportCharges: Checklist.TransportCharges,
		MedicineSupplies: Checklist.MedicineSupplies,
		FreshWaterSupply: Checklist.FreshWaterSupply,
		MarineAdvisory:   Checklist.MarineAdvisory,
		CourierServices:  Checklist.CourierServices,
		CrossHarbourFees: Checklist.CrossHarbourFees,
		SupplyBoat:       Checklist.SupplyBoat,
		Repairs:          Checklist.Repairs,
		Extras:           Checklist.Extras,
		CrewChange:       Checklist.CrewChange,
		ShipmentID:       Checklist.ShipmentID,
		CreatedAt:        Checklist.CreatedAt,
		UpdatedAt:        Checklist.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *ChecklistHandler) FilterChecklist(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	Checklist, err := h.ChecklistCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.ChecklistResponse{
		ID:               Checklist.ID,
		Tenant:           Checklist.Tenant,
		PortDues:         Checklist.PortDues,
		Pilotage:         Checklist.Pilotage,
		ServiceLaunch:    Checklist.ServiceLaunch,
		Logistics:        Checklist.Logistics,
		HotelCharges:     Checklist.HotelCharges,
		AirTickets:       Checklist.AirTickets,
		TransportCharges: Checklist.TransportCharges,
		MedicineSupplies: Checklist.MedicineSupplies,
		FreshWaterSupply: Checklist.FreshWaterSupply,
		MarineAdvisory:   Checklist.MarineAdvisory,
		CourierServices:  Checklist.CourierServices,
		CrossHarbourFees: Checklist.CrossHarbourFees,
		SupplyBoat:       Checklist.SupplyBoat,
		Repairs:          Checklist.Repairs,
		Extras:           Checklist.Extras,
		CrewChange:       Checklist.CrewChange,
		ShipmentID:       Checklist.ShipmentID,
		CreatedAt:        Checklist.CreatedAt,
		UpdatedAt:        Checklist.UpdatedAt,
	}
	render.JSON(w, r, response)
}
