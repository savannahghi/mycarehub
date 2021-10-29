package extension

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
	"github.com/savannahghi/interserviceclient"
	"gopkg.in/yaml.v2"
)

// InterServiceClient struct implements an interservice communication client
type InterServiceClient struct {
	client *interserviceclient.InterServiceClient
}

// ISC represents interservice client contract
type ISC interface {
	MakeRequest(
		ctx context.Context,
		method string,
		path string,
		body interface{},
	) (*http.Response, error)
}

// NewInterServiceClient initializes a new interservice client
func NewInterServiceClient(serviceName string) *InterServiceClient {
	file, err := ioutil.ReadFile(
		filepath.Clean(interserviceclient.PathToDepsFile()),
	)
	if err != nil {
		log.Errorf("error occurred while opening deps file %v", err)
		os.Exit(1)
	}
	var config interserviceclient.DepsConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	client, err := interserviceclient.SetupISCclient(config, serviceName)
	if err != nil {
		log.Panicf(
			"unable to initialize inter service client for %v service: %s",
			err,
			serviceName,
		)
	}
	return &InterServiceClient{client: client}
}

// MakeRequest calls `interservice's` MakeRequest to make the actual API call
func (i *InterServiceClient) MakeRequest(
	ctx context.Context,
	method string,
	path string,
	body interface{},
) (*http.Response, error) {
	return i.client.MakeRequest(ctx, method, path, body)
}
