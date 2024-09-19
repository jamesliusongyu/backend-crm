package whatsapp

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// WhatsAppClient defines the structure of the WhatsApp client
type WhatsAppClient struct {
	AccessToken   string
	WhatsAppApi   string
	PhoneNumberId string
}

// NewWhatsAppClient initializes a new WhatsApp client
func NewWhatsAppClient() (*WhatsAppClient, error) {
	configFile := "pkg/external/whatsapp/whatsapp_config.json"
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := struct {
		AccessToken   string `json:"accessToken"`
		WhatsAppApi   string `json:"whatsappApi"`
		PhoneNumberId string `json:"phoneNumberId"`
	}{}

	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	client := &WhatsAppClient{
		AccessToken:   config.AccessToken,
		WhatsAppApi:   config.WhatsAppApi,
		PhoneNumberId: config.PhoneNumberId,
	}

	return client, nil
}

// SendMessage sends a message using the WhatsApp client
func (client *WhatsAppClient) SendMessage(payload Payload) (*http.Response, error) {
	log.Println("hi")
	url := client.WhatsAppApi + client.PhoneNumberId + "/messages"
	log.Println(url, "url")
	log.Println(payload)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.AccessToken)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	log.Println(resp, "resp")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateShipmentWhatsAppMessage creates and sends a shipment WhatsApp message
func (client *WhatsAppClient) CreateShipmentWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime string, destinationContactNumber string) (*http.Response, error) {
	payload := CreateShipmentWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime, destinationContactNumber)
	return client.SendMessage(payload)
}

// CreateShipmentWhatsAppMessage creates and sends a shipment WhatsApp message
func (client *WhatsAppClient) RemindShipmentETAWhatsAppMessage(vesselName, imoNumber, hoursBeforeETA, destinationContactNumber string) (*http.Response, error) {
	payload := RemindShipmentETAWhatsAppMessage(vesselName, imoNumber, hoursBeforeETA, destinationContactNumber)
	return client.SendMessage(payload)
}

// CreateShipmentWhatsAppMessage creates and sends a shipment WhatsApp message
func (client *WhatsAppClient) RemindAGDWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime string, destinationContactNumber string) (*http.Response, error) {
	payload := RemindAGDWhatsAppMessage(agentName, vesselName, imoNumber, etaDate, etaTime, destinationContactNumber)
	return client.SendMessage(payload)
}

// CreateShipmentWhatsAppMessage creates and sends a shipment WhatsApp message
func (client *WhatsAppClient) RemindNORTenderedWhatsAppMessage(agentName, vesselName, imoNumber, destinationContactNumber string) (*http.Response, error) {
	payload := RemindNORTenderedWhatsAppMessage(agentName, vesselName, imoNumber, destinationContactNumber)
	return client.SendMessage(payload)
}
