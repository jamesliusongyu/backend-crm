package database

import (
	"backend-crm/pkg/core"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Supplier struct {
	Tenant                 string                      `json:"tenant"`
	SupplierSpecifications core.SupplierSpecifications `json:"supplier_specifications"`
	CreatedAt              time.Time                   `json:"created_at"`
	UpdatedAt              time.Time                   `json:"updated_at"`
}

type SupplierManagementResponse struct {
	ID                     string                      `bson:"_id"`
	Tenant                 string                      `json:"tenant"`
	SupplierSpecifications core.SupplierSpecifications `json:"supplier_specifications"`
	CreatedAt              time.Time                   `json:"created_at"`
	UpdatedAt              time.Time                   `json:"updated_at"`
}

type SupplierCollection struct {
	*GenericCollection[Supplier, SupplierManagementResponse]
}

func NewSupplierCollection(collection *mongo.Collection) *SupplierCollection {
	return &SupplierCollection{
		GenericCollection: NewGenericCollection[Supplier, SupplierManagementResponse](collection),
	}
}

var _ Collection[Supplier, SupplierManagementResponse] = (*SupplierCollection)(nil)

func (r *SupplierCollection) Create(ctx context.Context, entity Supplier) (string, error) {

	Supplier := Supplier{
		Tenant:                 entity.Tenant,
		SupplierSpecifications: entity.SupplierSpecifications,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	return r.GenericCollection.Create(ctx, Supplier)
}

func (r *SupplierCollection) Update(ctx context.Context, id string, entity Supplier) error {
	Supplier := Supplier{
		Tenant:                 entity.Tenant,
		SupplierSpecifications: entity.SupplierSpecifications,
		CreatedAt:              entity.CreatedAt,
		UpdatedAt:              time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, Supplier)
}

func (r *SupplierCollection) GetAll(ctx context.Context, tenant string) ([]SupplierManagementResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *SupplierCollection) GetByID(ctx context.Context, id string, tenant string) (SupplierManagementResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *SupplierCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
