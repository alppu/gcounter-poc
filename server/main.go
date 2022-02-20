package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// internal types
type Message struct {
	Message string `json:'message'`
}

var messages []Message

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func addMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging
	httpMethod := r.Method
	if httpMethod == "POST" {
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		messages = append(messages, message)
	}
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging
	httpMethod := r.Method
	if httpMethod == "GET" {
		json.NewEncoder(w).Encode(messages)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/message", addMessage)
	http.HandleFunc("/messages", getMessages)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
