package mock

// FakeReportErrorToSentryImpl mocks the sentry logic
type FakeReportErrorToSentryImpl struct {
	MockFakeReportErrorToSentryFn func(err error)
}

// NewFakeHelper initializes a new instance of helper mock
func NewFakeHelper() *FakeReportErrorToSentryImpl {
	return &FakeReportErrorToSentryImpl{
		MockFakeReportErrorToSentryFn: func(err error) {},
	}
}

//ReportErrorToSentry mocks the implementation reporting error
func (f *FakeReportErrorToSentryImpl) ReportErrorToSentry(err error) {
	f.MockFakeReportErrorToSentryFn(err)
}
