package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// type CategoryManagementProductType struct {
// 	Products core.Product `json:"products"`
// 	Tenant   string       `json:"tenant"`
// }

type CategoryManagementProductType struct {
	ProductType     string   `json:"product_type"`
	SubProductsType []string `json:"sub_products_type"`
	Tenant          string   `json:"tenant"`
}

type CategoryManagementProductTypeResponse struct {
	ID              string   `bson:"_id"`
	ProductType     string   `json:"product_type"`
	SubProductsType []string `json:"sub_products_type"`
	Tenant          string   `json:"tenant"`
}

type ProductTypeCollection struct {
	*GenericCollection[CategoryManagementProductType, CategoryManagementProductTypeResponse]
}

func NewProductTypeCollection(collection *mongo.Collection) *ProductTypeCollection {
	return &ProductTypeCollection{
		GenericCollection: NewGenericCollection[CategoryManagementProductType, CategoryManagementProductTypeResponse](collection),
	}
}

var _ Collection[CategoryManagementProductType, CategoryManagementProductTypeResponse] = (*ProductTypeCollection)(nil)

func (r *ProductTypeCollection) Create(ctx context.Context, entity CategoryManagementProductType) (string, error) {
	ProductType := CategoryManagementProductType{
		ProductType:     entity.ProductType,
		SubProductsType: entity.SubProductsType,
		Tenant:          entity.Tenant,
	}
	return r.GenericCollection.Create(ctx, ProductType)
}

func (r *ProductTypeCollection) Update(ctx context.Context, id string, entity CategoryManagementProductType) error {

	ProductType := CategoryManagementProductType{
		ProductType:     entity.ProductType,
		SubProductsType: entity.SubProductsType,
		Tenant:          entity.Tenant,
	}
	return r.GenericCollection.Update(ctx, id, ProductType)
}

func (r *ProductTypeCollection) GetAll(ctx context.Context, tenant string) ([]CategoryManagementProductTypeResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *ProductTypeCollection) GetByID(ctx context.Context, id string, tenant string) (CategoryManagementProductTypeResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *ProductTypeCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
