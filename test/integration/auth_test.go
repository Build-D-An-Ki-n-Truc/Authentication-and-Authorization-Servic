package test

import (
	"encoding/json"
	"fmt"
)

// Define the structs

type Payload struct {
	Type   []string          `json:"type"`
	Status int               `json:"status"`
	Data   map[string]string `json:"data "`
}

type Data struct {
	Headers       map[string]interface{} `json:"headers"`
	Authorization map[string]interface{} `json:"authorization"`
	Params        map[string]string      `json:"params"`
	Payload       Payload                `json:"payload"`
}

type Pattern struct {
	Service  string `json:"service"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type Request struct {
	Pattern Pattern `json:"pattern"`
	Data    Data    `json:"data"`
	ID      string  `json:"id"`
}

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
