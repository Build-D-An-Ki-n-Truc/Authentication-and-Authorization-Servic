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
// auth/login/user POST
func LoginSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("login/user", "POST", "auth")
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			userMap := request.Data.Payload.Data.(map[string]string)

			username := userMap["username"]
			password := userMap["password"]

			role, check := auth.Login(username, password)

			if check {
				// Create a token
				token, err := jwtFunc.GenerateToken(username, role)
				if err != nil {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusInternalServerError,
							Data:   err.Error(),
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				} else {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusOK,
							Data:   token,
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				}
			} else {
				response := Response{
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusUnauthorized,
						Data:   "Unauthorized",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			}
		}
	})

	if err != nil {
		log.Println(err)
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
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
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
		log.Println(err)
	}
}

// Subcriber for register, Request should have data:
//
//	Payload: Payload{
//		Data:  {
//			"username": "username",
//			"password": "password",
//			"name": "name",
//			"email": "email",
//			"role": "role",
//			"phone": "phone",
//			"isLocked": false,
//		},
//	},
//
// auth/register/user POST
func RegisterSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("register/user", "POST", "auth")
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		if unmarshalErr != nil {
			logrus.Println(unmarshalErr)
			response := Response{
				Headers:       request.Data.Headers,
				Authorization: request.Data.Authorization,
				Payload: Payload{
					Type:   []string{"info"},
					Status: http.StatusBadRequest,
					Data:   "Wrong format",
				},
			}
			message, _ := json.Marshal(response)
			m.Respond(message)
			return
		} else {
			userMap := request.Data.Payload.Data.(map[string]interface{})

			username := userMap["username"].(string)
			password := userMap["password"].(string)
			name := userMap["name"].(string)
			email := userMap["email"].(string)
			role := userMap["role"].(string)
			phone := userMap["phone"].(string)
			isLocked := userMap["isLocked"].(bool)

			// Register account (check if username already exists)
			// check == true when register success
			check, err := auth.RegisterAccount(username, password, name, email, role, phone, isLocked)

			if err != nil {
				response := Response{
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadGateway,
						Data:   err.Error(),
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
			} else {
				if check {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusAccepted,
							Data:   "Register Success",
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				} else {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusConflict,
							Data:   "Username already exists",
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				}
			}
		}
	})

	if err != nil {
		log.Println(err)
	}

}
