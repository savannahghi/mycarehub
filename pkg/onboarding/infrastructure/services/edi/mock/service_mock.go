package mock

import (
	"context"
	"net/http"
)

// FakeServiceEDI is an `EDI` service mock
type FakeServiceEDI struct {
	LinkCoverFn func(
		ctx context.Context,
		phoneNumber string,
		uid string,
	) (*http.Response, error)
}

// LinkCover ...
func (f *FakeServiceEDI) LinkCover(
	ctx context.Context,
	phoneNumber string,
	uid string,
) (*http.Response, error) {
	return f.LinkCoverFn(ctx, phoneNumber, uid)
}
