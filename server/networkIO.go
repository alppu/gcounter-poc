package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func registerService(registry_base_url string, serviceId string, ch chan<- Service) {
	var (
		response *http.Response
	)

	log.Println("Trying to register service")

	for response == nil {
		body := []byte(`{"id": "` + serviceId + `" }`)
		req, _ := http.NewRequest("POST", registry_base_url+"/register", bytes.NewBuffer(body))

		// host, _ := os.Hostname()
		req.Header.Set("x-forwarded-for", "localhost"+":"+os.Getenv("PORT"))
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println("Failed to register service")
			log.Println(err)
			log.Println("Retrying... in 2 seconds")
			time.Sleep(2 * time.Second)

		} else {
			var service Service
			err2 := json.NewDecoder(resp.Body).Decode(&service)

			if err2 != nil {
				log.Println("Failed to parse response")
				log.Println(err)
				log.Println("Retrying... in 2 seconds")
				time.Sleep(2 * time.Second)

			} else {
				log.Printf("Service %+v got registered \n", service)
				ch <- service
			}

		}
	}

}
