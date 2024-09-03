package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/db/mongodb"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/messaging/api"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	url, exists := os.LookupEnv("NATS_URL")
	if !exists {
		url = nats.DefaultURL
	} else {
		url = strings.TrimSpace(url)
	}

	if strings.TrimSpace(url) == "" {
		url = nats.DefaultURL
	}

	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		logrus.Fatal(err)
		return
	}

	err = mongodb.InitializeMongoDBClient()

	if err != nil {
		logrus.Fatal(err)
	}

	// Subcribe to each service
	api.LoginSubcriber(nc)
	api.VerifySubcriber(nc)
	api.RegisterSubcriber(nc)
	api.SendOTPSubcriber(nc)
	// Initialize MongoDB

	fmt.Println("Auth service running at port 3005")
	select {}
}
