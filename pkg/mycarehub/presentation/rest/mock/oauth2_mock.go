package mock

import (
	"context"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/ory/fosite"
	"golang.org/x/text/language"
)

type FositeOAuth2Mock struct {
	MockNewAuthorizeRequestFn    func(ctx context.Context, req *http.Request) (fosite.AuthorizeRequester, error)
	MockNewAuthorizeResponseFn   func(ctx context.Context, requester fosite.AuthorizeRequester, session fosite.Session) (fosite.AuthorizeResponder, error)
	MockWriteAuthorizeErrorFn    func(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, err error)
	MockWriteAuthorizeResponseFn func(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, responder fosite.AuthorizeResponder)

	MockNewAccessRequestFn    func(ctx context.Context, req *http.Request, session fosite.Session) (fosite.AccessRequester, error)
	MockNewAccessResponseFn   func(ctx context.Context, requester fosite.AccessRequester) (fosite.AccessResponder, error)
	MockWriteAccessResponseFn func(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, responder fosite.AccessResponder)
	MockWriteAccessErrorFn    func(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, err error)

	MockNewRevocationRequestFn    func(ctx context.Context, r *http.Request) error
	MockWriteRevocationResponseFn func(ctx context.Context, rw http.ResponseWriter, err error)

	MockNewIntrospectionRequestFn    func(ctx context.Context, r *http.Request, session fosite.Session) (fosite.IntrospectionResponder, error)
	MockWriteIntrospectionResponseFn func(ctx context.Context, rw http.ResponseWriter, r fosite.IntrospectionResponder)
	MockWriteIntrospectionErrorFn    func(ctx context.Context, rw http.ResponseWriter, err error)
}

func NewFositeOAuth2Mock() *FositeOAuth2Mock {
	return &FositeOAuth2Mock{
		MockNewAuthorizeRequestFn: func(ctx context.Context, req *http.Request) (fosite.AuthorizeRequester, error) {
			return &fosite.AuthorizeRequest{
				Request: fosite.Request{
					Client: &fosite.DefaultClient{
						ID: gofakeit.UUID(),
					},
				},
			}, nil
		},
		MockNewAuthorizeResponseFn: func(ctx context.Context, requester fosite.AuthorizeRequester, session fosite.Session) (fosite.AuthorizeResponder, error) {
			return &fosite.AuthorizeResponse{
				Header:     map[string][]string{},
				Parameters: map[string][]string{},
			}, nil
		},
		MockWriteAuthorizeErrorFn: func(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, err error) {},
		MockWriteAuthorizeResponseFn: func(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, responder fosite.AuthorizeResponder) {
		},
		MockNewAccessRequestFn: func(ctx context.Context, req *http.Request, session fosite.Session) (fosite.AccessRequester, error) {
			return &fosite.AccessRequest{
				GrantTypes:       []string{},
				HandledGrantType: []string{},
				Request:          fosite.Request{},
			}, nil
		},
		MockNewAccessResponseFn: func(ctx context.Context, requester fosite.AccessRequester) (fosite.AccessResponder, error) {
			return &fosite.AccessResponse{
				Extra:       map[string]interface{}{},
				AccessToken: "",
				TokenType:   "",
			}, nil
		},
		MockWriteAccessResponseFn: func(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, responder fosite.AccessResponder) {
		},
		MockWriteAccessErrorFn: func(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, err error) {},

		MockNewRevocationRequestFn: func(ctx context.Context, r *http.Request) error {
			return nil
		},
		MockWriteRevocationResponseFn: func(ctx context.Context, rw http.ResponseWriter, err error) {},
		MockNewIntrospectionRequestFn: func(ctx context.Context, r *http.Request, session fosite.Session) (fosite.IntrospectionResponder, error) {
			return &fosite.IntrospectionResponse{
				Active:          false,
				AccessRequester: nil,
				TokenUse:        "",
				AccessTokenType: "",
				Lang:            language.Tag{},
			}, nil
		},
		MockWriteIntrospectionResponseFn: func(ctx context.Context, rw http.ResponseWriter, r fosite.IntrospectionResponder) {},
		MockWriteIntrospectionErrorFn:    func(ctx context.Context, rw http.ResponseWriter, err error) {},
	}
}

// NewAuthorizeRequest mocks the implementation of NewAuthorizeRequest method
func (m *FositeOAuth2Mock) NewAuthorizeRequest(ctx context.Context, req *http.Request) (fosite.AuthorizeRequester, error) {
	return m.MockNewAuthorizeRequestFn(ctx, req)
}

// NewAuthorizeResponse mocks the implementation of NewAuthorizeResponse method
func (m *FositeOAuth2Mock) NewAuthorizeResponse(ctx context.Context, requester fosite.AuthorizeRequester, session fosite.Session) (fosite.AuthorizeResponder, error) {
	return m.MockNewAuthorizeResponseFn(ctx, requester, session)
}

// WriteAuthorizeError mocks the implementation of WriteAuthorizeError method
func (m *FositeOAuth2Mock) WriteAuthorizeError(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, err error) {
	m.MockWriteAuthorizeErrorFn(ctx, rw, requester, err)
}

// WriteAuthorizeResponse mocks the implementation of WriteAuthorizeResponse method
func (m *FositeOAuth2Mock) WriteAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, requester fosite.AuthorizeRequester, responder fosite.AuthorizeResponder) {
	m.MockWriteAuthorizeResponseFn(ctx, rw, requester, responder)
}

// NewAccessRequest mocks the implementation of NewAccessRequest method
func (m *FositeOAuth2Mock) NewAccessRequest(ctx context.Context, req *http.Request, session fosite.Session) (fosite.AccessRequester, error) {
	return m.MockNewAccessRequestFn(ctx, req, session)
}

// NewAccessResponse mocks the implementation of NewAccessResponse method
func (m *FositeOAuth2Mock) NewAccessResponse(ctx context.Context, requester fosite.AccessRequester) (fosite.AccessResponder, error) {
	return m.MockNewAccessResponseFn(ctx, requester)
}

// WriteAccessResponse mocks the implementation of WriteAccessResponse method
func (m *FositeOAuth2Mock) WriteAccessResponse(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, responder fosite.AccessResponder) {
	m.MockWriteAccessResponseFn(ctx, rw, requester, responder)
}

// WriteAccessError mocks the implementation of WriteAccessError method
func (m *FositeOAuth2Mock) WriteAccessError(ctx context.Context, rw http.ResponseWriter, requester fosite.AccessRequester, err error) {
	m.MockWriteAccessErrorFn(ctx, rw, requester, err)
}

// NewRevocationRequest mocks the implementation of NewRevocationRequest method
func (m *FositeOAuth2Mock) NewRevocationRequest(ctx context.Context, r *http.Request) error {
	return m.MockNewRevocationRequestFn(ctx, r)
}

// WriteRevocationResponse mocks the implementation of WriteRevocationResponse method
func (m *FositeOAuth2Mock) WriteRevocationResponse(ctx context.Context, rw http.ResponseWriter, err error) {
	m.MockWriteRevocationResponseFn(ctx, rw, err)
}

// NewIntrospectionRequest mocks the implementation of NewIntrospectionRequest method
func (m *FositeOAuth2Mock) NewIntrospectionRequest(ctx context.Context, r *http.Request, session fosite.Session) (fosite.IntrospectionResponder, error) {
	return m.MockNewIntrospectionRequestFn(ctx, r, session)
}

// WriteIntrospectionResponse mocks the implementation of WriteIntrospectionResponse method
func (m *FositeOAuth2Mock) WriteIntrospectionResponse(ctx context.Context, rw http.ResponseWriter, r fosite.IntrospectionResponder) {
	m.MockWriteIntrospectionResponseFn(ctx, rw, r)
}

// WriteIntrospectionError mocks the implementation of WriteIntrospectionError method
func (m *FositeOAuth2Mock) WriteIntrospectionError(ctx context.Context, rw http.ResponseWriter, err error) {
	m.MockWriteIntrospectionErrorFn(ctx, rw, err)
}
