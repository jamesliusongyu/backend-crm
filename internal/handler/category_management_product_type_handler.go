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

type CategoryManagementProductTypeHandler struct {
	CategoryManagementProductTypeCollection database.Collection[database.CategoryManagementProductType, database.CategoryManagementProductTypeResponse]
}

func NewProductTypeHandler(
	categoryManagementProductTypeCollection database.Collection[database.CategoryManagementProductType, database.CategoryManagementProductTypeResponse],
) *CategoryManagementProductTypeHandler {
	return &CategoryManagementProductTypeHandler{
		CategoryManagementProductTypeCollection: categoryManagementProductTypeCollection,
	}
}

type getAllCategoryManagementProductTypeResponse struct {
	CategoryManagementProductType []database.CategoryManagementProductTypeResponse `json:"product_types"`
}

func (h *CategoryManagementProductTypeHandler) GetAllCategoryManagementProductTypes(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var productTypeList []database.CategoryManagementProductTypeResponse

	categoryManagementProductTypeList, err := h.CategoryManagementProductTypeCollection.GetAll(r.Context(), tenant)
	log.Println(productTypeList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := getAllCategoryManagementProductTypeResponse{CategoryManagementProductType: categoryManagementProductTypeList}
	render.JSON(w, r, response)
}

type GetAllOnlySubProductTypesResponse struct {
	OnlySubProductTypes []string `json:"only_sub_product_types"`
}

func (h *CategoryManagementProductTypeHandler) GetAllOnlySubProductTypes(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	log.Println(tenant)
	var subProductTypeListOnly []string

	categoryManagementProductTypeList, err := h.CategoryManagementProductTypeCollection.GetAll(r.Context(), tenant)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	for _, product := range categoryManagementProductTypeList {
		subProductTypeListOnly = append(subProductTypeListOnly, product.SubProductsType...)
	}

	response := GetAllOnlySubProductTypesResponse{OnlySubProductTypes: subProductTypeListOnly}
	render.JSON(w, r, response)
}

func (h *CategoryManagementProductTypeHandler) CreateProductType(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)

	var createParams database.CategoryManagementProductType
	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	createParams.Tenant = tenant

	log.Println(createParams)
	_, err := h.CategoryManagementProductTypeCollection.Create(r.Context(), createParams)

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

func (h *CategoryManagementProductTypeHandler) DeleteProductTypeById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "product_type_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	err := h.CategoryManagementProductTypeCollection.Delete(r.Context(), _id)
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

func (h *CategoryManagementProductTypeHandler) GetProductTypeFromId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := chi.URLParam(r, "product_type_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	productType, err := h.CategoryManagementProductTypeCollection.GetByID(r.Context(), _id, tenant)
	log.Printf("Retrieved ProductType info from DB: %v", productType) // Confirm ProductType retrieval

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CategoryManagementProductTypeResponse{
		ID:              productType.ID,
		Tenant:          productType.Tenant,
		ProductType:     productType.ProductType,
		SubProductsType: productType.SubProductsType,
	}

	render.JSON(w, r, response)
}

func (h *CategoryManagementProductTypeHandler) UpdateProductTypeById(w http.ResponseWriter, r *http.Request) {
	_id := chi.URLParam(r, "product_type_id")
	log.Printf("Retrieved ID from URL: %s", _id) // Confirm ID retrieval

	var updateParams database.CategoryManagementProductType
	if err := json.NewDecoder(r.Body).Decode(&updateParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err := h.CategoryManagementProductTypeCollection.Update(r.Context(), _id, updateParams)

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

func (h *CategoryManagementProductTypeHandler) FilterProductType(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	_id := r.URL.Query().Get("_id")
	log.Println(_id)

	productType, err := h.CategoryManagementProductTypeCollection.GetByID(r.Context(), _id, tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	response := database.CategoryManagementProductTypeResponse{
		ID:              productType.ID,
		Tenant:          productType.Tenant,
		ProductType:     productType.ProductType,
		SubProductsType: productType.SubProductsType,
	}

	render.JSON(w, r, response)
}
