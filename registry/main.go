package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

/*
Registry:
Holds in-memory database of services and their locations.
Exposes two endpoints, one to register and second one to query all
registered services
*/

type Service struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Registry []Service

type RegisterRequestBody struct {
	Id string `json:"id"`
}

var services = []Service{}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Registering service with")
	var body RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// check if registry already includes service with same name.
	service := Service{
		Name:    body.Id,
		Address: r.Header.Get("x-forwarded-for"),
	}
	log.Printf("id: %+s \n", service.Name)
	log.Printf("address: %+s \n", service.Address)

	services = append(services, service)
	log.Printf("Service %+v registered \n", service)
	json.NewEncoder(w).Encode(service)
}

func query(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(services)
}

func main() {
	log.SetOutput(os.Stdout)
	log.Println("Starting service registry in port: " + os.Getenv("PORT"))
	http.HandleFunc("/register", register)
	http.HandleFunc("/query", query)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
