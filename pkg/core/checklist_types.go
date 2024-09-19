package core

type PortDues struct {
	ChecklistInformation
}

type Pilotage struct {
	ChecklistInformation
}

type ServiceLaunch struct {
	ChecklistInformation
}

type Logistics struct {
	ChecklistInformation
}

type HotelCharges struct {
	ChecklistInformation
}

type AirTickets struct {
	ChecklistInformation
}

type TransportCharges struct {
	ChecklistInformation
}

type MedicineSupplies struct {
	ChecklistInformation
}

type FreshWaterSupply struct {
	ChecklistInformation
}

type MarineAdvisory struct {
	ChecklistInformation
}

type CourierServices struct {
	ChecklistInformation
}

type CrossHarbourFees struct {
	ChecklistInformation
}

type SupplyBoat struct {
	ChecklistInformation
}

type Repairs struct {
	Deslopping Deslopping `json:"deslopping"`
	LiftRepair LiftRepair `json:"lift_repair"`
	UWClean    UWClean    `json:"uw_clean"`
}

type CrewChange struct {
	SignOn  []string `json:"sign_on"`
	SignOff []string `json:"sign_off"`
}

type ExtrasInformation struct {
	Name            string `json:"name"`
	Supplier        string `json:"supplier"`
	ServiceProvided bool   `json:"service_provided"`
}

type Deslopping struct {
	ChecklistInformation
}

type LiftRepair struct {
	ChecklistInformation
}

type UWClean struct {
	ChecklistInformation
}

type ChecklistInformation struct {
	Supplier        string `json:"supplier"`
	ServiceProvided bool   `json:"service_provided"`
}
