package clients

import (
	"backend-crm/pkg/external/openai"
	"backend-crm/pkg/external/whatsapp"
	"log"
)

var WhatsAppClient *whatsapp.WhatsAppClient
var OpenAIClient *openai.OpenAIClient

func InitClients() {
	var err error
	WhatsAppClient, err = whatsapp.NewWhatsAppClient()
	if err != nil {
		log.Fatalf("Error initializing WhatsApp client: %v", err)
	}
	OpenAIClient, err = openai.NewOpenAIClient()
	if err != nil {
		log.Fatalf("Error initializing OpenAI client: %v", err)
	}

}
