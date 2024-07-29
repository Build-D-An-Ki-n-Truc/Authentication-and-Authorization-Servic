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
	ContentType   string `json:"Content-Type"`
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
	Type   []string          `json:"type"`
	Status int               `json:"status"`
	Data   map[string]string `json:"data"`
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

// Login Subcriber for attaching token to user
func LoginSubcriber(nc *nats.Conn) {
	subject := createSubscriptionString("login/user", "POST", "auth")
	_, err := nc.Subscribe(subject, func(m *nats.Msg) {
		var request Request
		// parsing message to Request format
		unmarshalErr := json.Unmarshal(m.Data, &request)
		fmt.Println(m.Data)
		if unmarshalErr != nil {
			logrus.Panic(unmarshalErr)
		} else {
			// Get username and passwrod from user payload
			username := string(request.Data.Payload.Data["username"])
			password := string(request.Data.Payload.Data["password"])
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
