package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/auth"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/jwtFunc"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// Message Pattern
/*
	{
	  "pattern": {
	    "service": "example-nestjs",
	    "endpoint": "hello",
	    "method": "GET"
	  },
	  "data": {
	    "headers": {},
	    "authorization": {},
	    "params": {
	      "name": "hai"
	    },
	    "payload": {}
	  },
	  "id": "5cb26e8dfd533783314c4"
	}
*/

type Pattern struct {
	Service  string `json:"service"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
}

type Header struct {
	Authorization string `json:"Authorization"`
}

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Authorization struct {
	User User `json:"user"`
}

// In Data should have username and password for login
type Payload struct {
	Type   []string    `json:"type"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type Data struct {
	Headers       Header                 `json:"headers"`
	Authorization Authorization          `json:"authorization"`
	Params        map[string]interface{} `json:"params"`
	Payload       Payload                `json:"payload"`
}

// Struct for a Request
type Request struct {
	Pattern Pattern `json:"pattern"`
	Data    Data    `json:"data"`
	ID      string  `json:"id"`
}

// Struct for a Response

type Response struct {
	Headers       Header                 `json:"headers"`
	Authorization Authorization          `json:"authorization"`
	Params        map[string]interface{} `json:"params"`
	Payload       Payload                `json:"payload"`
}

func createSubscriptionString(endpoint, method, service string) string {
	return fmt.Sprintf(`{"endpoint":"%s","method":"%s","service":"%s"}`, endpoint, method, service)
}

// Subcriber for login, Payload should have data:
//
//	Payload: Payload{
//		Data:   {
//			"username": "username",
//			"password": "password",
//		},
//	},
//
// ["roleRequired"] Ex: ["admin", "user","brand"]
func LoginSubcriber(nc *nats.Conn) {
	subjectUser := createSubscriptionString("login/user", "POST", "auth")
	subjectAdmin := createSubscriptionString("login/admin", "POST", "auth")
	subjectBrand := createSubscriptionString("login/brand", "POST", "auth")

	// Common function that be used between each subcriber
	// Get username and password from user payload
	getUserInfo := func(request Request) (string, string) {
		userMap := request.Data.Payload.Data.(map[string]interface{})

		username := userMap["username"].(string)
		password := userMap["password"].(string)

		return username, password
	}

	// Send Respond to client (through API Gateway)
	sendRespond := func(username string, role string, check bool) (Response, error) {
		// Login successfully
		if check {

			token, tokenErr := jwtFunc.GenerateToken(username, role)
			if tokenErr != nil {
				return Response{}, tokenErr
			}

			response := Response{
				Headers: Header{
					Authorization: "Bearer " + token,
				},
				Authorization: Authorization{
					User: User{
						Username: username,
						Role:     role,
					},
				},
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusAccepted,
					Data: map[string]string{
						"Login": "Success",
					},
				},
			}
			return response, nil
		} else { // login failed
			response := Response{
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusOK,
					Data: map[string]string{
						"Login": "Failed",
					},
				},
			}

			return response, nil
		}
	}

	// end of common function
	// Subscribe to login/user
	_, errUser := nc.Subscribe(subjectUser, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			// Get username and password from user payload
			username, password := getUserInfo(request)
			role, check := auth.Login(username, password, "casual")

			response, _ := sendRespond(username, role, check)

			message, _ := json.Marshal(response)
			m.Respond(message)
		}
	})

	if errUser != nil {
		log.Fatal(errUser)
	}

	// Subscribe to login/admin
	_, errAdmin := nc.Subscribe(subjectAdmin, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			// Get username and password from user payload
			username, password := getUserInfo(request)
			role, check := auth.Login(username, password, "admin")

			response, _ := sendRespond(username, role, check)

			message, _ := json.Marshal(response)
			m.Respond(message)
		}
	})

	if errAdmin != nil {
		log.Fatal(errAdmin)
	}

	// Subscribe to login/brand
	_, errBrand := nc.Subscribe(subjectBrand, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			// Get username and password from user payload
			username, password := getUserInfo(request)
			role, check := auth.Login(username, password, "brand")

			response, _ := sendRespond(username, role, check)

			message, _ := json.Marshal(response)
			m.Respond(message)
		}
	})

	if errBrand != nil {
		log.Fatal(errBrand)
	}
}

// Subcriber for verifying token, Request should have data:
//
//	Headers: Header{
//		Authorization: "Bearer " + token,
//	},
//
//	Authorization: Authorization{
//		User: User{
//			Username: username,
//			Role:     role,
//		},
//	},
//
//	Payload: Payload{
//		Data:   ["roleRequired"] Ex: ["admin","user","brand"]
//	},
//
// ["roleRequired"] Ex: ["admin", "user","brand"]
func VerifySubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("verify", "GET", "auth")
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			token := request.Data.Headers.Authorization[7:]
			username := request.Data.Authorization.User.Username
			role := request.Data.Authorization.User.Role

			roleRequired := request.Data.Payload.Data.([]string)

			// Verify token
			_, err := auth.VerifyRequest(token, username, role, roleRequired)

			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusUnauthorized,
						Data:   err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusOK,
						Data:   "authorized",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}

func RegisterSubcriber(nc *nats.Conn) {
	//subject := createSubscriptionString("register/user", "POST", "auth")
}
