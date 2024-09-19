package handler

import (
	"backend-crm/internal/clients"
	"backend-crm/internal/config"
	database "backend-crm/internal/database/mongodb"
	"backend-crm/nlp"
	"backend-crm/pkg/auth"
	"backend-crm/pkg/core"
	"backend-crm/pkg/enum"
	"backend-crm/pkg/external/openai"
	"backend-crm/pkg/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedHandler struct {
	FeedEmailCollection database.Collection[database.FeedEmail, database.FeedEmailResponse]
	ShipmentCollection  database.Collection[database.Shipment, database.ShipmentResponse]
	ChecklistCollection database.Collection[database.Checklist, database.ChecklistResponse]
}

func NewFeedHandler(
	feedEmailCollection database.Collection[database.FeedEmail, database.FeedEmailResponse],
	shipmentCollection database.Collection[database.Shipment, database.ShipmentResponse],
	checklistCollection database.Collection[database.Checklist, database.ChecklistResponse],
) *FeedHandler {
	return &FeedHandler{
		FeedEmailCollection: feedEmailCollection,
		ShipmentCollection:  shipmentCollection,
		ChecklistCollection: checklistCollection,
	}
}

type getAllFeedEmailsResponse struct {
	FeedEmails []database.FeedEmailResponse `json:"feed_emails"`
}

func (h *FeedHandler) GetAllFeedEmails(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	var feedEmailsList []database.FeedEmailResponse

	feedEmailsList, err := h.FeedEmailCollection.GetAll(r.Context(), tenant)
	log.Println(feedEmailsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var feedEmails []database.FeedEmailResponse
	for _, feedEmail := range feedEmailsList {
		feedEmails = append(feedEmails, database.FeedEmailResponse{
			ID:               feedEmail.ID,
			Tenant:           feedEmail.Tenant,
			MasterEmail:      feedEmail.MasterEmail,
			ReceivedDateTime: feedEmail.ReceivedDateTime,
			ToEmailAddress:   feedEmail.ToEmailAddress,
			Subject:          feedEmail.Subject,
			BodyContent:      feedEmail.BodyContent,
			ShipmentId:       feedEmail.ShipmentId,
		})
	}

	// Creating the response object
	response := getAllFeedEmailsResponse{FeedEmails: feedEmails}
	render.JSON(w, r, response)
}

func (h *FeedHandler) GetFeedEmailsByMasterEmail(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	masterEmail := r.URL.Query().Get("master_email")
	var feedEmailsList []database.FeedEmailResponse

	feedEmailsList, err := h.FeedEmailCollection.GetAllByKeyValue(r.Context(), "masteremail", masterEmail, tenant)

	log.Println("List of Feed Emails:", feedEmailsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var feedEmails []database.FeedEmailResponse
	for _, feedEmail := range feedEmailsList {
		feedEmails = append(feedEmails, database.FeedEmailResponse{
			ID:               feedEmail.ID,
			Tenant:           feedEmail.Tenant,
			MasterEmail:      feedEmail.MasterEmail,
			ReceivedDateTime: feedEmail.ReceivedDateTime,
			ToEmailAddress:   feedEmail.ToEmailAddress,
			Subject:          feedEmail.Subject,
			BodyContent:      feedEmail.BodyContent,
			ShipmentId:       feedEmail.ShipmentId,
		})
	}

	// Creating the response object
	response := getAllFeedEmailsResponse{FeedEmails: feedEmails}
	render.JSON(w, r, response)
}

func (h *FeedHandler) GetFeedEmailsByShipmentId(w http.ResponseWriter, r *http.Request) {
	email := auth.GetEmailFromToken(r.Context())
	tenant := config.MakeMapping(email)
	shipmentId := chi.URLParam(r, "shipment_id")
	log.Printf("Retrieved shipment ID from URL: %s", shipmentId) // Confirm shipmentId retrieval
	var feedEmailsList []database.FeedEmailResponse

	feedEmailsList, err := h.FeedEmailCollection.GetAllByKeyValue(r.Context(), "shipmentid", shipmentId, tenant)

	log.Println(feedEmailsList)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	var feedEmails []database.FeedEmailResponse
	for _, feedEmail := range feedEmailsList {
		feedEmails = append(feedEmails, database.FeedEmailResponse{
			ID:               feedEmail.ID,
			Tenant:           feedEmail.Tenant,
			MasterEmail:      feedEmail.MasterEmail,
			ReceivedDateTime: feedEmail.ReceivedDateTime,
			ToEmailAddress:   feedEmail.ToEmailAddress,
			Subject:          feedEmail.Subject,
			BodyContent:      feedEmail.BodyContent,
			ShipmentId:       feedEmail.ShipmentId,
		})
	}

	// Creating the response object
	response := getAllFeedEmailsResponse{FeedEmails: feedEmails}
	render.JSON(w, r, response)
}

func (h *FeedHandler) CreateFeedMessage(w http.ResponseWriter, r *http.Request) {
	var createFeedEmailParams database.FeedEmail
	log.Println(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&createFeedEmailParams); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	// Get tenant directly from the "to_email_address"
	tenant := config.MakeMapping(createFeedEmailParams.ToEmailAddress)
	createFeedEmailParams.Tenant = tenant

	updated, err := h.updateChecklistBasedOnEmailContent(w, r, createFeedEmailParams, tenant)
	if err != nil {
		log.Println(err)
		return
	}

	if !updated {
		// Get shipmentId according to master email address + status not completed
		//    - there should only be one shipment+masteremail that is currently in progress

		var shipmentsList []database.ShipmentResponse
		shipmentsList, err := h.ShipmentCollection.GetAll(r.Context(), tenant)
		log.Println(shipmentsList)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				render.Render(w, r, ErrNotFound)
			} else {
				render.Render(w, r, ErrInternalServerError)
			}
			return
		}
		for _, shipment := range shipmentsList {
			if shipment.MasterEmail == createFeedEmailParams.MasterEmail && (shipment.CurrentStatus != enum.SHIPMENT_STATUS_COSP && shipment.CurrentStatus != enum.SHIPMENT_STATUS_ACTIVITY_COMPLETED) {
				createFeedEmailParams.ShipmentId = shipment.ID
			}
		}

		log.Println(createFeedEmailParams)
		if createFeedEmailParams.ShipmentId == "" {
			// if no corresponding shipment meeting criteria, just skip
			// WIP 2 Aug 2024 - need to add a whatsapp reminder to the team lead to create shipment
			render.Render(w, r, SuccessCreated)
			return
		} else {
			_, err = h.FeedEmailCollection.Create(r.Context(), createFeedEmailParams)
			if err != nil {
				if errors.Is(err, database.ErrDuplicateKey) {
					render.Render(w, r, ErrDuplicate(err))
					return
				}
				render.Render(w, r, ErrInternalServerError)
				return
			}

			// now we update the shipment collection
			currentShipmentResponse, invalid := h.ShipmentCollection.GetByID(r.Context(), createFeedEmailParams.ShipmentId, tenant)
			if invalid != nil {
				if invalid == mongo.ErrNoDocuments {
					render.Render(w, r, ErrNotFound)
				} else {
					render.Render(w, r, ErrInternalServerError)
				}
				return
			}

			// get the ETA and status from Master Email Content

			_, newStatus := nlp.TranslateMasterEmailContentToStatusUpdate(createFeedEmailParams.Subject, createFeedEmailParams.BodyContent, currentShipmentResponse)
			// use OpenAI to extract ETA
			ETAPrompt := nlp.GetETAFromMasterEmailOpenAIPrompt()

			// check openAI if the email contains any information that should update the checklist
			parsedETA, invalid := clients.OpenAIClient.ExtractEntityFromText(openai.OpenAIRequest{
				Model: "gpt-4o-mini",
				Messages: []openai.Message{
					{
						Role:    "user",
						Content: createFeedEmailParams.BodyContent + ETAPrompt,
					},
				},
			})

			if invalid != nil {
				log.Println(invalid)
				return
			}

			log.Println("openai", parsedETA)

			systemTimeNow := time.Now()
			systemTimeNowLocation := systemTimeNow.Location()
			// Check if the location is not "Local"
			if systemTimeNowLocation.String() != "Local" {
				// Define a fixed time zone with +8 hours offset from UTC
				systemTimeNowLocation = time.FixedZone("UTC+8", 8*60*60)

			}
			// Define the time format that matches the ETAPrompt
			layout := "02 Jan 2006 15:04"
			// Parse the ETAPrompt string into a time.Time object
			newETA, err := time.ParseInLocation(layout, parsedETA, systemTimeNowLocation)

			log.Println(newETA, "prasedTime")
			if err != nil {
				log.Println("Error parsing ETA:", err)
				return
			}
			log.Println("Parsed Time (UTC+8):", newETA)

			// Convert the parsed time to the specified time zone
			// newETA := parsedTime.In(systemTimeNowLocation)

			// log.Println(newETA, "eta from chatgpt")
			log.Println(currentShipmentResponse.CreatedAt, "CREATEDAT")
			newShipment := database.Shipment{
				Tenant:               currentShipmentResponse.Tenant,
				MasterEmail:          currentShipmentResponse.MasterEmail,
				InitialETA:           currentShipmentResponse.InitialETA,
				CurrentETA:           newETA,
				VoyageNumber:         currentShipmentResponse.VoyageNumber,
				CurrentStatus:        newStatus,
				ShipmentDetails:      currentShipmentResponse.ShipmentDetails,
				VesselSpecifications: currentShipmentResponse.VesselSpecifications,
				ShipmentType:         currentShipmentResponse.ShipmentType,
				CreatedAt:            currentShipmentResponse.CreatedAt,
				UpdatedAt:            time.Now(),
			}

			invalid = h.ShipmentCollection.Update(r.Context(), createFeedEmailParams.ShipmentId, newShipment)
			if invalid != nil {
				if invalid == mongo.ErrNoDocuments {
					render.Render(w, r, ErrNotFound)
				} else {
					render.Render(w, r, ErrInternalServerError)
				}
				return
			}
			log.Println("shipment updated success after master email")
		}
	}

	w.WriteHeader(201)
	w.Write(nil)

	render.Render(w, r, SuccessCreated)
}

func (h *FeedHandler) updateChecklistBasedOnEmailContent(w http.ResponseWriter, r *http.Request, createFeedEmailParams database.FeedEmail, tenant string) (bool, error) {
	prompt := nlp.GetIntentionsFromChecklistOpenAIPrompt()

	// check openAI if the email contains any information that should update the checklist
	parsedIntention, invalid := clients.OpenAIClient.ExtractEntityFromText(openai.OpenAIRequest{
		Model: "gpt-4o-mini",
		Messages: []openai.Message{
			{
				Role:    "user",
				Content: createFeedEmailParams.BodyContent + prompt,
			},
		},
	})

	if invalid != nil {
		log.Println(invalid)
		return false, invalid
	}

	log.Println("openai", parsedIntention)

	var data []interface{}
	err := json.Unmarshal([]byte(parsedIntention), &data)
	if err != nil {
		log.Println("Error unmarshalling JSON: %v", err)
	}

	if data[0] != "intention not found" {
		var shipmentsList []database.ShipmentResponse
		var checklistShipmentID string

		// Find matching shipment based on vessel name in email subject
		shipmentsList, err := h.ShipmentCollection.GetAll(r.Context(), tenant)
		if err != nil {
			log.Println("Error getting shipments:", err)
			render.Render(w, r, ErrInternalServerError)
			return false, invalid
		}
		for _, shipment := range shipmentsList {

			normalizedSubject := utils.NormalizeString(createFeedEmailParams.Subject)
			normalizedVesselName := utils.NormalizeString(shipment.VesselSpecifications.VesselName)
			log.Println("enter")
			if strings.Contains(normalizedSubject, normalizedVesselName) &&
				(shipment.CurrentStatus != enum.SHIPMENT_STATUS_COSP && shipment.CurrentStatus != enum.SHIPMENT_STATUS_ACTIVITY_COMPLETED) {
				log.Println("enter2")

				checklistShipmentID = shipment.ID
				break
			}
		}

		if checklistShipmentID == "" {
			log.Println("No matching shipment found")
			render.Render(w, r, ErrNotFound)
			return false, invalid
		}

		// Get current checklist
		currentChecklistResponses, err := h.ChecklistCollection.GetAllByKeyValue(r.Context(), "shipmentid", checklistShipmentID, tenant)
		if err != nil || len(currentChecklistResponses) == 0 {
			log.Println("Error getting checklist:", err)
			render.Render(w, r, ErrNotFound)
			return false, invalid
		}
		currentChecklistResponse := currentChecklistResponses[0]
		log.Println(currentChecklistResponse, "currentChecklistResponse")
		// var updatedChecklist database.Checklist

		updatedChecklist := database.Checklist{
			PortDues:         currentChecklistResponse.PortDues,
			Pilotage:         currentChecklistResponse.Pilotage,
			ServiceLaunch:    currentChecklistResponse.ServiceLaunch,
			Logistics:        currentChecklistResponse.Logistics,
			HotelCharges:     currentChecklistResponse.HotelCharges,
			AirTickets:       currentChecklistResponse.AirTickets,
			TransportCharges: currentChecklistResponse.TransportCharges,
			MedicineSupplies: currentChecklistResponse.MedicineSupplies,
			FreshWaterSupply: currentChecklistResponse.FreshWaterSupply,
			MarineAdvisory:   currentChecklistResponse.MarineAdvisory,
			CourierServices:  currentChecklistResponse.CourierServices,
			CrossHarbourFees: currentChecklistResponse.CrossHarbourFees,
			SupplyBoat:       currentChecklistResponse.SupplyBoat,
			Repairs:          currentChecklistResponse.Repairs,
			Extras:           currentChecklistResponse.Extras,
			CrewChange:       currentChecklistResponse.CrewChange,
			ShipmentID:       currentChecklistResponse.ShipmentID,
			Tenant:           currentChecklistResponse.Tenant,
			CreatedAt:        currentChecklistResponse.CreatedAt,
			UpdatedAt:        time.Now(),
		}

		if len(data) > 0 {
			switch data[0] {
			case "crew_change":
				if len(data) == 3 { // must be 3
					log.Println(data, "data")
					signOn, ok := data[1].([]interface{})
					if !ok {
						log.Println("Error parsing sign on data")
					}

					signOff, ok := data[2].([]interface{})
					if !ok {
						log.Println("Error parsing sign off data")
					}

					var signOnList []string
					for _, name := range signOn {
						signOnList = append(signOnList, name.(string))
					}

					var signOffList []string
					for _, name := range signOff {
						signOffList = append(signOffList, name.(string))
					}

					updatedChecklist.CrewChange = core.CrewChange{
						SignOn:  signOnList,
						SignOff: signOffList,
					}
				} else {
					log.Println("Invalid crew change data format")
				}

			default:
				intentionsList := data

				for _, intention := range intentionsList {
					intentionStr, ok := intention.(string)
					if !ok {
						log.Println("Error parsing intention")
					}

					// Set the service provided flag based on intention
					switch intentionStr {
					case "logistics":
						if !updatedChecklist.Logistics.ChecklistInformation.ServiceProvided {
							updatedChecklist.Logistics.ChecklistInformation.ServiceProvided = true
						}
					case "hotel_charges":
						if !updatedChecklist.Logistics.ChecklistInformation.ServiceProvided {
							updatedChecklist.Logistics.ChecklistInformation.ServiceProvided = true
						}
					case "air_tickets":
						if !updatedChecklist.AirTickets.ChecklistInformation.ServiceProvided {
							updatedChecklist.AirTickets.ChecklistInformation.ServiceProvided = true
						}
					case "transport_charges":
						if !updatedChecklist.TransportCharges.ChecklistInformation.ServiceProvided {
							updatedChecklist.TransportCharges.ChecklistInformation.ServiceProvided = true
						}
					case "medicine_supplies":
						if !updatedChecklist.MedicineSupplies.ChecklistInformation.ServiceProvided {
							updatedChecklist.MedicineSupplies.ChecklistInformation.ServiceProvided = true
						}
					case "fresh_water_supply":
						if !updatedChecklist.FreshWaterSupply.ChecklistInformation.ServiceProvided {
							updatedChecklist.FreshWaterSupply.ChecklistInformation.ServiceProvided = true
						}
					case "marine_advisory":
						if !updatedChecklist.MarineAdvisory.ChecklistInformation.ServiceProvided {
							updatedChecklist.MarineAdvisory.ChecklistInformation.ServiceProvided = true
						}
					case "courier_services":
						if !updatedChecklist.CourierServices.ChecklistInformation.ServiceProvided {
							updatedChecklist.CourierServices.ChecklistInformation.ServiceProvided = true
						}
					case "cross_harbour_fees":
						if !updatedChecklist.CrossHarbourFees.ChecklistInformation.ServiceProvided {
							updatedChecklist.CrossHarbourFees.ChecklistInformation.ServiceProvided = true
						}
					case "supply_boat":
						if !updatedChecklist.SupplyBoat.ChecklistInformation.ServiceProvided {
							updatedChecklist.SupplyBoat.ChecklistInformation.ServiceProvided = true
						}
					case "deslopping":
						if !updatedChecklist.Repairs.Deslopping.ChecklistInformation.ServiceProvided {
							updatedChecklist.Repairs.Deslopping.ChecklistInformation.ServiceProvided = true
						}
					case "uw_clean":
						if !updatedChecklist.Repairs.UWClean.ChecklistInformation.ServiceProvided {
							updatedChecklist.Repairs.UWClean.ChecklistInformation.ServiceProvided = true
						}
					case "lift_repair":
						if !updatedChecklist.Repairs.LiftRepair.ChecklistInformation.ServiceProvided {
							updatedChecklist.Repairs.LiftRepair.ChecklistInformation.ServiceProvided = true
						}
					}
				}
			}
		}
		log.Println(updatedChecklist)
		invalid = h.ChecklistCollection.Update(r.Context(), currentChecklistResponse.ID, updatedChecklist)
		if invalid != nil {
			if invalid == mongo.ErrNoDocuments {
				render.Render(w, r, ErrNotFound)
			} else {
				render.Render(w, r, ErrInternalServerError)
			}
			return false, invalid
		}
		log.Println("checklist updated success after parsedIntention")

		// appends the email in feed as well
		createFeedEmailParams.ShipmentId = checklistShipmentID

		_, err = h.FeedEmailCollection.Create(r.Context(), createFeedEmailParams)
		if err != nil {
			if errors.Is(err, database.ErrDuplicateKey) {
				render.Render(w, r, ErrDuplicate(err))
				return false, err
			}
			render.Render(w, r, ErrInternalServerError)
			return false, err
		}
		log.Println("feed updated success after parsedIntention")

		return true, nil
	}
	log.Println("detected intention not found, hence moving on")
	return false, nil
}
