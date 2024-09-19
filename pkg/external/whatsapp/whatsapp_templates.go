package whatsapp

import "log"

// TemplateParameter defines the structure for template parameters
type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// TemplateComponent defines the structure for template components
type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters"`
}

// Template defines the structure for the template
type Template struct {
	Name       string              `json:"name"`
	Language   Language            `json:"language"`
	Components []TemplateComponent `json:"components"`
}

// Language defines the structure for the language
type Language struct {
	Policy string `json:"policy"`
	Code   string `json:"code"`
}

// Payload defines the structure for the message payload
type Payload struct {
	MessagingProduct string   `json:"messaging_product"`
	RecipientType    string   `json:"recipient_type"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Template         Template `json:"template"`
}

// CreateShipmentWhatsAppMessage creates a payload for a shipment WhatsApp message
func CreateShipmentWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime string, destinationContactNumber string) Payload {
	log.Println(destinationContactNumber, "handphone")
	return Payload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               destinationContactNumber, // "6591896934"
		Type:             "template",
		Template: Template{
			// this value is from Whatsapp Cloud
			Name: "columbus_shipment_creation",
			Language: Language{
				Policy: "deterministic",
				Code:   "en_US",
			},
			Components: []TemplateComponent{
				{
					Type: "body",
					Parameters: []TemplateParameter{
						{Type: "text", Text: agentName},
						{Type: "text", Text: vesselName},
						{Type: "text", Text: imoNumber},
						{Type: "text", Text: etaDate},
						{Type: "text", Text: etaTime},
					},
				},
			},
		},
	}
}

// RemindShipmentETAWhatsAppMessage creates a payload for a reminder WhatsApp message
func RemindShipmentETAWhatsAppMessage(vesselName, imoNumber, hoursBeforeETA, destinationContactNumber string) Payload {
	log.Println(destinationContactNumber, "handphone")
	log.Println(hoursBeforeETA, "hoursBeforeETA")
	return Payload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               destinationContactNumber, // "6591896934"
		Type:             "template",
		Template: Template{
			// this value is from Whatsapp Cloud
			Name: "columbus_eta_reminder",
			Language: Language{
				Policy: "deterministic",
				Code:   "en_US",
			},
			Components: []TemplateComponent{
				{
					Type: "body",
					Parameters: []TemplateParameter{
						{Type: "text", Text: hoursBeforeETA},
						{Type: "text", Text: vesselName},
						{Type: "text", Text: imoNumber},
					},
				},
			},
		},
	}
}

// RemindNORTenderedWhatsAppMessage creates a payload for a reminder WhatsApp message
func RemindNORTenderedWhatsAppMessage(agentName, vesselName, imoNumber, destinationContactNumber string) Payload {
	log.Println(destinationContactNumber, "handphone")
	return Payload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               destinationContactNumber, // "6591896934"
		Type:             "template",
		Template: Template{
			// this value is from Whatsapp Cloud
			Name: "columbus_nor_tender",
			Language: Language{
				Policy: "deterministic",
				Code:   "en_US",
			},
			Components: []TemplateComponent{
				{
					Type: "body",
					Parameters: []TemplateParameter{
						{Type: "text", Text: agentName},
						{Type: "text", Text: vesselName},
						{Type: "text", Text: imoNumber},
					},
				},
			},
		},
	}
}

// RemindAGDWhatsAppMessage creates a payload for a reminder WhatsApp message
func RemindAGDWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime string, destinationContactNumber string) Payload {
	log.Println(destinationContactNumber, "handphone")
	return Payload{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               destinationContactNumber, // "6591896934"
		Type:             "template",
		Template: Template{
			// this value is from Whatsapp Cloud
			Name: "columbus_agd_reminder",
			Language: Language{
				Policy: "deterministic",
				Code:   "en_US",
			},
			Components: []TemplateComponent{
				{
					Type: "body",
					Parameters: []TemplateParameter{
						{Type: "text", Text: agentName},
						{Type: "text", Text: vesselName},
						{Type: "text", Text: imoNumber},
						{Type: "text", Text: etaDate},
						{Type: "text", Text: etaTime},
					},
				},
			},
		},
	}
}
