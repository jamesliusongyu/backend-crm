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

type SupplierHandler struct {
	SupplierCollection database.Collection[database.Supplier, database.SupplierManagementResponse]
}

func NewSupplierHandler(
	SupplierCollection database.Collection[database.Supplier, database.SupplierManagementResponse],
) *SupplierHandler {
	return &SupplierHandler{
		SupplierCollection: SupplierCollection,
	}
}

type getAllSuppliersResponse struct {
	Suppliers []database.SupplierManagementResponse `json:"suppliers"`
}

func (h *SupplierHandler) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var SuppliersList []database.SupplierManagementResponse

	SuppliersList, err := h.SupplierCollection.GetAll(r.Context(), tenant)
	log.Println(SuppliersList)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var Suppliers []database.SupplierManagementResponse
	for _, Supplier := range SuppliersList {
		Suppliers = append(Suppliers, database.SupplierManagementResponse{
			ID:                     Supplier.ID,
			Tenant:                 Supplier.Tenant,
			SupplierSpecifications: Supplier.SupplierSpecifications,
			CreatedAt:              Supplier.CreatedAt,
			UpdatedAt:              Supplier.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllSuppliersResponse{Suppliers: Suppliers}
	render.JSON(w, r, response)
}

func (h *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.Supplier
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.SupplierCollection.Create(r.Context(), createParams)

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

func (h *SupplierHandler) DeleteSupplierById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "supplier_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.SupplierCollection.Delete(r.Context(), _id)
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

func (h *SupplierHandler) UpdateSupplierById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "supplier_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Supplier
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.SupplierCollection.Update(r.Context(), _id, updateParams)

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

func (h *SupplierHandler) GetSupplierById(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "supplier_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	Supplier, err := h.SupplierCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved Supplier info from DB: %v", Supplier) // Confirm Supplier retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.SupplierManagementResponse{
		ID:                     Supplier.ID,
		Tenant:                 Supplier.Tenant,
		SupplierSpecifications: Supplier.SupplierSpecifications,
		CreatedAt:              Supplier.CreatedAt,
		UpdatedAt:              Supplier.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *SupplierHandler) FilterSupplier(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	Supplier, err := h.SupplierCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.SupplierManagementResponse{
		ID:                     Supplier.ID,
		Tenant:                 Supplier.Tenant,
		SupplierSpecifications: Supplier.SupplierSpecifications,
		CreatedAt:              Supplier.CreatedAt,
		UpdatedAt:              Supplier.UpdatedAt,
	}
	render.JSON(w, r, response)
}
