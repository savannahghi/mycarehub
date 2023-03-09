package communities

import (
	"context"
	"strings"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
)

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error)
	ListCommunities(ctx context.Context) ([]string, error)
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	ExternalExt extension.ExternalMethodsExtension
	Matrix      matrix.Matrix
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	externalExtension extension.ExternalMethodsExtension,
	matrix matrix.Matrix,
) UseCasesCommunities {
	return &UseCasesCommunitiesImpl{
		Create:      create,
		Query:       query,
		ExternalExt: externalExtension,
		Matrix:      matrix,
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
