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

type CategoryManagementActivityTypeHandler struct {
	CategoryManagementActivityTypeCollection database.Collection[database.CategoryManagementActivityType, database.CategoryManagementActivityTypeResponse]
}

func NewActivityTypeHandler(
	categoryManagementActivityTypeCollection database.Collection[database.CategoryManagementActivityType, database.CategoryManagementActivityTypeResponse],
) *CategoryManagementActivityTypeHandler {
	return &CategoryManagementActivityTypeHandler{
		CategoryManagementActivityTypeCollection: categoryManagementActivityTypeCollection,
	}
}

type getAllCategoryManagementActivityTypeResponse struct {
	CategoryManagementActivityType []database.CategoryManagementActivityTypeResponse `json:"activity_types"`
}

func (h *CategoryManagementActivityTypeHandler) GetAllCategoryManagementActivityTypes(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var activityTypeList []database.CategoryManagementActivityTypeResponse

	categoryManagementActivityTypeList, err := h.CategoryManagementActivityTypeCollection.GetAll(r.Context(), tenant)
	log.Println(activityTypeList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := getAllCategoryManagementActivityTypeResponse{CategoryManagementActivityType: categoryManagementActivityTypeList}
	render.JSON(w, r, response)
}

func (h *CategoryManagementActivityTypeHandler) CreateActivityType(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.CategoryManagementActivityType
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.CategoryManagementActivityTypeCollection.Create(r.Context(), createParams)

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

func (h *CategoryManagementActivityTypeHandler) DeleteActivityTypeById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "activity_type_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.CategoryManagementActivityTypeCollection.Delete(r.Context(), _id)
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

// func (h *ActivityTypeHandler) UpdateActivityTypeById(w http.ResponseWriter, r *http.Request) {
// 	_id := chi.URLParam(r, "activity_type_id")
// 	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

// 	var updateParams database.ActivityType
// 	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
// 		render.Render(w, r, ErrBadRequest)
// 		return
// 	}

// 	err := h.ActivityTypeCollection.Update(r.Context(), _id, updateParams)

// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			render.Render(w, r, ErrNotFound)
// 		} else {
// 			render.Render(w, r, ErrInternalServerError)
// 		}
// 		return
// 	}

// 	w.WriteHeader(200)
// 	w.Write(nil)

// 	render.Render(w, r, SuccessOK)
// }

func (h *CategoryManagementActivityTypeHandler) GetActivityTypeFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "activity_type_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	activityType, err := h.CategoryManagementActivityTypeCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved ActivityType info from DB: %v", activityType) // Confirm ActivityType retrieval

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CategoryManagementActivityTypeResponse{
		ID:           activityType.ID,
		Tenant:       activityType.Tenant,
		ActivityType: activityType.ActivityType,
	}

	render.JSON(w, r, response)
}

func (h *CategoryManagementActivityTypeHandler) FilterActivityType(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	activityType, err := h.CategoryManagementActivityTypeCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CategoryManagementActivityTypeResponse{
		ID:           activityType.ID,
		Tenant:       activityType.Tenant,
		ActivityType: activityType.ActivityType,
	}

	render.JSON(w, r, response)
}
