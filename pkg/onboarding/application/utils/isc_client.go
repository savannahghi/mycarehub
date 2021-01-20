package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
	"gitlab.slade360emr.com/go/base"
	"gopkg.in/yaml.v2"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func NewInterServiceClient(serviceName string) *base.InterServiceClient {
	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(base.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occurred while opening deps file %v", err)
		os.Exit(1)
	}
	var config base.DepsConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	client, err := base.SetupISCclient(config, serviceName)
	if err != nil {
		log.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
	}
	return client
}
