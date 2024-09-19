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
	"go.mongodb.org/mongo-driver/mongo"
)

type AgentHandler struct {
	AgentCollection database.Collection[database.Agent, database.AgentManagementResponse]
}

func NewAgentHandler(
	agentCollection database.Collection[database.Agent, database.AgentManagementResponse],
) *AgentHandler {
	return &AgentHandler{
		AgentCollection: agentCollection,
	}
}

type getAllAgentsResponse struct {
	Agents []database.AgentManagementResponse `json:"agents"`
}

func (h *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var agentsList []database.AgentManagementResponse

	agentsList, err := h.AgentCollection.GetAll(r.Context(), tenant)
	log.Println(agentsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var agents []database.AgentManagementResponse
	for _, agent := range agentsList {
		agents = append(agents, database.AgentManagementResponse{
			ID:        agent.ID,
			Tenant:    agent.Tenant,
			Name:      agent.Name,
			Email:     agent.Email,
			Contact:   agent.Contact,
			CreatedAt: agent.CreatedAt,
			UpdatedAt: agent.UpdatedAt,
		})
	}

	// Creating the response object
	response := getAllAgentsResponse{Agents: agents}
	render.JSON(w, r, response)
}

func (h *AgentHandler) CreateAgent(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.Agent
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.AgentCollection.Create(r.Context(), createParams)

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

func (h *AgentHandler) DeleteAgentById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "agent_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.AgentCollection.Delete(r.Context(), _id)
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

func (h *AgentHandler) UpdateAgentById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "agent_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.Agent
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.AgentCollection.Update(r.Context(), _id, updateParams)

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

func (h *AgentHandler) GetAgentFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "agent_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	agent, err := h.AgentCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved Agent info from DB: %v", agent) // Confirm Agent retrieval

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.AgentManagementResponse{
		ID:        agent.ID,
		Tenant:    agent.Tenant,
		Name:      agent.Name,
		Email:     agent.Email,
		Contact:   agent.Contact,
		CreatedAt: agent.CreatedAt,
		UpdatedAt: agent.UpdatedAt,
	}

	render.JSON(w, r, response)
}

func (h *AgentHandler) FilterAgent(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	agent, err := h.AgentCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.AgentManagementResponse{
		ID:        agent.ID,
		Tenant:    agent.Tenant,
		Name:      agent.Name,
		Email:     agent.Email,
		Contact:   agent.Contact,
		CreatedAt: agent.CreatedAt,
		UpdatedAt: agent.UpdatedAt,
	}

	render.JSON(w, r, response)
}
