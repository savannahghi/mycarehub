package mock

import (
	"context"
	"net/http"

	"gitlab.slade360emr.com/go/apiclient"
)

// FakeServiceEDI is an `EDI` service mock
type FakeServiceEDI struct {
	LinkCoverFn func(
		ctx context.Context,
		phoneNumber string,
		uid string,
		pushToken []string,
	) (*http.Response, error)

	GetSladerDataFn func(
		ctx context.Context,
		phoneNumber string,
	) (*[]apiclient.MarketingData, error)

	LinkEDIMemberCoverFn func(
		ctx context.Context,
		phoneNumber string,
		membernumber string,
		payersladecode int,
	) (*http.Response, error)
}

// LinkCover ...
func (f *FakeServiceEDI) LinkCover(
	ctx context.Context,
	phoneNumber string,
	uid string,
	pushToken []string,
) (*http.Response, error) {
	return f.LinkCoverFn(ctx, phoneNumber, uid, pushToken)
}

// GetSladerData represents a mock for the GetSladerData method
func (f *FakeServiceEDI) GetSladerData(
	ctx context.Context,
	phoneNumber string,
) (*[]apiclient.MarketingData, error) {
	return f.GetSladerDataFn(ctx, phoneNumber)
}

// LinkEDIMemberCover represents a mock for the LinkEDIMemberCover method
func (f *FakeServiceEDI) LinkEDIMemberCover(
	ctx context.Context,
	phoneNumber string,
	membernumber string,
	payersladecode int,
) (*http.Response, error) {
	return f.LinkEDIMemberCoverFn(ctx, phoneNumber, membernumber, payersladecode)
}
