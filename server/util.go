package main

import (
	"log"
	"os"
)

func requireEnvVar(variable Env) string {
	val, present := os.LookupEnv(string(variable))
	if !present {
		log.Fatal("Required env variable " + variable + " not set!")
	}
	return val
}
