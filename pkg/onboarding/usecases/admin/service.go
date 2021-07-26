package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization/permission"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"go.opentelemetry.io/otel"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/sirupsen/logrus"
)

var tracer = otel.Tracer("github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin")

const (
	validURLSchema         = "https"
	validGrapthEndpoinTail = "/graphql"
	invalidGraphEndpoint   = "Invalid graph detected. Please add a endpoint that ends with > graphql"
)

// Usecase ...
type Usecase interface {
	RegisterMicroservice(
		ctx context.Context,
		input domain.Microservice,
	) (*domain.Microservice, error)
	CheckHealthEndpoint(healthEndpoint string) bool
	ListMicroservices(ctx context.Context) ([]*domain.Microservice, error)
	DeregisterMicroservice(ctx context.Context, id string) (bool, error)
	DeregisterAllMicroservices(ctx context.Context) (bool, error)
	FindMicroserviceByID(ctx context.Context, id string) (*domain.Microservice, error)
	PollMicroservicesStatus(ctx context.Context) ([]*domain.MicroserviceStatus, error)
}

// Service is an admin tools service e.g for registering and deregistering micro-services
type Service struct {
	baseExt extension.BaseExtension
}

// NewService initializes a valid OTP service
func NewService(ext extension.BaseExtension) *Service {
	return &Service{
		baseExt: ext,
	}
}

// CheckPreconditions ...
func (s *Service) CheckPreconditions() {
	if s.baseExt == nil {
		log.Panicf("admin service infrastructure is nil")
	}
}

// RegisterMicroservice registers a micro-service
func (s *Service) RegisterMicroservice(
	ctx context.Context,
	input domain.Microservice,
) (*domain.Microservice, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "RegisterMicroservice")
	defer span.End()

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.MicroserviceCreate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !isAuthorized {
		return nil, fmt.Errorf("user not authorized to access this resource")
	}
	//validate endpoint

	parseURL, err := url.Parse(input.URL)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("expected a valid URL, Found %v ", input.URL)
	}

	if parseURL.Scheme != validURLSchema {
		return nil, fmt.Errorf("expected a secure URL, Found %v ", input.URL)
	}

	if parseURL.Path != validGrapthEndpoinTail {
		return nil, fmt.Errorf("%v", invalidGraphEndpoint)
	}

	endPoint, err := utils.ServiceHealthEndPoint(input.URL)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("%s url is incorrect", input.Name)
	}

	if !s.CheckHealthEndpoint(endPoint) {
		return nil, fmt.Errorf("%s service is not online", input.Name)
	}

	// check if the service exists first
	filter := &firebasetools.FilterInput{
		FilterBy: []*firebasetools.FilterParam{
			{
				FieldName:           "url",
				FieldType:           enumutils.FieldTypeString,
				ComparisonOperation: enumutils.OperationEqual,
				FieldValue:          input.URL,
			},
		},
	}
	existing, _, err := firebasetools.QueryNodes(ctx, nil, filter, nil, &domain.Microservice{})
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to query microservices: %w", err)
	}
	if len(existing) > 0 {
		return nil, fmt.Errorf("a service(s) with the URL %s is already registered", input.URL)
	}

	node := &domain.Microservice{
		Name:        input.Name,
		URL:         input.URL,
		Description: input.Description,
	}
	id, _, err := firebasetools.CreateNode(ctx, node)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to register microservice: %w", err)
	}
	node.SetID(id)
	return node, nil
}

//CheckHealthEndpoint Check if service is reachable
func (s Service) CheckHealthEndpoint(healthEndpoint string) bool {
	req, err := http.NewRequest(http.MethodGet, healthEndpoint, nil)
	if err != nil {
		log.Printf("Failed to create action request with error; %v", err)
		return false
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Connection", "keep-alive")

	c := &http.Client{Timeout: time.Second * 300}
	resp, err := c.Do(req)

	if err != nil {
		log.Print(err)
		return false
	}

	defer resp.Body.Close()

	var b bool

	/* #nosec*/
	json.NewDecoder(resp.Body).Decode(&b)

	if resp.StatusCode == http.StatusOK && b {
		return true
	}
	return false
}

// ListMicroservices returns all registered micro-services
func (s *Service) ListMicroservices(ctx context.Context) ([]*domain.Microservice, error) {
	s.CheckPreconditions()

	ctx, span := tracer.Start(ctx, "ListMicroservices")
	defer span.End()

	docs, _, err := firebasetools.QueryNodes(ctx, nil, nil, nil, &domain.Microservice{})
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("unable to query microservices: %w", err)
	}

	services := []*domain.Microservice{}
	for _, doc := range docs {
		service := &domain.Microservice{}
		err = doc.DataTo(service)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("unable to unmarshal Firebase doc to microservice: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

// DeregisterMicroservice removes a micro-service
func (s *Service) DeregisterMicroservice(ctx context.Context, id string) (bool, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "DeregisterMicroservice")
	defer span.End()

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.MicroserviceDelete)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}
	return firebasetools.DeleteNode(ctx, id, &domain.Microservice{})
}

// DeregisterAllMicroservices removes all services at once. This is called internally when running in CLI mode
func (s *Service) DeregisterAllMicroservices(ctx context.Context) (bool, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "DeregisterAllMicroservices")
	defer span.End()

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.MicroserviceDelete)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}
	services, err := s.ListMicroservices(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to list all microservices: %w", err)
	}

	successCount := 0

	for _, srv := range services {
		_, err := firebasetools.DeleteNode(ctx, srv.ID, &domain.Microservice{})
		if err != nil {
			utils.RecordSpanError(span, err)
			// silent
			logrus.Errorf("failed to remove %v with error : %v", srv.Name, err)
		}
		successCount++
	}

	if successCount == len(services) {
		return true, nil
	}

	// not fatal. Recreation will happen in the next step of CLI mode
	return false, fmt.Errorf("unable to deregiseter all services")
}

// FindMicroserviceByID retrieves a micro-service by it's ID
func (s *Service) FindMicroserviceByID(ctx context.Context, id string) (*domain.Microservice, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "FindMicroserviceByID")
	defer span.End()

	node, err := firebasetools.RetrieveNode(ctx, id, &domain.Microservice{})
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get microservice with ID %s: %w", id, err)
	}
	service, ok := node.(*domain.Microservice)
	if !ok {
		return nil, fmt.Errorf("can't convert retrieved node to *Microservice")
	}
	return service, nil
}

// PollMicroservicesStatus checks if the registered microservices are serving HTTP requests/ healthy
func (s *Service) PollMicroservicesStatus(ctx context.Context) ([]*domain.MicroserviceStatus, error) {
	s.CheckPreconditions()
	ctx, span := tracer.Start(ctx, "PollMicroservicesStatus")
	defer span.End()

	services, err := s.ListMicroservices(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get registered microservices")
	}

	statuses := []*domain.MicroserviceStatus{}

	for _, service := range services {
		url, err := utils.ServiceHealthEndPoint(service.URL)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot form health url: %v", err)
		}

		reachable := s.CheckHealthEndpoint(url)

		status := &domain.MicroserviceStatus{
			Service: service,
			Active:  reachable,
		}

		statuses = append(statuses, status)

	}

	return statuses, nil
}
