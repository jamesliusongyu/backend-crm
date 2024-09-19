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

type CustomerHandler struct {
	CustomerCollection database.Collection[database.Customer, database.CustomerResponse]
}

func NewCustomerHandler(
	customerCollection database.Collection[database.Customer, database.CustomerResponse],
) *CustomerHandler {
	return &CustomerHandler{
		CustomerCollection: customerCollection,
	}
}

type getAllCustomersResponse struct {
	Customers []database.CustomerResponse `json:"customers"`
}

func (h *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Context())
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var customersList []database.CustomerResponse

	customersList, err := h.CustomerCollection.GetAll(r.Context(), tenant)
	log.Println(customersList)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var customers []database.CustomerResponse
	for _, customer := range customersList {
		customers = append(customers, database.CustomerResponse{
			ID:        customer.ID,
			Tenant:    customer.Tenant,
			Customer:  customer.Customer,
			Company:   customer.Company,
			Email:     customer.Email,
			Contact:   customer.Contact,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllCustomersResponse{Customers: customers}
	render.JSON(w, r, response)
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	log.Println(r.Context())
	log.Println("ss")
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	log.Println(tenant)
	var createParams database.Customer
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.CustomerCollection.Create(r.Context(), createParams)

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

func (h *CustomerHandler) DeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "customer_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.CustomerCollection.Delete(r.Context(), _id)
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

func (h *CustomerHandler) UpdateCustomerById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "customer_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Customer
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.CustomerCollection.Update(r.Context(), _id, updateParams)

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

func (h *CustomerHandler) GetCustomerFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	_id := chi.URLParam(r, "customer_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	customer, err := h.CustomerCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved customer info from DB: %v", customer) // Confirm customer retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CustomerResponse{
		ID:        customer.ID,
		Tenant:    customer.Tenant,
		Customer:  customer.Customer,
		Company:   customer.Company,
		Email:     customer.Email,
		Contact:   customer.Contact,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *CustomerHandler) GetCustomerFromName(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	_name := chi.URLParam(r, "customer_name")
	log.Printf("Retrieved Name from URL: %s", _name) // Confirm Name retrieval

	customer, err := h.CustomerCollection.GetByKeyValue(r.Context(), "customer", _name, tenant)
	log.Printf("Retrieved customer info from DB: %v", customer) // Confirm customer retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CustomerResponse{
		ID:        customer.ID,
		Tenant:    customer.Tenant,
		Customer:  customer.Customer,
		Company:   customer.Company,
		Email:     customer.Email,
		Contact:   customer.Contact,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *CustomerHandler) FilterCustomer(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	customer, err := h.CustomerCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.CustomerResponse{
		ID:        customer.ID,
		Tenant:    customer.Tenant,
		Customer:  customer.Customer,
		Company:   customer.Company,
		Email:     customer.Email,
		Contact:   customer.Contact,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	render.JSON(w, r, response)
}
