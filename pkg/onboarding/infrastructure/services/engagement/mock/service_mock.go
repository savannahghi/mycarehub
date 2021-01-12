package mock

import (
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceEngagement is an `engagement` service mock .
type FakeServiceEngagement struct {
	PublishKYCNudgeFn    func(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItemFn func(uid string, payload base.Item) (*http.Response, error)
}

// PublishKYCNudge ...
func (f *FakeServiceEngagement) PublishKYCNudge(uid string, payload base.Nudge) (*http.Response, error) {
	return f.PublishKYCNudgeFn(uid, payload)
}

// PublishKYCFeedItem ...
func (f *FakeServiceEngagement) PublishKYCFeedItem(uid string, payload base.Item) (*http.Response, error) {
	return f.PublishKYCFeedItemFn(uid, payload)
}
