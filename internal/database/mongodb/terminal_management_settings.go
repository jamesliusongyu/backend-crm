package database

import (
	"backend-crm/pkg/core"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Terminal struct {
	Tenant                 string                      `json:"tenant"`
	TerminalSpecifications core.TerminalSpecifications `json:"terminal_specifications"`
	CreatedAt              time.Time                   `json:"created_at"`
	UpdatedAt              time.Time                   `json:"updated_at"`
}

type TerminalManagementResponse struct {
	ID                     string                      `bson:"_id"`
	Tenant                 string                      `json:"tenant"`
	TerminalSpecifications core.TerminalSpecifications `json:"terminal_specifications"`
	CreatedAt              time.Time                   `json:"created_at"`
	UpdatedAt              time.Time                   `json:"updated_at"`
}

type TerminalCollection struct {
	*GenericCollection[Terminal, TerminalManagementResponse]
}

func NewTerminalCollection(collection *mongo.Collection) *TerminalCollection {
	return &TerminalCollection{
		GenericCollection: NewGenericCollection[Terminal, TerminalManagementResponse](collection),
	}
}

var _ Collection[Terminal, TerminalManagementResponse] = (*TerminalCollection)(nil)

func (r *TerminalCollection) Create(ctx context.Context, entity Terminal) (string, error) {

	Terminal := Terminal{
		Tenant:                 entity.Tenant,
		TerminalSpecifications: entity.TerminalSpecifications,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	return r.GenericCollection.Create(ctx, Terminal)
}

func (r *TerminalCollection) Update(ctx context.Context, id string, entity Terminal) error {
	Terminal := Terminal{
		Tenant:                 entity.Tenant,
		TerminalSpecifications: entity.TerminalSpecifications,
		CreatedAt:              entity.CreatedAt,
		UpdatedAt:              time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, Terminal)
}

func (r *TerminalCollection) GetAll(ctx context.Context, tenant string) ([]TerminalManagementResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *TerminalCollection) GetByID(ctx context.Context, id string, tenant string) (TerminalManagementResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *TerminalCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
