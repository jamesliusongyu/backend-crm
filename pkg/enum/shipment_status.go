package enum

type ShipmentStatus string

const (
	SHIPMENT_STATUS_NOT_STARTED        ShipmentStatus = "Not Started"
	SHIPMENT_STATUS_EN_ROUTE           ShipmentStatus = "En Route"
	SHIPMENT_STATUS_EOSP               ShipmentStatus = "EOSP"
	SHIPMENT_STATUS_AT_ANCHORAGE       ShipmentStatus = "At Anchorage"
	SHIPMENT_STATUS_NOR_TENDERED       ShipmentStatus = "NOR Tendered"
	SHIPMENT_STATUS_NOR_RETENDERED     ShipmentStatus = "NOR Re-Tendered"
	SHIPMENT_STATUS_BERTHED            ShipmentStatus = "Berthed"
	SHIPMENT_STATUS_ACTIVITY_COMMENCED ShipmentStatus = "Activity Commenced"
	SHIPMENT_STATUS_ACTIVITY_COMPLETED ShipmentStatus = "Activity Completed"
	SHIPMENT_STATUS_COSP               ShipmentStatus = "COSP"
)

// ShipmentStatuses is a slice of all shipment statuses
var ShipmentStatuses = []ShipmentStatus{
	SHIPMENT_STATUS_NOT_STARTED,
	SHIPMENT_STATUS_EN_ROUTE,
	SHIPMENT_STATUS_EOSP,
	SHIPMENT_STATUS_AT_ANCHORAGE,
	SHIPMENT_STATUS_NOR_TENDERED,
	SHIPMENT_STATUS_NOR_RETENDERED,
	SHIPMENT_STATUS_BERTHED,
	SHIPMENT_STATUS_ACTIVITY_COMMENCED,
	SHIPMENT_STATUS_ACTIVITY_COMPLETED,
	SHIPMENT_STATUS_COSP,
}

// ShipmentStatusesWithColours is a map of shipment statuses to their respective colors
var ShipmentStatusesWithColours = map[ShipmentStatus]string{
	SHIPMENT_STATUS_NOT_STARTED:        "lime",
	SHIPMENT_STATUS_EN_ROUTE:           "blue",
	SHIPMENT_STATUS_EOSP:               "green",
	SHIPMENT_STATUS_AT_ANCHORAGE:       "yellow",
	SHIPMENT_STATUS_NOR_TENDERED:       "pink",
	SHIPMENT_STATUS_NOR_RETENDERED:     "purple",
	SHIPMENT_STATUS_BERTHED:            "black",
	SHIPMENT_STATUS_ACTIVITY_COMMENCED: "red",
	SHIPMENT_STATUS_ACTIVITY_COMPLETED: "cyan",
	SHIPMENT_STATUS_COSP:               "orange",
}

// GetShipmentStatuses returns a slice of all shipment statuses
func GetShipmentStatuses() []string {
	statuses := make([]string, len(ShipmentStatuses))
	for i, status := range ShipmentStatuses {
		statuses[i] = string(status)
	}
	return statuses
}

// GetShipmentStatusesWithColours returns a map of shipment statuses to their respective colors
func GetShipmentStatusesWithColours() map[string]string {
	statusesWithColors := make(map[string]string, len(ShipmentStatusesWithColours))
	for status, color := range ShipmentStatusesWithColours {
		statusesWithColors[string(status)] = color
	}
	return statusesWithColors
}
