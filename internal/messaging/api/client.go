package api

import (
	"encoding/json"
	"fmt"
)

// Define the structs

func TestMain() {
	// JSON string to unmarshal
	jsonStr := `{
		"pattern": {
			"service": "auth",
			"endpoint": "login",
			"method": "POST"
		},
		"data": {
			"headers": {},
			"authorization": {},
			"params": {},
			"payload": {
				"type": [
					"info"
				],
				"status": 200,
				"data ": {
					"username": "sampleUser",
					"password": "123456"
				}
			}
		},
		"id": "255211389e207b0049f5f"
	}`

	// Unmarshal the JSON string into a Request struct
	var request Request
	unmarshalErr := json.Unmarshal([]byte(jsonStr), &request)
	if unmarshalErr != nil {
		fmt.Println("Error unmarshalling JSON:", unmarshalErr)
		return
	}

	// Print the unmarshalled struct
	fmt.Printf("%+v\n", request)
	fmt.Printf("Payload Data: %+v\n", request.Data.Payload.Data)
}
