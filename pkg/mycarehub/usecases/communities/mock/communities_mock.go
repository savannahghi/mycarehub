package mock

import (
	"context"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CommunityUsecaseMock contains the community usecase mock methods
type CommunityUsecaseMock struct {
	MockListMembersFn              func(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error)
	MockCreateCommunityFn          func(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
	MockListCommunityMembers       func(ctx context.Context, communityID string) ([]*domain.CommunityMember, error)
	MockRejectInviteFn             func(ctx context.Context, userID string, channelID string) (bool, error)
	MockAddMembersToCommunityFn    func(ctx context.Context, memberIDs []string, communityID string) (*stream.Response, error)
	MockBanUserFn                  func(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error)
	MockDeleteCommunityMessageFn   func(ctx context.Context, messageID string) (bool, error)
	MockValidateGetStreamRequestFn func(ctx context.Context, body []byte, signature string) bool
	MockProcessGetstreamEventsFn   func(ctx context.Context, event *dto.GetStreamEvent) error
}

// NewCommunityUsecaseMock initializes a new instance of the Community usecase happy cases
func NewCommunityUsecaseMock() *CommunityUsecaseMock {
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
		MockListCommunityMembers: func(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
			return []*domain.CommunityMember{
				{
					UserID: uuid.New().String(),
					User: domain.Member{
						ID: uuid.New().String(),
					},
				},
			}, nil
		},

		MockRejectInviteFn: func(ctx context.Context, userID string, channelID string) (bool, error) {
			return true, nil
		},
		MockAddMembersToCommunityFn: func(ctx context.Context, memberIDs []string, communityID string) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     10,
					Remaining: 10,
					Reset:     10,
				},
			}, nil
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
func (c CommunityUsecaseMock) ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
	return c.MockListCommunityMembers(ctx, communityID)
}

// RejectInvite mocks the implementation of rejecting an invitation into a community
func (c CommunityUsecaseMock) RejectInvite(ctx context.Context, userID string, channelID string, message string) (bool, error) {
	return c.MockRejectInviteFn(ctx, userID, channelID)
}

// AddMembersToCommunity mocks the implementation of adding members to a community
func (c CommunityUsecaseMock) AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (*stream.Response, error) {
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
