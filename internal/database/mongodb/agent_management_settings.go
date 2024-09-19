package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Agent struct {
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Contact   string     `json:"contact"`
	Tenant    string     `json:"tenant"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type AgentManagementResponse struct {
	ID        string    `bson:"_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Contact   string    `json:"contact"`
	Tenant    string    `json:"tenant"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgentCollection struct {
	*GenericCollection[Agent, AgentManagementResponse]
}

func NewAgentCollection(collection *mongo.Collection) *AgentCollection {
	return &AgentCollection{
		GenericCollection: NewGenericCollection[Agent, AgentManagementResponse](collection),
	}
}

var _ Collection[Agent, AgentManagementResponse] = (*AgentCollection)(nil)

func (r *AgentCollection) Create(ctx context.Context, entity Agent) (string, error) {
	now := time.Now()
	agent := Agent{
		Name:      entity.Name,
		Email:     entity.Email,
		Contact:   entity.Contact,
		Tenant:    entity.Tenant,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	return r.GenericCollection.Create(ctx, agent)
}

func (r *AgentCollection) Update(ctx context.Context, id string, entity Agent) error {
	now := time.Now()

	agent := Agent{
		Name:      entity.Name,
		Email:     entity.Email,
		Contact:   entity.Contact,
		Tenant:    entity.Tenant,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: &now,
	}
	return r.GenericCollection.Update(ctx, id, agent)
}

func (r *AgentCollection) GetAll(ctx context.Context, tenant string) ([]AgentManagementResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *AgentCollection) GetByID(ctx context.Context, id string, tenant string) (AgentManagementResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *AgentCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
