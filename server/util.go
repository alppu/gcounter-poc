package main

import (
	"log"
	"net/http"
	"os"
)

func requireEnvVar(variable Env) string {
	val, present := os.LookupEnv(string(variable))
	if !present {
		log.Fatal("Required env variable " + variable + " not set!")
	}
	return val
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // remove this after debugging
}
