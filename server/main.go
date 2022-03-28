package main

import (
	"log"
	"net/http"
	"server/gcounter"

	"github.com/google/uuid"
)

type Env string

const (
	ENV_REGISTRY_URL Env = "REGISTRY_URL"
	ENV_PORT         Env = "PORT"
)

type Service struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var cachedServices map[string]Service = make(map[string]Service)

var serviceCircuitBreakerCount map[string]int = make(map[string]int)

var counter gcounter.GCounter
var serviceName string

func registerRoutes() {
	http.HandleFunc("/counter/", getCounter)
	http.HandleFunc("/increment/", getIncrement)
	http.HandleFunc("/value", getValue)
	http.HandleFunc("/merge", getMerge)
}

func main() {
	serviceName = uuid.New().String()
	counter = gcounter.Initial(serviceName)

	port := requireEnvVar(ENV_PORT)
	log.Println("Starting server in port: " + port)

	registry_base_url := requireEnvVar(ENV_REGISTRY_URL)
	ch := make(chan Service)
	go registerService(registry_base_url, serviceName, ch)

	registerRoutes()

	log.Println("Ready to serve")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
