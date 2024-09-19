package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type InvoicePricing struct {
	Tenant                string            `json:"tenant"`
	ShipmentID            string            `json:"shipment_id"`
	InvoicePricingDetails map[string]string `json:"invoice_pricing_details"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

type InvoicePricingResponse struct {
	ID                    string            `bson:"_id"`
	Tenant                string            `json:"tenant"`
	ShipmentID            string            `json:"shipment_id"`
	InvoicePricingDetails map[string]string `json:"invoice_pricing_details"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

type InvoicePricingCollection struct {
	*GenericCollection[InvoicePricing, InvoicePricingResponse]
}

func NewInvoicePricingCollection(collection *mongo.Collection) *InvoicePricingCollection {
	return &InvoicePricingCollection{
		GenericCollection: NewGenericCollection[InvoicePricing, InvoicePricingResponse](collection),
	}
}

var _ Collection[InvoicePricing, InvoicePricingResponse] = (*InvoicePricingCollection)(nil)

func (r *InvoicePricingCollection) GetByID(ctx context.Context, id string, tenant string) (InvoicePricingResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *InvoicePricingCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (InvoicePricingResponse, error) {
	log.Println(key)
	log.Println(value)
	return r.GenericCollection.GetByKeyValue(ctx, key, value, tenant)
}

func (r *InvoicePricingCollection) Create(ctx context.Context, entity InvoicePricing) (string, error) {
	invoicePricing := InvoicePricing{
		Tenant:                entity.Tenant,
		ShipmentID:            entity.ShipmentID,
		InvoicePricingDetails: entity.InvoicePricingDetails,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
	return r.GenericCollection.Create(ctx, invoicePricing)
}

func (r *InvoicePricingCollection) Update(ctx context.Context, id string, entity InvoicePricing) error {
	log.Println(entity, "enettt")
	invoicePricing := InvoicePricing{
		Tenant:                entity.Tenant,
		ShipmentID:            entity.ShipmentID,
		InvoicePricingDetails: entity.InvoicePricingDetails,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, invoicePricing)
}
