package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"server/gcounter"
)


func getCounter(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	json.NewEncoder(w).Encode(counter)
}

func getIncrement(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

	counter = gcounter.Inc(counter)

	go mergeReplicas()

	json.NewEncoder(w).Encode(counter)
}

func getValue(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)

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

			// Clear the cache so that we don't end up with dead connections if the registry
			// if we have cached a service that is no longer up and running
			for k := range cachedServices {
				delete(cachedServices, k)
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
				// Once the other replica cannot be synced to for 5 consecutive times
				// we will determine that it has gone dead and we will remove it from the
				// service cache. If it somehow gets back alive we can fetch it again from the registry
				// There will be no data loss since every replica stores each replicas counters in addition to own
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
	setDefaultHeaders(w)

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
