package usecases

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	"gopkg.in/yaml.v2"
)

const engagementService = "engagement"

const (
	// engagement ISC paths
	publishNudge = "feed/%s/PRO/false/nudges/"
	publishItem  = "feed/%s/PRO/false/items/"
)

// EngagementUseCases  ...
type EngagementUseCases interface {
	PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error)
}

// EngagementUseCasesImpl represents engagement usecases
type EngagementUseCasesImpl struct {
	Engage *base.InterServiceClient
}

// NewEngagementUseCasesImpl ...
func NewEngagementUseCasesImpl(r repository.OnboardingRepository) EngagementUseCases {

	var config base.DepsConfig
	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(base.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occured while opening deps file %v", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	var client *base.InterServiceClient
	client, err = base.SetupISCclient(config, engagementService)
	if err != nil {
		log.Panicf("unable to initialize otp inter service client: %s", err)

	}

	return &EngagementUseCasesImpl{Engage: client}
}

// PublishKYCNudge ...
func (en *EngagementUseCasesImpl) PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishNudge, uid), payload)
}

// PublishKYCFeedItem ...
func (en *EngagementUseCasesImpl) PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error) {
	return en.Engage.MakeRequest("POST", fmt.Sprintf(publishItem, uid), payload)
}
