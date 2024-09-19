package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Customer struct {
	Tenant    string    `json:"tenant"`
	Customer  string    `json:"customer"`
	Company   string    `json:"company"`
	Email     string    `json:"email"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerResponse struct {
	ID        string    `bson:"_id"`
	Tenant    string    `json:"tenant"`
	Customer  string    `json:"customer"`
	Company   string    `json:"company"`
	Email     string    `json:"email"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerCollection struct {
	*GenericCollection[Customer, CustomerResponse]
}

func NewCustomerCollection(collection *mongo.Collection) *CustomerCollection {
	return &CustomerCollection{
		GenericCollection: NewGenericCollection[Customer, CustomerResponse](collection),
	}
}

// Ensure CustomerCollection implements Collection[CustomerResponse]
var _ Collection[Customer, CustomerResponse] = (*CustomerCollection)(nil)

func (r *CustomerCollection) Create(ctx context.Context, entity Customer) (string, error) {
	customer := Customer{
		Tenant:    entity.Tenant,
		Customer:  entity.Customer,
		Company:   entity.Company,
		Email:     entity.Email,
		Contact:   entity.Contact,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return r.GenericCollection.Create(ctx, customer)
}

func (r *CustomerCollection) Update(ctx context.Context, id string, entity Customer) error {
	customer := Customer{
		Tenant:    entity.Tenant,
		Customer:  entity.Customer,
		Company:   entity.Company,
		Email:     entity.Email,
		Contact:   entity.Contact,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, customer)
}

func (r *CustomerCollection) GetAll(ctx context.Context, tenant string) ([]CustomerResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *CustomerCollection) GetByID(ctx context.Context, id string, tenant string) (CustomerResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *CustomerCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
