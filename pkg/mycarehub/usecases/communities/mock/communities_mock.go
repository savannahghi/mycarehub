package mock

import (
	"context"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CommunityUsecaseMock contains the community usecase mock methods
type CommunityUsecaseMock struct {
	MockListMembersFn                func(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error)
	MockCreateCommunityFn            func(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
	MockListCommunityMembersFn       func(ctx context.Context, communityID string, input *stream.QueryOption) ([]*domain.CommunityMember, error)
	MockRejectInviteFn               func(ctx context.Context, memberID string, channelID string) (bool, error)
	MockAddMembersToCommunityFn      func(ctx context.Context, memberIDs []string, communityID string) (bool, error)
	MockBanUserFn                    func(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error)
	MockDeleteCommunityMessageFn     func(ctx context.Context, messageID string) (bool, error)
	MockValidateGetStreamRequestFn   func(ctx context.Context, body []byte, signature string) bool
	MockProcessGetstreamEventsFn     func(ctx context.Context, event *dto.GetStreamEvent) error
	MockGetUserProfileByMemberIDFn   func(ctx context.Context, memberID string) (*domain.User, error)
	MockAcceptInviteFn               func(ctx context.Context, memberID string, channelID string) (bool, error)
	MockRemoveMembersFromCommunityFn func(ctx context.Context, channelID string, memberIDs []string) (bool, error)
	MockAddModeratorsWithMessageFn   func(ctx context.Context, memberIDs []string, communityID string) (bool, error)
	MockDemoteModeratorsFn           func(ctx context.Context, communityID string, memberIDs []string) (bool, error)
	MockUnBanUserFn                  func(ctx context.Context, targetID string, communityID string) (bool, error)
	MockListCommunityBannedMembersFn func(ctx context.Context, communityID string) ([]*domain.Member, error)
	MockListFlaggedMessagesFn        func(ctx context.Context, communityCID *string, memberIDs []*string) ([]*domain.MessageFlag, error)
	MockRecommendedCommunitiesFn     func(ctx context.Context, clientID string, limit int) ([]*domain.Community, error)
	MockDeleteCommunitiesFn          func(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error)
	MockInviteMembersFn              func(ctx context.Context, communityID string, memberIDs []string) (bool, error)
	MockListCommunitiesFn            func(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error)
	MockListPendingInvitesFn         func(ctx context.Context, memberID string, input *stream.QueryOption) ([]*domain.Community, error)
}

// NewCommunityUsecaseMock initializes a new instance of the Community usecase happy cases
func NewCommunityUsecaseMock() *CommunityUsecaseMock {
	UUID := gofakeit.UUID()
	now := time.Now()
	ok := true
	member := domain.Member{
		ID:       UUID,
		UserID:   UUID,
		Name:     gofakeit.Name(),
		Role:     string(enums.UserRoleTypeCommunityManagement),
		Username: gofakeit.BS(),
		Gender:   enumutils.GenderMale,
		UserType: gofakeit.BS(),
	}
	messageFlag := domain.MessageFlag{
		CreatedByAutomod: &ok,
		ModerationResult: &domain.ModerationResult{},
		Message:          &domain.GetstreamMessage{},
		User:             &member,
		CreatedAt:        &now,
	}
	community := domain.Community{
		ID:          UUID,
		CID:         UUID,
		Name:        gofakeit.BS(),
		Disabled:    false,
		Frozen:      false,
		MemberCount: 1,
		CreatedAt:   now,
		Description: gofakeit.BS(),
		InviteOnly:  false,
		CreatedBy:   &member,
		ProgramID:   UUID,
	}
	return &CommunityUsecaseMock{
		MockListMembersFn: func(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error) {
			return []*domain.Member{
				{
					ID:   uuid.New().String(),
					Role: "user",
				},
			}, nil
		},
		MockCreateCommunityFn: func(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.New().String(),
				Name:        "test",
				Description: "test",
				AgeRange: &domain.AgeRange{
					LowerBound: 1,
					UpperBound: 3,
				},
				Gender:     []enumutils.Gender{enumutils.AllGender[0]},
				ClientType: []enums.ClientType{enums.AllClientType[0]},
				InviteOnly: true,
			}, nil
		},
		MockListCommunityMembersFn: func(ctx context.Context, communityID string, input *stream.QueryOption) ([]*domain.CommunityMember, error) {
			return []*domain.CommunityMember{
				{
					UserID: uuid.New().String(),
					User: domain.Member{
						ID: uuid.New().String(),
					},
				},
			}, nil
		},

		MockRejectInviteFn: func(ctx context.Context, memberID string, channelID string) (bool, error) {
			return true, nil
		},
		MockAddMembersToCommunityFn: func(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
			return true, nil
		},
		MockBanUserFn: func(ctx context.Context, targetMemberID, bannedBy, communityID string) (bool, error) {
			return true, nil
		},
		MockDeleteCommunityMessageFn: func(ctx context.Context, messageID string) (bool, error) {
			return true, nil
		},
		MockValidateGetStreamRequestFn: func(ctx context.Context, body []byte, signature string) bool {
			return true
		},
		MockProcessGetstreamEventsFn: func(ctx context.Context, event *dto.GetStreamEvent) error {
			return nil
		},
		MockGetUserProfileByMemberIDFn: func(ctx context.Context, memberID string) (*domain.User, error) {
			ID := uuid.New().String()
			return &domain.User{
				ID:       &ID,
				Username: "test",
			}, nil
		},
		MockAcceptInviteFn: func(ctx context.Context, memberID string, channelID string) (bool, error) {
			return true, nil
		},
		MockRemoveMembersFromCommunityFn: func(ctx context.Context, channelID string, memberIDs []string) (bool, error) {
			return true, nil
		},
		MockAddModeratorsWithMessageFn: func(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
			return true, nil
		},
		MockDemoteModeratorsFn: func(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
			return true, nil
		},
		MockUnBanUserFn: func(ctx context.Context, targetID string, communityID string) (bool, error) {
			return true, nil
		},
		MockListCommunityBannedMembersFn: func(ctx context.Context, communityID string) ([]*domain.Member, error) {
			return []*domain.Member{&member}, nil
		},
		MockListFlaggedMessagesFn: func(ctx context.Context, communityCID *string, memberIDs []*string) ([]*domain.MessageFlag, error) {
			return []*domain.MessageFlag{&messageFlag}, nil
		},
		MockRecommendedCommunitiesFn: func(ctx context.Context, clientID string, limit int) ([]*domain.Community, error) {
			return []*domain.Community{&community}, nil
		},
		MockDeleteCommunitiesFn: func(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error) {
			return true, nil
		},
		MockInviteMembersFn: func(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
			return true, nil
		},
		MockListCommunitiesFn: func(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error) {
			return []*domain.Community{&community}, nil
		},
		MockListPendingInvitesFn: func(ctx context.Context, memberID string, input *stream.QueryOption) ([]*domain.Community, error) {
			return []*domain.Community{&community}, nil
		},
	}
}

// ListMembers mocks the implementation for listing getstream users
func (c CommunityUsecaseMock) ListMembers(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error) {
	return c.MockListMembersFn(ctx, input)
}

// CreateCommunity mocks the implementation of creating communities
func (c CommunityUsecaseMock) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	return c.MockCreateCommunityFn(ctx, input)
}

// ListCommunityMembers mocks the implementation of listing members
func (c CommunityUsecaseMock) ListCommunityMembers(ctx context.Context, communityID string, input *stream.QueryOption) ([]*domain.CommunityMember, error) {
	return c.MockListCommunityMembersFn(ctx, communityID, input)
}

// RejectInvite mocks the implementation of rejecting an invitation into a community
func (c CommunityUsecaseMock) RejectInvite(ctx context.Context, memberID string, channelID string) (bool, error) {
	return c.MockRejectInviteFn(ctx, memberID, channelID)
}

// AddMembersToCommunity mocks the implementation of adding members to a community
func (c CommunityUsecaseMock) AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	return c.MockAddMembersToCommunityFn(ctx, memberIDs, communityID)
}

// BanUser mocks the implementation banning a user from a specified channel
func (c CommunityUsecaseMock) BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error) {
	return c.MockBanUserFn(ctx, targetMemberID, bannedBy, communityID)
}

// DeleteCommunityMessage mocks the implementation of deleting messages
func (c CommunityUsecaseMock) DeleteCommunityMessage(ctx context.Context, messageID string) (bool, error) {
	return c.MockDeleteCommunityMessageFn(ctx, messageID)
}

// ValidateGetStreamRequest mocks the implementation of verifying a stream webhook http request
func (c CommunityUsecaseMock) ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool {
	return c.MockValidateGetStreamRequestFn(ctx, body, signature)
}

// ProcessGetstreamEvents mocks the implementation of a getstream event that has been published to our endpoint
func (c CommunityUsecaseMock) ProcessGetstreamEvents(ctx context.Context, event *dto.GetStreamEvent) error {
	return c.MockProcessGetstreamEventsFn(ctx, event)
}

// GetUserProfileByMemberID mocks the implementation of getting a user profile by member ID
func (c CommunityUsecaseMock) GetUserProfileByMemberID(ctx context.Context, memberID string) (*domain.User, error) {
	return c.MockGetUserProfileByMemberIDFn(ctx, memberID)
}

// AcceptInvite mock the implementation of the AcceptInvite method
func (c CommunityUsecaseMock) AcceptInvite(ctx context.Context, memberID string, channelID string) (bool, error) {
	return c.MockAcceptInviteFn(ctx, memberID, channelID)
}

// RemoveMembersFromCommunity mock the implementation of the RemoveMembersFromCommunity method
func (c CommunityUsecaseMock) RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string) (bool, error) {
	return c.MockRemoveMembersFromCommunityFn(ctx, channelID, memberIDs)
}

// AddModeratorsWithMessage mock the implementation of the AddModeratorsWithMessage method
func (c CommunityUsecaseMock) AddModeratorsWithMessage(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	return c.MockAddModeratorsWithMessageFn(ctx, memberIDs, communityID)
}

// DemoteModerators mock the implementation of the DemoteModerators method
func (c CommunityUsecaseMock) DemoteModerators(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	return c.MockDemoteModeratorsFn(ctx, communityID, memberIDs)
}

// UnBanUser mock the implementation of the UnBanUser method
func (c CommunityUsecaseMock) UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error) {
	return c.MockUnBanUserFn(ctx, targetID, communityID)
}

// ListCommunityBannedMembers mock the implementation of the ListCommunityBannedMembers method
func (c CommunityUsecaseMock) ListCommunityBannedMembers(ctx context.Context, communityID string) ([]*domain.Member, error) {
	return c.MockListCommunityBannedMembersFn(ctx, communityID)
}

// ListFlaggedMessages mock the implementation of the ListFlaggedMessages method
func (c CommunityUsecaseMock) ListFlaggedMessages(ctx context.Context, communityCID *string, memberIDs []*string) ([]*domain.MessageFlag, error) {
	return c.MockListFlaggedMessagesFn(ctx, communityCID, memberIDs)
}

// RecommendedCommunities mock the implementation of the RecommendedCommunities method
func (c CommunityUsecaseMock) RecommendedCommunities(ctx context.Context, clientID string, limit int) ([]*domain.Community, error) {
	return c.MockRecommendedCommunitiesFn(ctx, clientID, limit)
}

// DeleteCommunities mock the implementation of the DeleteCommunities method
func (c CommunityUsecaseMock) DeleteCommunities(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error) {
	return c.MockDeleteCommunitiesFn(ctx, communityIDs, hardDelete)
}

// InviteMembers mock the implementation of the InviteMembers method
func (c CommunityUsecaseMock) InviteMembers(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	return c.MockInviteMembersFn(ctx, communityID, memberIDs)
}

// ListCommunities mock the implementation of the ListCommunities method
func (c CommunityUsecaseMock) ListCommunities(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error) {
	return c.MockListCommunitiesFn(ctx, input)
}

// ListPendingInvites mock the implementation of the ListPendingInvites method
func (c CommunityUsecaseMock) ListPendingInvites(ctx context.Context, memberID string, input *stream.QueryOption) ([]*domain.Community, error) {
	return c.MockListPendingInvitesFn(ctx, memberID, input)
}
