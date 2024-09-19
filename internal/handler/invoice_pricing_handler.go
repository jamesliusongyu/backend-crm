package handler

import (
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/auth"
	"backend-crm/pkg/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-chi/render"
)

type InvoicePricingHandler struct {
	InvoicePricingCollection database.Collection[database.InvoicePricing, database.InvoicePricingResponse]
}

func NewInvoicePricingHandler(
	invoicePricingCollection database.Collection[database.InvoicePricing, database.InvoicePricingResponse],
) *InvoicePricingHandler {
	return &InvoicePricingHandler{
		InvoicePricingCollection: invoicePricingCollection,
	}
}

func (h *InvoicePricingHandler) GetTenant(w http.ResponseWriter, r *http.Request) {

	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	response := struct {
		Tenant string `json:"tenant"`
	}{
		Tenant: tenant,
	}
	render.JSON(w, r, response)
}

func (h *InvoicePricingHandler) GetInvoiceFeesFromPortAuthority(w http.ResponseWriter, r *http.Request) {

	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	// Load BluShipping data
	data, err := utils.LoadInvoiceFeesData(tenant)
	if err != nil {
		http.Error(w, "Failed to load BluShipping data", http.StatusInternalServerError)
		return
	}

	// Create a response map
	response := map[string]interface{}{
		"tenant":      tenant,
		"invoiceFees": data,
	}

	// Render the response as JSON
	render.JSON(w, r, response)
}

func (h *InvoicePricingHandler) CreatePDAInvoicePricing(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createPDAInvoicePricingParams database.InvoicePricing
	// Read the body into a byte slice
	// body, invalid := io.ReadAll(r.Body)
	// if invalid != nil {
	// 	render.Render(w, r, ErrBadRequest)
	// 	return
	// }

	// Log the request body
	// log.Println("Request Body:", string(body))

	if err := json.NewDecoder(r.Body).Decode(&createPDAInvoicePricingParams); err != nil {
		log.Println(err, "err")
		render.Render(w, r, ErrBadRequest)
		return
	}
	createPDAInvoicePricingParams.Tenant = tenant

	log.Println(createPDAInvoicePricingParams)
	_, err := h.InvoicePricingCollection.Create(r.Context(), createPDAInvoicePricingParams)

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

func (h *InvoicePricingHandler) GetPDAInvoicePricingFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	_id := chi.URLParam(r, "invoice_id")
	log.Printf("Retrieved invoice ID from URL: %s", _id) // Confirm ID retrieval

	invoice, err := h.InvoicePricingCollection.GetByKeyValue(r.Context(), "shipmentid", _id, tenant)
	log.Println(err, "err")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.JSON(w, r, struct{}{}) // Return empty JSON payload so that frontend can handle
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.InvoicePricingResponse{
		ID:                    invoice.ID,
		Tenant:                invoice.Tenant,
		ShipmentID:            invoice.ShipmentID,
		InvoicePricingDetails: invoice.InvoicePricingDetails,
		CreatedAt:             invoice.CreatedAt,
		UpdatedAt:             invoice.UpdatedAt,
	}

	render.JSON(w, r, response)
}

func (h *InvoicePricingHandler) EditPDAInvoicePricing(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "invoice_id")
	log.Printf("Retrieved invoice ID from URL: %s", _id) // Confirm ID retrieval
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var updatePDAInvoicePricingParams database.InvoicePricing
	if err := json.NewDecoder(r.Body).Decode(&updatePDAInvoicePricingParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	invoice, err := h.InvoicePricingCollection.GetByKeyValue(r.Context(), "shipmentid", _id, tenant)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	invalid := h.InvoicePricingCollection.Update(r.Context(), invoice.ID, updatePDAInvoicePricingParams)
	log.Println(err, "err")
	if err != nil {
		if invalid == mongo.ErrNoDocuments {
			render.JSON(w, r, struct{}{}) // Return empty JSON payload so that frontend can handle
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)

	render.Render(w, r, SuccessOK)
}
