package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type OpenAIClient struct {
	AccessToken string
	OrgID       string
	ProjectID   string
	OpenAIApi   string
}

// NewOpenAIClient initializes a new OpenAIClient client
func NewOpenAIClient() (*OpenAIClient, error) {
	configFile := "pkg/external/openai/openai_config.json"
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := struct {
		AccessToken string `json:"accessToken"`
		OrgID       string `json:"orgID"`
		ProjectID   string `json:"projectID"`
		OpenAIApi   string `json:"openAIApi"`
	}{}

	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	client := &OpenAIClient{
		AccessToken: config.AccessToken,
		OrgID:       config.OrgID,
		ProjectID:   config.ProjectID,
		OpenAIApi:   config.OpenAIApi,
	}

	return client, nil
}

// Calls OpenAI API to extract entity from text
func (client *OpenAIClient) ExtractEntityFromText(request OpenAIRequest) (string, error) {
	const maxRetries = 5
	const initialBackoff = time.Second
	const maxBackoff = 30 * time.Second

	var backoff time.Duration
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Marshal the request data into JSON
		requestBody, err := json.Marshal(request)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("POST", client.OpenAIApi, bytes.NewBuffer(requestBody))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+client.AccessToken)

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		if resp.StatusCode == http.StatusOK {
			// Unmarshal the response into a struct
			var response OpenAIResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				return "", err
			}

			// Print the generated text
			log.Println("Chat response:")
			parsedIntention := response.Choices[0].Message.Content
			return parsedIntention, nil
		}

		// Handle non-OK status codes
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode >= 500 {
			// Exponential backoff logic
			if attempt < maxRetries-1 {
				backoff = initialBackoff * (1 << attempt)
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				// Optionally add some random jitter
				backoff += time.Duration(rand.Int63n(int64(backoff)))

				log.Printf("Retrying after %v due to status code %d", backoff, resp.StatusCode)
				time.Sleep(backoff)
				continue
			}
		}

		// If not retriable or max retries exceeded
		log.Printf("Failed after %d retries with status code %d", attempt+1, resp.StatusCode)
		return "", nil
	}

	log.Println("Unexpected error: max retries exceeded")
	return "", nil
}
