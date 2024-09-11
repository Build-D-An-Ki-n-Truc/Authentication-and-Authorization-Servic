package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/auth"
	emailSender "github.com/Build-D-An-Ki-n-Truc/auth/internal/email"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/hashing"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/jwtFunc"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Brand    string `json:"brand,omitempty"`
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
//		data:   {
//			"username": "username",
//			"password": "password",
//		},
//	},
//
// auth/login POST
func LoginSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("login", "POST", "auth")
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

			username, ok := userMap["username"].(string)
			if !ok {
				errorCastingResponse(request, "Username should be a string", m)
				return
			}
			password, ok := userMap["password"].(string)
			if !ok {
				errorCastingResponse(request, "Password should be a string", m)
				return
			}

			// check == true if login success
			role, check := auth.Login(username, password)

			if check {
				// Create a token
				token, err := jwtFunc.GenerateToken(username, role)
				if err != nil {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusInternalServerError,
							Data:   "Error:" + err.Error(),
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				} else {
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
							Status: http.StatusOK,
							Data:   token,
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				}
				// Wrong password
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

// Subcriber for login Brand, Payload should have data:
//
//	Payload: Payload{
//		data:   {
//			"username": "username",
//			"password": "password",
//		},
//	},
//
// auth/loginBrand POST
func LoginBrandSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("loginBrand", "POST", "auth")
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

			username, ok := userMap["username"].(string)
			if !ok {
				errorCastingResponse(request, "Username should be a string", m)
				return
			}
			password, ok := userMap["password"].(string)
			if !ok {
				errorCastingResponse(request, "Password should be a string", m)
				return
			}

			// check == true if login success
			role, check, brandID := auth.LoginBrand(username, password)

			if check {
				// Create a token
				token, err := jwtFunc.GenerateToken(username, role)
				if err != nil {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusInternalServerError,
							Data:   "Error:" + err.Error(),
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				} else {
					response := Response{
						Headers: Header{
							Authorization: "Bearer " + token,
						},
						Authorization: Authorization{
							User: User{
								Username: username,
								Role:     role,
								Brand:    brandID,
							},
						},
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusOK,
							Data:   token,
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
				}
				// Wrong password
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
	subject := createSubscriptionString("verify", "POST", "auth")
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

			roleRequired, ok := request.Data.Payload.Data.([]string)
			if !ok {
				errorCastingResponse(request, "roleRequired should be an array of strings", m)
				return
			}

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
//	 Params: {
//		   crypted: bool // "true or false"
//	 }
//
//		Payload: Payload{
//
//			Data:  {
//				"username": "username",
//				"password": "password",
//				"name": "name",
//				"email": "email",
//				"role": "role",
//				"phone": "phone",
//				"isLocked": false,
//			},
//		},
//
// auth/register/user?crypted=boolean POST
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
			var crypted bool
			cryptedParam, ok := request.Data.Params["crypted"].(string)
			if !ok {
				errorCastingResponse(request, "crypted should be a string", m)
				return
			}
			if cryptedParam == "true" {
				crypted = true
			} else {
				crypted = false
			}

			username, ok := userMap["username"].(string)
			if !ok {
				errorCastingResponse(request, "Username should be a string", m)
				return
			}

			password, ok := userMap["password"].(string)
			if !ok {
				errorCastingResponse(request, "Password should be a string", m)
				return
			}
			if !crypted {
				hashedPassword, err := hashing.GenerateHash([]byte(password))
				if err != nil {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusBadRequest,
							Data:   "Password too long",
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
					return
				}
				password = string(hashedPassword)
			}
			// turn this to check if map has key then default value
			name, ok := userMap["name"].(string)
			if !ok {
				errorCastingResponse(request, "Name should be a string", m)
				return
			}

			email, ok := userMap["email"].(string)
			if !ok {
				errorCastingResponse(request, "Email should be a string", m)
				return
			}

			role, ok := userMap["role"].(string)
			if !ok {
				errorCastingResponse(request, "Role should be a string", m)
				return
			}

			phone, ok := userMap["phone"].(string)
			if !ok {
				errorCastingResponse(request, "Phone should be a string", m)
				return
			}

			isLocked, ok := userMap["isLocked"].(bool)
			if !ok {
				isLocked = false
			}
			// Register account (check if username already exists)
			// check == true when register success
			checkSuccess, err := auth.RegisterAccount(username, password, name, email, role, phone, isLocked)

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
				if checkSuccess {
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

// Send OTP email to user subcriber
// Request should have data:
//
//	Payload: Payload{
//		Data:   {
//			"email": "email",
//		},
//	},
//
// auth/sendOTP POST
// this will return a OTP if success and error if failed
func SendOTPSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("sendOTP", "POST", "auth")
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

			email, ok := userMap["email"].(string)
			if !ok {
				errorCastingResponse(request, "Email should be a string", m)
				return
			}
			//name := userMap["name"].(string)
			// Send OTP email
			otp, err := emailSender.SendEmail(email)

			if err != nil {
				response := Response{
					Headers:       request.Data.Headers,
					Authorization: request.Data.Authorization,
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadGateway,
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
						Data:   otp,
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

// Subcriber for create a Brand, Payload should have data:
//
//	 Params: {
//		   crypted: bool // "true or false"
//	 }
//
//	Payload: Payload{
//			Data:{
//				"username": username,
//				"password": password,
//				"name": name,
//				"email": email,
//				"role": brand,
//				"phone": phone,
//				"isLocked": isLocked,
//		        "brandID": brandID,
//			},
//		},
//
// auth/register/brand?crypted=bool POST	-> create a Brand User
func RegisterBrandBrandSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("register/brand", "POST", "auth")
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
			var crypted bool
			cryptedParam, ok := request.Data.Params["crypted"].(string)
			if !ok {
				errorCastingResponse(request, "crypted should be a string", m)
				return
			}
			if cryptedParam == "true" {
				crypted = true
			} else {
				crypted = false
			}

			userMap := request.Data.Payload.Data.(map[string]interface{})

			username, ok := userMap["username"].(string)
			if !ok {
				errorCastingResponse(request, "Username should be a string", m)
				return
			}

			password, ok := userMap["password"].(string)
			if !ok {
				errorCastingResponse(request, "Password should be a string", m)
				return
			}

			name, ok := userMap["name"].(string)
			if !ok {
				errorCastingResponse(request, "Name should be a string", m)
				return
			}

			email, ok := userMap["email"].(string)
			if !ok {
				errorCastingResponse(request, "Email should be a string", m)
				return
			}

			role, ok := userMap["role"].(string)
			if !ok {
				errorCastingResponse(request, "Role should be a string", m)
				return
			}

			phone, ok := userMap["phone"].(string)
			if !ok {
				errorCastingResponse(request, "Phone should be a string", m)
				return
			}

			isLocked, ok := userMap["isLocked"].(bool)
			if !ok {
				isLocked = false
			}

			brandID, ok := userMap["brandID"].(string)
			if !ok {
				errorCastingResponse(request, "BrandID should be a string", m)
				return
			}

			// Convert password to hash
			if !crypted {
				hashedPassword, err := hashing.GenerateHash([]byte(password))
				if err != nil {
					response := Response{
						Payload: Payload{
							Type:   []string{"info"},
							Status: http.StatusBadRequest,
							Data:   "Password too long",
						},
					}
					message, _ := json.Marshal(response)
					m.Respond(message)
					return
				}
				password = string(hashedPassword)
			}

			// Convert brandID to primitive.ObjectID
			brandIDObjectID, err := primitive.ObjectIDFromHex(brandID)
			if err != nil {
				response := Response{
					Payload: Payload{
						Type:   []string{"info"},
						Status: http.StatusBadRequest,
						Data:   "BrandID is not valid",
					},
				}
				message, _ := json.Marshal(response)
				m.Respond(message)
				return
			}

			// Register account (check if username already exists)
			// check == true when register success
			checkSuccess, err := auth.RegisterBrand(username, password, name, email, role, phone, isLocked, brandIDObjectID)

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
				if checkSuccess {
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

func errorCastingResponse(request Request, err string, m *nats.Msg) {
	response := Response{
		Headers:       request.Data.Headers,
		Authorization: request.Data.Authorization,
		Payload: Payload{
			Type:   []string{"info"},
			Status: http.StatusBadRequest,
			Data:   err,
		},
	}
	message, _ := json.Marshal(response)
	m.Respond(message)
}
