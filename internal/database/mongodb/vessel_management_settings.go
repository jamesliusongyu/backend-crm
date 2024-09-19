package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Vessel struct {
	Tenant               string               `json:"tenant"`
	VesselSpecifications VesselSpecifications `json:"vessel_specifications"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type VesselManagementResponse struct {
	ID                   string               `bson:"_id"`
	Tenant               string               `json:"tenant"`
	VesselSpecifications VesselSpecifications `json:"vessel_specifications"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type VesselCollection struct {
	*GenericCollection[Vessel, VesselManagementResponse]
}

func NewVesselCollection(collection *mongo.Collection) *VesselCollection {
	return &VesselCollection{
		GenericCollection: NewGenericCollection[Vessel, VesselManagementResponse](collection),
	}
}

var _ Collection[Vessel, VesselManagementResponse] = (*VesselCollection)(nil)

func (r *VesselCollection) Create(ctx context.Context, entity Vessel) (string, error) {

	vessel := Vessel{
		Tenant:               entity.Tenant,
		VesselSpecifications: entity.VesselSpecifications,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	return r.GenericCollection.Create(ctx, vessel)
}

func (r *VesselCollection) Update(ctx context.Context, id string, entity Vessel) error {
	vessel := Vessel{
		Tenant:               entity.Tenant,
		VesselSpecifications: entity.VesselSpecifications,
		CreatedAt:            entity.CreatedAt,
		UpdatedAt:            time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, vessel)
}

func (r *VesselCollection) GetAll(ctx context.Context, tenant string) ([]VesselManagementResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *VesselCollection) GetByID(ctx context.Context, id string, tenant string) (VesselManagementResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *VesselCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
