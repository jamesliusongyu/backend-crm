package database

import (
	"backend-crm/pkg/core"
	"backend-crm/pkg/enum"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Shipment struct {
	Tenant               string               `json:"tenant"`
	MasterEmail          string               `json:"master_email"`
	InitialETA           time.Time            `json:"initial_ETA"`
	CurrentETA           time.Time            `json:"current_ETA"`
	VoyageNumber         string               `json:"voyage_number"`
	CurrentStatus        enum.ShipmentStatus  `json:"current_status"`
	ShipmentDetails      ShipmentDetails      `json:"shipment_details"`
	VesselSpecifications VesselSpecifications `json:"vessel_specifications"`
	ShipmentType         ShipmentType         `json:"shipment_type"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type ShipmentType struct {
	CargoOperations CargoOperations `json:"cargo_operations"`
	Bunkering       Bunkering       `json:"bunkering"`
	OwnerMatters    OwnerMatters    `json:"owner_matters"`
}

type CargoOperations struct {
	CargoOperations         bool                            `json:"cargo_operations"`
	CargoOperationsActivity []*core.CargoOperationsActivity `json:"cargo_operations_activity"`
}

type Bunkering struct {
	Bunkering         bool                      `json:"bunkering"`
	BunkeringActivity []*core.BunkeringActivity `json:"bunkering_activity"`
}

type OwnerMatters struct {
	OwnerMatters bool             `json:"owner_matters"`
	Activity     []*core.Activity `json:"activity"`
}

type VesselSpecifications struct {
	ImoNumber  int64   `json:"imo_number"`
	VesselName string  `json:"vessel_name"`
	CallSign   string  `json:"call_sign"`
	SDWT       int64   `json:"sdwt"`
	NRT        int64   `json:"nrt"`
	Flag       string  `json:"flag"`
	GRT        int64   `json:"grt"`
	LOA        float64 `json:"loa"`
}

type ShipmentDetails struct {
	Agent Agent `json:"agent_details"`
}

type ShipmentResponse struct {
	ID                   string               `bson:"_id"`
	Tenant               string               `json:"tenant"`
	MasterEmail          string               `json:"master_email"`
	InitialETA           time.Time            `json:"initial_ETA"`
	CurrentETA           time.Time            `json:"current_ETA"`
	VoyageNumber         string               `json:"voyage_number"`
	CurrentStatus        enum.ShipmentStatus  `json:"current_status"`
	ShipmentType         ShipmentType         `json:"shipment_type"`
	VesselSpecifications VesselSpecifications `json:"vessel_specifications"`
	ShipmentDetails      ShipmentDetails      `json:"shipment_details"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type ShipmentCollection struct {
	*GenericCollection[Shipment, ShipmentResponse]
}

func NewShipmentCollection(collection *mongo.Collection) *ShipmentCollection {
	return &ShipmentCollection{
		GenericCollection: NewGenericCollection[Shipment, ShipmentResponse](collection),
	}
}

var _ Collection[Shipment, ShipmentResponse] = (*ShipmentCollection)(nil)

func (r *ShipmentCollection) Create(ctx context.Context, entity Shipment) (string, error) {
	shipment := Shipment{
		Tenant:               entity.Tenant,
		MasterEmail:          entity.MasterEmail,
		InitialETA:           entity.InitialETA,
		CurrentETA:           entity.CurrentETA,
		VoyageNumber:         entity.VoyageNumber,
		CurrentStatus:        entity.CurrentStatus,
		ShipmentType:         entity.ShipmentType,
		VesselSpecifications: entity.VesselSpecifications,
		ShipmentDetails:      entity.ShipmentDetails,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	id, err := r.GenericCollection.Create(ctx, shipment)
	return id, err
}

func (r *ShipmentCollection) Update(ctx context.Context, id string, entity Shipment) error {
	shipment := Shipment{
		Tenant:               entity.Tenant,
		MasterEmail:          entity.MasterEmail,
		InitialETA:           entity.InitialETA,
		CurrentETA:           entity.CurrentETA,
		VoyageNumber:         entity.VoyageNumber,
		CurrentStatus:        entity.CurrentStatus,
		ShipmentType:         entity.ShipmentType,
		VesselSpecifications: entity.VesselSpecifications,
		ShipmentDetails:      entity.ShipmentDetails,
		CreatedAt:            entity.CreatedAt,
		UpdatedAt:            time.Now(),
	}
	log.Println(entity)
	return r.GenericCollection.Update(ctx, id, shipment)
}

func (r *ShipmentCollection) GetAll(ctx context.Context, tenant string) ([]ShipmentResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *ShipmentCollection) GetByID(ctx context.Context, id string, tenant string) (ShipmentResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *ShipmentCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
