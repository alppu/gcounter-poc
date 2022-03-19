package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/gcounter"
)

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

type Service struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func register() {
	log.Println("Registering service")
	req, _ := http.NewRequest("POST", "http://localhost:5000/register", nil)
	host, _ := os.Hostname()
	req.Header.Set("x-forwarded-for", host+":"+os.Getenv("PORT"))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Failed to register service")
		log.Fatal(err)
	}

	var response Service
	err2 := json.NewDecoder(resp.Body).Decode(&response)

	if err2 != nil {
		log.Println("Failed to parse response")
		log.Fatal(err)
	}

	log.Printf("Service %+v got registered \n", response)
	serviceName = response.Name
	counter = gcounter.Initial(serviceName)
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

func main() {
	log.Println("Starting server in port: " + os.Getenv("PORT"))
	register()
	http.HandleFunc("/counter/", getCounter)
	http.HandleFunc("/increment/", increment)
	http.HandleFunc("/value", value)
	http.HandleFunc("/merge", merge)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
