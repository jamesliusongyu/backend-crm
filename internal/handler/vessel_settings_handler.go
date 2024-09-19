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

type VesselHandler struct {
	VesselCollection database.Collection[database.Vessel, database.VesselManagementResponse]
}

func NewVesselHandler(
	vesselCollection database.Collection[database.Vessel, database.VesselManagementResponse],
) *VesselHandler {
	return &VesselHandler{
		VesselCollection: vesselCollection,
	}
}

type getAllVesselsResponse struct {
	Vessels []database.VesselManagementResponse `json:"vessels"`
}

func (h *VesselHandler) GetAllVessels(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var VesselsList []database.VesselManagementResponse

	VesselsList, err := h.VesselCollection.GetAll(r.Context(), tenant)
	log.Println(VesselsList)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var Vessels []database.VesselManagementResponse
	for _, Vessel := range VesselsList {
		Vessels = append(Vessels, database.VesselManagementResponse{
			ID:                   Vessel.ID,
			Tenant:               Vessel.Tenant,
			VesselSpecifications: Vessel.VesselSpecifications,
			CreatedAt:            Vessel.CreatedAt,
			UpdatedAt:            Vessel.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllVesselsResponse{Vessels: Vessels}
	render.JSON(w, r, response)
}

func (h *VesselHandler) CreateVessel(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.Vessel
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.VesselCollection.Create(r.Context(), createParams)

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

func (h *VesselHandler) DeleteVesselById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "vessel_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.VesselCollection.Delete(r.Context(), _id)
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

func (h *VesselHandler) UpdateVesselById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "vessel_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Vessel
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.VesselCollection.Update(r.Context(), _id, updateParams)

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

func (h *VesselHandler) GetVesselById(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "vessel_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	Vessel, err := h.VesselCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved Vessel info from DB: %v", Vessel) // Confirm Vessel retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.VesselManagementResponse{
		ID:                   Vessel.ID,
		Tenant:               Vessel.Tenant,
		VesselSpecifications: Vessel.VesselSpecifications,
		CreatedAt:            Vessel.CreatedAt,
		UpdatedAt:            Vessel.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *VesselHandler) FilterVessel(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	Vessel, err := h.VesselCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.VesselManagementResponse{
		ID:                   Vessel.ID,
		Tenant:               Vessel.Tenant,
		VesselSpecifications: Vessel.VesselSpecifications,
		CreatedAt:            Vessel.CreatedAt,
		UpdatedAt:            Vessel.UpdatedAt,
	}
	render.JSON(w, r, response)
}
