package handler

import (
	"backend-crm/internal/clients"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/enum"
	"context"
	"log"
	"strconv"
	"sync"
	"time"
)

type ShipmentReminders struct {
	Within24Hours bool
	Within12Hours bool
	Within6Hours  bool
	AGDReminder   bool
	NORReminder   bool
}

// [
//
//	{
//		"shipment1": {
//		Within24Hours: true,
//		Within12Hours: false,
//		Within6Hours:  false,
//		AGDReminder:  false,
//		NORReminder:  false,
//		},
//		"shipment2": {
//		Within24Hours: true,
//		Within12Hours: true,
//		Within6Hours:  true,
//		AGDReminder:  false,
//		NORReminder:  false,
//		}
//	}
//
// ]
var sentReminders = make(map[string]*ShipmentReminders)

func StartTimer(stopChan chan struct{}, wg *sync.WaitGroup, shipmentCollection *database.ShipmentCollection, interval time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		log.Println("timer log")
		select {
		case <-ticker.C:
			// Fetch all shipments and process ETA reminders
			// can put this in a cache next time to avoid hitting db everytime
			shipments, err := shipmentCollection.GetAll(context.Background(), "")
			if err != nil {
				log.Printf("Error fetching shipments: %v", err)
				continue
			}

			// Process the shipments
			CheckETAsAndSendWhatsApp(shipments)
			// SendAGDs(shipments)
		case <-stopChan:
			log.Println("Timer stopped")
			return
		}
	}
}

func CheckETAsAndSendWhatsApp(shipments []database.ShipmentResponse) {
	currentTime := time.Now()

	bufferTimeForMessages := 15 * time.Minute

	for _, shipment := range shipments {
		log.Println(currentTime, "curr")
		log.Println(shipment.CurrentETA, "ETA")
		log.Println(shipment.CurrentETA.Sub(currentTime))
		hoursBeforeETA := formatDuration(shipment.CurrentETA.Sub(currentTime))
		if shipment.CurrentStatus == enum.SHIPMENT_STATUS_NOR_TENDERED && !hasSentReminder(shipment.ID, "NOR") {
			log.Printf("Shipment %s status is NOR Tendered \n", shipment.ID)
			SendNORTenderedWhatsApp(shipment)
			markReminderSent(shipment.ID, "NOR")
		}

		if shipment.CurrentETA.Sub(currentTime) > 0 {
			switch {
			case shipment.CurrentETA.Sub(currentTime) <= 5*time.Minute && !hasSentReminder(shipment.ID, "AGD"):
				log.Printf("Shipment %s passed ETA, sends AGD reminder\n", shipment.ID)
				SendAGDWhatsApp(shipment)
				markReminderSent(shipment.ID, "AGD")
				// delete entry from map to avoid redundant memory usage
				// delete(sentReminders, shipment.ID)
				// log.Printf("Deleted reminder for shipment %s as all reminders have been sent\n", shipment.ID)

			case shipment.CurrentETA.Sub(currentTime) <= 6*time.Hour+bufferTimeForMessages && !hasSentReminder(shipment.ID, "6"):
				log.Printf("Shipment %s is within 6 hours of ETA\n", shipment.ID)
				SendETAWhatsApp(shipment, hoursBeforeETA)
				markReminderSent(shipment.ID, "6")

			case shipment.CurrentETA.Sub(currentTime) <= 12*time.Hour+bufferTimeForMessages && !hasSentReminder(shipment.ID, "12"):
				log.Printf("Shipment %s is within 12 hours of ETA\n", shipment.ID)
				SendETAWhatsApp(shipment, hoursBeforeETA)
				markReminderSent(shipment.ID, "12")

			case shipment.CurrentETA.Sub(currentTime) <= 24*time.Hour+bufferTimeForMessages && !hasSentReminder(shipment.ID, "24"):
				log.Printf("Shipment %s is within 24 hours of ETA\n", shipment.ID)
				SendETAWhatsApp(shipment, hoursBeforeETA)
				markReminderSent(shipment.ID, "24")
			}
			// put the delete item from hashmap logic to AFTER shipment is marked as completed!
			// In the future, we have may
			// some logic which states that if a particular status such as NOR tendered is stuck for too long
			// we want to trigger a second notification, hence do not rush to delete the item

			// } else {
			// 	if hasSentReminder(shipment.ID, "AGD") {
			// 		// delete entry from map to avoid redundant memory usage
			// 		delete(sentReminders, shipment.ID)
			// 		log.Printf("Deleted reminder for shipment %s as all reminders have been sent\n", shipment.ID)
			// 	}

		}
	}
}

// SendETAWhatsApp sends a WhatsApp message to remind about the shipment ETA.
func SendETAWhatsApp(shipment database.ShipmentResponse, hoursBeforeETA string) {
	_, err := clients.WhatsAppClient.RemindShipmentETAWhatsAppMessage(
		shipment.VesselSpecifications.VesselName,
		strconv.FormatInt(shipment.VesselSpecifications.ImoNumber, 10), // Convert int64 to string
		hoursBeforeETA,
		shipment.ShipmentDetails.Agent.Contact,
	)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
	}
}

// Send sends a WhatsApp message to remind about the AGD (Arrival General Declaration).
func SendAGDWhatsApp(shipment database.ShipmentResponse) {
	_, err := clients.WhatsAppClient.RemindAGDWhatsAppMessage(
		shipment.ShipmentDetails.Agent.Name,
		shipment.VesselSpecifications.VesselName,
		strconv.FormatInt(shipment.VesselSpecifications.ImoNumber, 10), // Convert int64 to string
		shipment.CurrentETA.Local().Format("02-Jan-2006"),
		shipment.CurrentETA.Local().Format("15:04"),
		shipment.ShipmentDetails.Agent.Contact,
	)
	if err != nil {
		log.Printf("Error sending AGD WhatsApp message: %v", err)
	}
}

// SendNORTenderedWhatsApp sends a WhatsApp message to remind about the shipment ETA.
func SendNORTenderedWhatsApp(shipment database.ShipmentResponse) {
	_, err := clients.WhatsAppClient.RemindNORTenderedWhatsAppMessage(
		shipment.ShipmentDetails.Agent.Name,
		shipment.VesselSpecifications.VesselName,
		strconv.FormatInt(shipment.VesselSpecifications.ImoNumber, 10), // Convert int64 to string
		shipment.ShipmentDetails.Agent.Contact,
	)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
	}
}

// formatDuration converts a time.Duration to a string in hours and minutes.
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	return strconv.Itoa(hours)

}

// markReminderSent marks that a reminder has been sent for the given shipment and time interval.
func markReminderSent(shipmentID string, key string) {
	if _, exists := sentReminders[shipmentID]; !exists {
		sentReminders[shipmentID] = &ShipmentReminders{}
	}
	switch key {
	case "24":
		sentReminders[shipmentID].Within24Hours = true
	case "12":
		sentReminders[shipmentID].Within12Hours = true
	case "6":
		sentReminders[shipmentID].Within6Hours = true
	case "AGD":
		sentReminders[shipmentID].AGDReminder = true
	case "NOR":
		sentReminders[shipmentID].NORReminder = true
	}

}

// hasSentReminder checks if a reminder has been sent for the given shipment and time interval.
func hasSentReminder(shipmentID string, key string) bool {
	if reminders, exists := sentReminders[shipmentID]; exists {
		switch key {
		case "24":
			return reminders.Within24Hours
		case "12":
			return reminders.Within12Hours
		case "6":
			return reminders.Within6Hours
		case "AGD":
			return reminders.AGDReminder
		case "NOR":
			return reminders.NORReminder

		}

	}
	return false
}
