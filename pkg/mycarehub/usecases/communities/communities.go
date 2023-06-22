package communities

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
	"github.com/savannahghi/serverutils"
)

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error)
	ListCommunities(ctx context.Context) ([]string, error)
	SearchUsers(ctx context.Context, limit *int, searchTerm string) (*domain.MatrixUserSearchResult, error)
	SetPusher(ctx context.Context, flavour feedlib.Flavour) (bool, error)
	PushNotify(ctx context.Context, input *dto.MatrixNotifyInput) error
	AuthenticateUserToCommunity(ctx context.Context) (*domain.CommunityProfile, error)
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	Create       infrastructure.Create
	Query        infrastructure.Query
	ExternalExt  extension.ExternalMethodsExtension
	Matrix       matrix.Matrix
	Notification notification.UseCaseNotification
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	externalExtension extension.ExternalMethodsExtension,
	matrix matrix.Matrix,
	notification notification.UseCaseNotification,
) UseCasesCommunities {
	return &UseCasesCommunitiesImpl{
		Create:       create,
		Query:        query,
		ExternalExt:  externalExtension,
		Matrix:       matrix,
		Notification: notification,
	}
}

// CreateCommunity is used to create a new Matrix community(room)
// The aim of this API is to allow use of our backend as a means of achieving multi-tenancy.
// An example of this is allowing users to only see communities(rooms) w.r.t to organisation, program and
// facility they are currently logged in into.
func (uc *UseCasesCommunitiesImpl) CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {
	loggedInUser, err := uc.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	userProfile, err := uc.Query.GetUserProfileByUserID(ctx, loggedInUser)
	if err != nil {
		return nil, err
	}

	genders := []enumutils.Gender{}
	for _, k := range communityInput.Gender {
		genders = append(genders, enumutils.Gender(strings.ToUpper(k.String())))
	}

	clientTypes := []enums.ClientType{}
	for _, k := range communityInput.ClientType {
		clientTypes = append(clientTypes, enums.ClientType(strings.ToUpper(k.String())))
	}

	auth := &domain.MatrixAuth{
		Username: userProfile.Username,
		Password: *userProfile.ID,
	}

	roomID, err := uc.Matrix.CreateCommunity(ctx, auth, communityInput)
	if err != nil {
		return nil, err
	}

	// Set push rule
	pathValues := &domain.QueryPathValues{
		Scope:  "global",
		RuleID: "m.room.message",
		Kind:   "room",
	}

	pushRulePayload := &domain.PushRulePayload{
		Conditions: []domain.Conditions{
			{
				Kind:    "event_match",
				Key:     "type",
				Pattern: "m.room.message",
			},
		},
		Actions: []any{
			"notify",
			map[string]interface{}{
				"set_tweak": "highlight",
			},
			map[string]interface{}{
				"set_tweak": "sound",
				"value":     "default",
			},
		},
		Kind: "room",
	}

	err = uc.Matrix.SetPushRule(ctx, auth, pathValues, pushRulePayload)
	if err != nil {
		return nil, err
	}

	communityPayload := domain.Community{
		Name:        communityInput.Name,
		RoomID:      roomID,
		Description: communityInput.Topic,
		AgeRange: &domain.AgeRange{
			LowerBound: communityInput.AgeRange.LowerBound,
			UpperBound: communityInput.AgeRange.UpperBound,
		},
		Gender:         genders,
		ClientType:     clientTypes,
		OrganisationID: userProfile.CurrentOrganizationID,
		ProgramID:      userProfile.CurrentProgramID,
	}

	community, err := uc.Create.CreateCommunity(ctx, &communityPayload)
	if err != nil {
		return nil, err
	}

	return community, nil
}

// ListCommunities is used to list Matrix communities that the currently logged in user is in
func (uc *UseCasesCommunitiesImpl) ListCommunities(ctx context.Context) ([]string, error) {
	loggedInUser, err := uc.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	userProfile, err := uc.Query.GetUserProfileByUserID(ctx, loggedInUser)
	if err != nil {
		return nil, err
	}

	communities, err := uc.Query.ListCommunities(ctx, userProfile.CurrentProgramID, userProfile.CurrentOrganizationID)
	if err != nil {
		return nil, err
	}

	var communityIDs []string
	for _, community := range communities {
		communityIDs = append(communityIDs, community.RoomID)
	}

	return communityIDs, nil
}

// SearchUsers searches for users from Matrix server
func (uc *UseCasesCommunitiesImpl) SearchUsers(ctx context.Context, limit *int, searchTerm string) (*domain.MatrixUserSearchResult, error) {
	if len(searchTerm) < 3 {
		return nil, fmt.Errorf("search term must be at least 3 characters long")
	}

	loggedInUserID, err := uc.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	loggedInUserProfile, err := uc.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return nil, err
	}

	auth := &domain.MatrixAuth{
		Username: loggedInUserProfile.Username,
		Password: *loggedInUserProfile.ID,
	}

	searchResults, err := uc.Matrix.SearchUsers(ctx, *limit, searchTerm, auth)
	if err != nil {
		return nil, err
	}

	var output domain.MatrixUserSearchResult

	for _, result := range searchResults.Results {
		username := utils.TruncateMatrixUserID(result.UserID)

		userProfile, err := uc.Query.GetUserProfileByUsername(ctx, username)
		if err != nil {
			return nil, err
		}

		// if logged in user's profile is not equal to user profile of the Matrix user, skip the result
		if loggedInUserProfile.CurrentProgramID != userProfile.CurrentProgramID {
			continue
		}

		output.Results = append(output.Results, result)
	}

	return &output, nil
}

// SetPusher allows the creation, modification and deletion of pushers for a Matrix user
func (uc *UseCasesCommunitiesImpl) SetPusher(ctx context.Context, flavour feedlib.Flavour) (bool, error) {
	loggedInUserID, err := uc.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	loggedInUserProfile, err := uc.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	pusherKind := "http"
	pushGatewayURL := fmt.Sprintf("%s/_matrix/push/v1/notify", serverutils.MustGetEnvVar("SERVICE_HOST"))
	pusherPayload := &domain.PusherPayload{
		Append: true,
		PusherData: domain.PusherData{
			Format: "event_id_only",
			URL:    pushGatewayURL,
		},
		DeviceDisplayName: "myCareHub-v2", // TODO: Discuss the appropriate name for this
		Kind:              &pusherKind,
		Lang:              "en-US",
		Pushkey:           loggedInUserProfile.PushTokens[0],
	}

	switch flavour {
	case feedlib.FlavourPro:
		pusherPayload.AppDisplayName = "myCareHub Pro-v2"
		pusherPayload.AppID = serverutils.MustGetEnvVar("MYCAREHUB_PRO_APP_ID")

	case feedlib.FlavourConsumer:
		pusherPayload.AppDisplayName = "myCareHub Consumer-v2"
		pusherPayload.AppID = serverutils.MustGetEnvVar("MYCAREHUB_CONSUMER_APP_ID")

	default:
		return false, fmt.Errorf("invalid flavour")
	}

	matrixAuth := &domain.MatrixAuth{
		Username: loggedInUserProfile.Username,
		Password: *loggedInUserProfile.ID,
	}

	err = uc.Matrix.SetPusher(ctx, matrixAuth, pusherPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// PushNotify acts as the entry point to receive notifications from our Matrix chat server
func (uc *UseCasesCommunitiesImpl) PushNotify(ctx context.Context, input *dto.MatrixNotifyInput) error {
	userNotification := &domain.Notification{
		Title: "Your have a new chat message.",
		Type:  enums.NotificationTypeCommunities,
	}

	for _, device := range input.Notification.Devices {
		userProfile, err := uc.Query.GetUserProfileByPushToken(ctx, device.Pushkey)
		if err != nil {
			return err
		}

		err = uc.Notification.NotifyUser(ctx, userProfile, userNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}

	return nil
}

// AuthenticateUserToCommunity enables a user to access the community feature. It returns a user community profile
func (uc *UseCasesCommunitiesImpl) AuthenticateUserToCommunity(ctx context.Context) (*domain.CommunityProfile, error) {
	loggedInUserID, err := uc.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	loggedInUserProfile, err := uc.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	communityProfile, err := uc.Matrix.Login(ctx, loggedInUserProfile.Username, *loggedInUserProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	return communityProfile, nil
}
