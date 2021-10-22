package infrastructure

import (
	"context"
	"time"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
	SetUserPIN(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
	RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error)
	RegisterClient(
		ctx context.Context,
		userInput *dto.UserInput,
		clientInput *dto.ClientProfileInput,
	) (*domain.ClientUserProfile, error)
	AddIdentifier(
		ctx context.Context,
		clientID string,
		idType enums.IdentifierType,
		idValue string,
		isPrimary bool,
	) (*domain.Identifier, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id string) (bool, error)
}

// ServiceCreateImpl represents create contract implementation object
type ServiceCreateImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceCreateImpl returns new instance of ServiceCreateImpl
func NewServiceCreateImpl(on pg.OnboardingDb) Create {
	return &ServiceCreateImpl{
		onboarding: on,
	}
}

// GetOrCreateFacility is responsible for creating a representation of a facility
func (f ServiceCreateImpl) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.onboarding.GetOrCreateFacility(ctx, &facility)
}

// CollectMetrics is responsible for creating a representation of a metric
func (f ServiceCreateImpl) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
	return f.onboarding.CollectMetrics(ctx, metric)
}

// SetUserPIN saves user's PIN data
func (f ServiceCreateImpl) SetUserPIN(ctx context.Context, input *domain.UserPIN) (bool, error) {
	return f.onboarding.SetUserPIN(ctx, input)
}

// RegisterStaffUser is responsible for creating a representation of a staff user
func (f ServiceCreateImpl) RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	return f.onboarding.RegisterStaffUser(ctx, user, staff)
}

// AddIdentifier adds an identifier that is associated to a given client
func (f ServiceCreateImpl) AddIdentifier(
	ctx context.Context,
	clientID string,
	idType enums.IdentifierType,
	idValue string,
	isPrimary bool,
) (*domain.Identifier, error) {
	return f.onboarding.AddIdentifier(ctx, clientID, idType, idValue, isPrimary)
}

// RegisterClient creates a client user and saves the details in the database
func (f ServiceCreateImpl) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	return f.onboarding.RegisterClient(ctx, userInput, clientInput)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error)
}

// ServiceQueryImpl contains implementation for the Query interface
type ServiceQueryImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceQueryImpl is the initializer for Service query
func NewServiceQueryImpl(on pg.OnboardingDb) *ServiceQueryImpl {
	return &ServiceQueryImpl{
		onboarding: on,
	}
}

// RetrieveFacility  is a repository implementation method for RetrieveFacility
func (q ServiceQueryImpl) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveFacility(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode  is a repository implementation method for RetrieveFacilityByMFLCode
func (q ServiceQueryImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveByFacilityMFLCode(ctx, MFLCode, isActive)
}

//GetFacilities is responsible for returning a slice of healthcare facilities in the platform.
func (q ServiceQueryImpl) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return q.onboarding.GetFacilities(ctx)
}

// DeleteFacility is responsible for deletion of a facility from the database using the facility's id
func (f ServiceDeleteImpl) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return f.onboarding.DeleteFacility(ctx, id)
}

// GetUserProfileByUserID gets user profile by user ID
func (q ServiceQueryImpl) GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*domain.User, error) {
	return q.onboarding.GetUserProfileByUserID(ctx, userID, flavour)
}

// GetUserPINByUserID gets user PIN by user ID
func (q ServiceQueryImpl) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return q.onboarding.GetUserPINByUserID(ctx, userID)
}

// GetClientProfileByClientID fetches a client profile using the client ID
func (q ServiceQueryImpl) GetClientProfileByClientID(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
	return q.onboarding.GetClientProfileByClientID(ctx, clientID)
}

// ServiceDeleteImpl represents delete facility implementation object
type ServiceDeleteImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceDeleteImpl returns new instance of NewServiceDeleteImpl
func NewServiceDeleteImpl(on pg.OnboardingDb) Delete {
	return &ServiceDeleteImpl{
		onboarding: on,
	}
}

// Update contains all update methods
type Update interface {
	UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error
	UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error
	UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error
}

// ServiceUpdateImpl represents update user implementation object
type ServiceUpdateImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceUpdateImpl returns new instance of NewServiceUpdateImpl
func NewServiceUpdateImpl(on pg.OnboardingDb) Update {
	return &ServiceUpdateImpl{
		onboarding: on,
	}
}

// UpdateUserLastSuccessfulLogin ...
func (u *ServiceUpdateImpl) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	return u.onboarding.UpdateUserLastSuccessfulLogin(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (u *ServiceUpdateImpl) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	return u.onboarding.UpdateUserLastFailedLogin(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (u *ServiceUpdateImpl) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	return u.onboarding.UpdateUserFailedLoginCount(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (u *ServiceUpdateImpl) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	return u.onboarding.UpdateUserNextAllowedLogin(ctx, userID, nextAllowedLoginTime, flavour)
}
