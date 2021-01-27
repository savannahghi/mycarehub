package utils

import (
	"os"

	"github.com/labstack/gommon/log"
	"gitlab.slade360emr.com/go/base"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string) *base.InterServiceClient {
	config, err := base.LoadDepsFromYAML()
	if err != nil {
		log.Errorf("occurred while opening deps file %v", err)
		os.Exit(1)
	}

	client, err := base.SetupISCclient(*config, serviceName)
	if err != nil {
		log.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
	}
	return client
}
