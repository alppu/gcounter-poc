package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/gcounter"
)

var counter = gcounter.Initial()

func getCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	json.NewEncoder(w).Encode(counter)
}

func increment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging

	counter = gcounter.Inc(counter)

	json.NewEncoder(w).Encode(counter)
}

/*
Merge incoming
*/
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

func main() {
	port := os.Args[1:][0]
	log.Println("Starting server in port: " + port)
	http.HandleFunc("/counter/", getCounter)
	http.HandleFunc("/increment/", increment)
	http.HandleFunc("/value", value)
	http.HandleFunc("/merge", merge)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
