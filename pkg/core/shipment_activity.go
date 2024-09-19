package core

import (
	"time"
)

type CargoOperationsActivity struct {
	ActivityType                string                      `json:"activity_type"`
	CustomerName                string                      `json:"customer_name"`
	AnchorageLocation           string                      `json:"anchorage_location"`
	TerminalName                string                      `json:"terminal_name"`
	ShipmentProduct             []*ShipmentProduct          `json:"shipment_product"`
	Readiness                   time.Time                   `json:"readiness"`
	ETB                         time.Time                   `json:"etb"`
	ETD                         time.Time                   `json:"etd"`
	ArrivalDepartureInformation ArrivalDepartureInformation `json:"arrival_departure_information"`
}

type BunkeringActivity struct {
	CustomerName               string                        `json:"customer_name"`
	Supplier                   string                        `json:"supplier"`
	SupplierContact            string                        `json:"supplier_contact"`
	AppointedSurveyor          string                        `json:"appointed_surveyor"`
	Docking                    string                        `json:"docking"`
	SupplierVessel             string                        `json:"supplier_vessel"`
	BunkerIntakeSpecifications []*BunkerIntakeSpecifications `json:"bunker_intake_specifications"`
	ShipmentProduct            []*ShipmentProduct            `json:"shipment_product"`
	Freeboard                  float64                       `json:"freeboard"`
	Readiness                  time.Time                     `json:"readiness"`
	ETB                        time.Time                     `json:"etb"`
	ETD                        time.Time                     `json:"etd"`
}

type BunkerIntakeSpecifications struct {
	SubProductType        string  `json:"sub_product_type"`
	MaximumQuantityIntake float64 `json:"maximum_quantity_intake"`
	MaximumHoseSize       float64 `json:"maximum_hose_size"`
}

type Activity struct {
	ActivityType                string                      `json:"activity_type"`
	CustomerName                string                      `json:"customer_name"`
	AnchorageLocation           string                      `json:"anchorage_location"`
	TerminalName                string                      `json:"terminal_name"`
	ShipmentProduct             []*ShipmentProduct          `json:"shipment_product"`
	Readiness                   time.Time                   `json:"readiness"`
	ETB                         time.Time                   `json:"etb"`
	ETD                         time.Time                   `json:"etd"`
	ArrivalDepartureInformation ArrivalDepartureInformation `json:"arrival_departure_information"`
}

type CustomerSpecifications struct {
	Customer string `json:"customer"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
}

type TerminalSpecifications struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
}

type SupplierSpecifications struct {
	Name    string   `json:"name"`
	Vessels []string `json:"vessels"`
	Email   string   `json:"email"`
	Contact string   `json:"contact"`
}

type ShipmentProduct struct {
	SubProductType string `json:"sub_product_type"`
	Quantity       int64  `json:"quantity"`
	QuantityCode   string `json:"quantity_code"`
	Percentage     int64  `json:"percentage"`
}

type ArrivalDepartureInformation struct {
	ArrivalDisplacement   float64 `json:"arrival_displacement"`
	DepartureDisplacement float64 `json:"departure_displacement"`
	ArrivalDraft          float64 `json:"arrival_draft"`
	DepartureDraft        float64 `json:"departure_draft"`
	ArrivalMastHeight     float64 `json:"arrival_mast_height"`
	DepartureMastHeight   float64 `json:"departure_mast_height"`
}
