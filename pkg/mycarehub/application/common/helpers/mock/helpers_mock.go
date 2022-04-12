package mocks

// FakeReportErrorToSentryImpl mocks the sentry logic
type FakeReportErrorToSentryImpl struct {
	MockFakeReportErrorToSentryFn func(err error)
}

// NewHelper initializes a new instance of helper mock
func NewHelper() *FakeReportErrorToSentryImpl {
	return &FakeReportErrorToSentryImpl{
		MockFakeReportErrorToSentryFn: func(err error) {},
	}
}

//ReportErrorToSentry mocks the implementation reporting error
func (f *FakeReportErrorToSentryImpl) ReportErrorToSentry(err error) {
	f.MockFakeReportErrorToSentryFn(err)
}
