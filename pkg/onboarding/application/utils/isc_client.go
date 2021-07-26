package utils

import (
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/sirupsen/logrus"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string, baseExt extension.BaseExtension) *interserviceclient.InterServiceClient {
	config, err := baseExt.LoadDepsFromYAML()
	if err != nil {
		logrus.Panicf("occurred while opening deps file %v", err)
		return nil
	}

	client, err := baseExt.SetupISCclient(*config, serviceName)
	if err != nil {
		logrus.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
		return nil
	}
	return client
}
