package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/db/mongodb"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/messaging/api"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics
var (
	natsMessagesReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_messages_received_total",
			Help: "Total number of messages received from NATS.",
		},
		[]string{"subject"},
	)
	mongoConnections = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "mongo_connections_total",
			Help: "Total number of connections to MongoDB.",
		},
	)
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(natsMessagesReceived)
	prometheus.MustRegister(mongoConnections)
}

func main() {
	// Expose Prometheus metrics via an HTTP endpoint
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":3050", nil); err != nil {
			logrus.Fatalf("Error starting Prometheus HTTP server: %v", err)
		}
	}()

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

	// Increment MongoDB connection counter
	mongoConnections.Inc()

	// Subcribe to each service
	api.LoginSubcriber(nc, natsMessagesReceived)
	api.VerifySubcriber(nc, natsMessagesReceived)
	api.RegisterSubcriber(nc, natsMessagesReceived)
	api.SendOTPSubcriber(nc, natsMessagesReceived)
	api.RegisterBrandBrandSubcriber(nc, natsMessagesReceived)
	// Initialize MongoDB

	fmt.Println("Auth service running at port 3005")
	select {}
}
