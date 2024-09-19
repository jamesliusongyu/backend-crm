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

type TerminalHandler struct {
	TerminalCollection database.Collection[database.Terminal, database.TerminalManagementResponse]
}

func NewTerminalHandler(
	TerminalCollection database.Collection[database.Terminal, database.TerminalManagementResponse],
) *TerminalHandler {
	return &TerminalHandler{
		TerminalCollection: TerminalCollection,
	}
}

type getAllTerminalsResponse struct {
	Terminals []database.TerminalManagementResponse `json:"terminals"`
}

func (h *TerminalHandler) GetAllTerminals(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var TerminalsList []database.TerminalManagementResponse

	TerminalsList, err := h.TerminalCollection.GetAll(r.Context(), tenant)
	log.Println(TerminalsList)

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var Terminals []database.TerminalManagementResponse
	for _, Terminal := range TerminalsList {
		Terminals = append(Terminals, database.TerminalManagementResponse{
			ID:                     Terminal.ID,
			Tenant:                 Terminal.Tenant,
			TerminalSpecifications: Terminal.TerminalSpecifications,
			CreatedAt:              Terminal.CreatedAt,
			UpdatedAt:              Terminal.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllTerminalsResponse{Terminals: Terminals}
	render.JSON(w, r, response)
}

func (h *TerminalHandler) CreateTerminal(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.Terminal
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.TerminalCollection.Create(r.Context(), createParams)

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

func (h *TerminalHandler) DeleteTerminalById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "terminal_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.TerminalCollection.Delete(r.Context(), _id)
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

func (h *TerminalHandler) UpdateTerminalById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "terminal_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Terminal
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.TerminalCollection.Update(r.Context(), _id, updateParams)

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

func (h *TerminalHandler) GetTerminalById(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "terminal_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	Terminal, err := h.TerminalCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved Terminal info from DB: %v", Terminal) // Confirm Terminal retrieval

	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.TerminalManagementResponse{
		ID:                     Terminal.ID,
		Tenant:                 Terminal.Tenant,
		TerminalSpecifications: Terminal.TerminalSpecifications,
		CreatedAt:              Terminal.CreatedAt,
		UpdatedAt:              Terminal.UpdatedAt,
	}
	render.JSON(w, r, response)
}

func (h *TerminalHandler) FilterTerminal(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	Terminal, err := h.TerminalCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		var rnfErr *database.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}
	response := database.TerminalManagementResponse{
		ID:                     Terminal.ID,
		Tenant:                 Terminal.Tenant,
		TerminalSpecifications: Terminal.TerminalSpecifications,
		CreatedAt:              Terminal.CreatedAt,
		UpdatedAt:              Terminal.UpdatedAt,
	}
	render.JSON(w, r, response)
}
