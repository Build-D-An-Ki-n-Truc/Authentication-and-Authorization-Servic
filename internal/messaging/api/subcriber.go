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
	Headers       Header            `json:"headers"`
	Authorization Authorization     `json:"authorization"`
	Params        map[string]string `json:"params"`
	Payload       Payload           `json:"payload"`
}

// Struct for a Request
type Request struct {
	Pattern Pattern `json:"pattern"`
	Data    Data    `json:"data"`
	ID      string  `json:"id"`
}

// Struct for a Response

type Response struct {
	Headers       Header            `json:"headers"`
	Authorization Authorization     `json:"authorization"`
	Params        map[string]string `json:"params"`
	Payload       Payload           `json:"payload"`
}

func createSubscriptionString(endpoint, method, service string) string {
	return fmt.Sprintf(`{"endpoint":"%s","method":"%s","service":"%s"}`, endpoint, method, service)
}

func RegisterSubcriber(nc *nats.Conn) {
	//subject := createSubscriptionString("register/user", "POST", "auth")

}

// Subcriber for login, Payload should have data:
//
//	Payload: Payload{
//		Type:   []string{"info"},
//		Status: http.StatusOK,
//		Data:   {
//			"username": "username",
//			"password": "password",
//		},
//	},
//
// ["roleRequired"] Ex: ["admin", "user","brand"]
func LoginSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("login/user", "POST", "auth")
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			// Get username and passwrod from user payload
			userMap := request.Data.Payload.Data.(map[string]string)

			username := string(userMap["username"])
			password := string(userMap["password"])
			fmt.Println("username: " + username)
			fmt.Println("password: " + password)
			role, check := auth.Login(username, password)

			// Login successfully
			if check {

				token, tokenErr := jwtFunc.GenerateToken(username, role)
				if tokenErr != nil {
					logrus.Panic(tokenErr)
					return
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
				message, _ := json.Marshal(response)
				m.Respond(message)
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

				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}

// Subcriber for verifying token, Payload should have data:
//
//	Payload: Payload{
//		Type:   []string{"info"},
//		Status: http.StatusOK,
//		Data:   ["roleRequired"] Ex: ["admin", "user","brand"]
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
			token := request.Data.Headers.Authorization
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
