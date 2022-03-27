package main

import (
	"encoding/json"
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

var counter gcounter.GCounter
var serviceName string

func getCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	json.NewEncoder(w).Encode(counter)
}

func increment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	counter = gcounter.Inc(counter)

	mergeFromOthers()

	json.NewEncoder(w).Encode(counter)
}

/*
Merge incoming
*/
// TODO: Make own function of the merge stuff,
// extract it to a goroutine and call it everytime the counter is incremented
func merge(w http.ResponseWriter, r *http.Request) {
	var replica gcounter.GCounter
	err := json.NewDecoder(r.Body).Decode(&replica)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	log.Println("start merge operation")
	log.Println(replica.Counter)
	counter = gcounter.Merge(counter, replica)
	json.NewEncoder(w).Encode(counter)
}

func value(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(gcounter.Value(counter))
}

func mergeFromOthers() {
	log.Println("Starting to merge values from others")

	resp, err := http.Get("http://localhost:5000/query")

	if err != nil {
		log.Println("Failed to merge others")
		log.Fatal(err)
	}

	var response []Service
	err2 := json.NewDecoder(resp.Body).Decode(&response)
	if err2 != nil {
		log.Println("Failed to parse response")
		log.Fatal(err)
	}
	log.Printf("Got %+v \n", response)
	for _, v := range response {
		if v.Name == serviceName {
			continue // dont try to merge values from the service which is currently running
		}
		resp, err := http.Get("http://" + v.Address + "/counter")

		if err != nil {
			log.Println("Failed to call other service")
			log.Fatal(err)
		}

		var replica gcounter.GCounter
		err2 := json.NewDecoder(resp.Body).Decode(&replica)
		if err2 != nil {
			log.Println("Failed to parse response")
			log.Fatal(err)
		}
		log.Println("start merge operation")
		log.Println(replica.Counter)
		counter = gcounter.Merge(counter, replica)
	}
}

func registerRoutes() {
	http.HandleFunc("/counter/", getCounter)
	http.HandleFunc("/increment/", increment)
	http.HandleFunc("/value", value)
	http.HandleFunc("/merge", merge)
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
