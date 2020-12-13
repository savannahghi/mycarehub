package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.slade360emr.com/go/base"
)

func main() {
	bearerToken, err := getInterserviceBearerTokenHeader()
	if err != nil {
		log.Printf("error: %s", err)
		os.Exit(-1)
	}
	log.Println(bearerToken)
}

func getInterserviceBearerTokenHeader() (string, error) {
	service := base.ISCService{} // name and domain not necessary for our use case
	isc, err := base.NewInterserviceClient(service)
	if err != nil {
		return "", fmt.Errorf("can't initialize interservice client: %w", err)
	}

	authToken, err := isc.CreateAuthToken()
	if err != nil {
		return "", fmt.Errorf("can't get auth token: %w", err)
	}
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader, nil
}
