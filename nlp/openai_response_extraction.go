package nlp

import (
	database "backend-crm/internal/database/mongodb"
	"backend-crm/pkg/external/openai"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func ParseGetIntentionsFromChecklistOpenAIResponse(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	log.Print(body, "body")
	var openAIresponse openai.OpenAIResponse
	err = json.Unmarshal(body, &openAIresponse)
	if err != nil {
		return "", err
	}

	// Print the generated text
	// log.Println("Chat response:")
	// for _, choice := range openAIresponse.Choices {
	// 	log.Printf("Message: %s\n", choice.Message.Content)
	// }

	parsedIntention := openAIresponse.Choices[0].Message.Content
	log.Println("Parsed intention:", parsedIntention)

	// update the checklist with the parsed intention

	// intentions := getIntentionsFromChecklist()
	// for _, intention := range intentions {
	// 	if parsedIntention == intention
	// }
	// if parsedIntention == "intention not found" {
	// 	log.Println("Intention not found")
	// 	return nil
	// } else if parsedIntention == "marine_advisory" {
	// 	log.Println("Intention found:", parsedIntention)
	// 	return nil
	// } else if parsedIntention == "" {
	// return nil
	return parsedIntention, nil
}

func GetIntentionsFromChecklistOpenAIPrompt() string {
	// Get intentions from database.Checklist
	intentions := getIntentionsFromChecklist()

	// Create the question prompt with the formatted intentions list
	questionPrompt := strings.Join([]string{
		"Given these intentions [" + strings.Join(intentions, ", ") + "].",
		"Which intention does the input given fulfil?",
		"Select only from within the given intentions and output an array containing the intentions.",
		"If it is \"marine_advisory\" and \"medicine_supplies\", just output [\"marine_advisory\",\"medicine_supplies\"].",
		"If it is \"crew_change\", output the intention and the names of the crew as well. The order of the names must be sign on and then sign off. For example",
		"[\"crew_change\" , [\"BAN CHI TAN\", \"NG HUATSONG\"],  [\"FUNG KHANT ZAW\", \"KI AISONG\"]]",
		"Do not say anything else.",
		"If the intention does not fall within the given intentions, output [\"intention not found\"].",
	}, " ")

	// Use the questionPrompt as needed
	log.Println(questionPrompt)

	return questionPrompt
}

func getIntentionsFromChecklist() []string {
	var intentions []string
	checklistType := reflect.TypeOf(database.Checklist{})

	for i := 0; i < checklistType.NumField(); i++ {
		field := checklistType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "extras" {
			intentions = append(intentions, jsonTag)
		}
	}

	return intentions
}

func GetETAFromMasterEmailOpenAIPrompt() string {

	// Create the question prompt with the formatted intentions list
	questionPrompt := "Output the ETA. For example, \"19 Jun 2023 10:00\" or  \"08 Dec 2024 22:30\""

	// Use the questionPrompt as needed
	log.Println(questionPrompt)

	return questionPrompt
}
