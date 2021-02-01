package utils

import (
	"github.com/labstack/gommon/log"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string, baseExt extension.BaseExtension) *base.InterServiceClient {
	config, err := baseExt.LoadDepsFromYAML()
	if err != nil {
		log.Panicf("occurred while opening deps file %v", err)
		return nil
	}

	client, err := baseExt.SetupISCclient(*config, serviceName)
	if err != nil {
		log.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
		return nil
	}
	return client
}
