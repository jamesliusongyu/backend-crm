package database

import (
	"backend-crm/pkg/core"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Checklist struct {
	PortDues         core.PortDues                      `json:"port_dues"`
	Pilotage         core.Pilotage                      `json:"pilotage"`
	ServiceLaunch    core.ServiceLaunch                 `json:"service_launch"`
	Logistics        core.Logistics                     `json:"logistics"`
	HotelCharges     core.HotelCharges                  `json:"hotel_charges"`
	AirTickets       core.AirTickets                    `json:"air_tickets"`
	TransportCharges core.TransportCharges              `json:"transport_charges"`
	MedicineSupplies core.MedicineSupplies              `json:"medicine_supplies"`
	FreshWaterSupply core.FreshWaterSupply              `json:"fresh_water_supply"`
	MarineAdvisory   core.MarineAdvisory                `json:"marine_advisory"`
	CourierServices  core.CourierServices               `json:"courier_services"`
	CrossHarbourFees core.CrossHarbourFees              `json:"cross_harbour_fees"`
	SupplyBoat       core.SupplyBoat                    `json:"supply_boat"`
	Repairs          core.Repairs                       `json:"repairs"`
	CrewChange       core.CrewChange                    `json:"crew_change"`
	Extras           map[string]*core.ExtrasInformation `json:"extras"`

	ShipmentID string    `json:"shipment_id"`
	Tenant     string    `json:"tenant"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ChecklistResponse struct {
	ID               string                             `bson:"_id"`
	PortDues         core.PortDues                      `json:"port_dues"`
	Pilotage         core.Pilotage                      `json:"pilotage"`
	ServiceLaunch    core.ServiceLaunch                 `json:"service_launch"`
	Logistics        core.Logistics                     `json:"logistics"`
	HotelCharges     core.HotelCharges                  `json:"hotel_charges"`
	AirTickets       core.AirTickets                    `json:"air_tickets"`
	TransportCharges core.TransportCharges              `json:"transport_charges"`
	MedicineSupplies core.MedicineSupplies              `json:"medicine_supplies"`
	FreshWaterSupply core.FreshWaterSupply              `json:"fresh_water_supply"`
	MarineAdvisory   core.MarineAdvisory                `json:"marine_advisory"`
	CourierServices  core.CourierServices               `json:"courier_services"`
	CrossHarbourFees core.CrossHarbourFees              `json:"cross_harbour_fees"`
	SupplyBoat       core.SupplyBoat                    `json:"supply_boat"`
	Repairs          core.Repairs                       `json:"repairs"`
	CrewChange       core.CrewChange                    `json:"crew_change"`
	Extras           map[string]*core.ExtrasInformation `json:"extras"`

	ShipmentID string    `json:"shipment_id"`
	Tenant     string    `json:"tenant"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ChecklistCollection struct {
	*GenericCollection[Checklist, ChecklistResponse]
}

func NewChecklistCollection(collection *mongo.Collection) *ChecklistCollection {
	return &ChecklistCollection{
		GenericCollection: NewGenericCollection[Checklist, ChecklistResponse](collection),
	}
}

var _ Collection[Checklist, ChecklistResponse] = (*ChecklistCollection)(nil)

func (r *ChecklistCollection) Create(ctx context.Context, entity Checklist) (string, error) {

	Checklist := Checklist{
		PortDues:         entity.PortDues,
		Pilotage:         entity.Pilotage,
		ServiceLaunch:    entity.ServiceLaunch,
		Logistics:        entity.Logistics,
		HotelCharges:     entity.HotelCharges,
		AirTickets:       entity.AirTickets,
		TransportCharges: entity.TransportCharges,
		MedicineSupplies: entity.MedicineSupplies,
		FreshWaterSupply: entity.FreshWaterSupply,
		MarineAdvisory:   entity.MarineAdvisory,
		CourierServices:  entity.CourierServices,
		CrossHarbourFees: entity.CrossHarbourFees,
		SupplyBoat:       entity.SupplyBoat,
		Repairs:          entity.Repairs,
		Tenant:           entity.Tenant,
		Extras:           entity.Extras,
		CrewChange:       entity.CrewChange,
		ShipmentID:       entity.ShipmentID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return r.GenericCollection.Create(ctx, Checklist)
}

func (r *ChecklistCollection) Update(ctx context.Context, id string, entity Checklist) error {
	Checklist := Checklist{
		PortDues:         entity.PortDues,
		Pilotage:         entity.Pilotage,
		ServiceLaunch:    entity.ServiceLaunch,
		Logistics:        entity.Logistics,
		HotelCharges:     entity.HotelCharges,
		AirTickets:       entity.AirTickets,
		TransportCharges: entity.TransportCharges,
		MedicineSupplies: entity.MedicineSupplies,
		FreshWaterSupply: entity.FreshWaterSupply,
		MarineAdvisory:   entity.MarineAdvisory,
		CourierServices:  entity.CourierServices,
		CrossHarbourFees: entity.CrossHarbourFees,
		SupplyBoat:       entity.SupplyBoat,
		Repairs:          entity.Repairs,
		ShipmentID:       entity.ShipmentID,
		Tenant:           entity.Tenant,
		Extras:           entity.Extras,
		CrewChange:       entity.CrewChange,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: time.Now(),
	}
	return r.GenericCollection.Update(ctx, id, Checklist)
}

func (r *ChecklistCollection) GetAll(ctx context.Context, tenant string) ([]ChecklistResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *ChecklistCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]ChecklistResponse, error) {
	return r.GenericCollection.GetAllByKeyValue(ctx, key, value, tenant)
}

func (r *ChecklistCollection) GetByID(ctx context.Context, id string, tenant string) (ChecklistResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *ChecklistCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
