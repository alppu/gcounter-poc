package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"server/gcounter"
)

func getCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	json.NewEncoder(w).Encode(counter)
}

func getIncrement(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	counter = gcounter.Inc(counter)

	go mergeReplicas()

	json.NewEncoder(w).Encode(counter)
}

func getValue(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(gcounter.Value(counter))
}

func mergeReplicas() {
	log.Println("Starting to merge values to others")
	registry_base_url := requireEnvVar(ENV_REGISTRY_URL)
	resp, err := http.Get(registry_base_url + "/query")

	if err != nil {
		log.Println("Failed to get services from service registry")
		log.Println(err)
	}

	// If the query endpoint succeeds we will then cache the registered services
	if err == nil {
		var response []Service
		decodeErr := json.NewDecoder(resp.Body).Decode(&response)
		if decodeErr != nil {
			log.Println("Failed to parse service response")
			log.Println(decodeErr)
		} else {
			if len(cachedServices) == 0 {
				log.Println("Creating service cache for the first time")
			}
			for _, service := range response {
				cachedServices[service.Name] = service
			}
		}
	}

	// Iterate trough the cached service list and try to push our replica to them
	if len(cachedServices) > 0 {
		for _, service := range cachedServices {
			if service.Name == serviceName {
				continue // dont try to merge values from the service which is currently running
			}
			body, _ := json.Marshal(counter)
			_, mergeError := http.Post("http://"+service.Address+"/merge", "application/jsonn", bytes.NewBuffer(body))
			if mergeError != nil {
				log.Println("Cannot connnect to service " + service.Name)
				if serviceCircuitBreakerCount[service.Name] > 5 {
					log.Println("Removing service " + service.Name + " from cache")

				} else {
					serviceCircuitBreakerCount[service.Name] += 1
					delete(cachedServices, service.Name)
					delete(serviceCircuitBreakerCount, service.Name)
				}
				log.Println("Failed to merge replica")
				log.Println(mergeError)
			} else {
				serviceCircuitBreakerCount[service.Name] = 0
			}
		}
	} else {
		log.Println("No other services available")
	}

	log.Println("Finished merging values")
}

func getMerge(w http.ResponseWriter, r *http.Request) {
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