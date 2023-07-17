package presentation

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

var (
	clientID     = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_ID")
	clientSecret = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_SECRET")
	// sanity check to ensure it is present
	_ = serverutils.MustGetEnvVar("MYCAREHUB_INTROSPECT_URL")
)

type IntrospectResponse struct {
	Active bool   `json:"active"`
	UserID string `json:"user_id"`
}

type IntrospectFunc func(ctx context.Context, token string) (*IntrospectResponse, error)

func Introspector(ctx context.Context, token string) (*IntrospectResponse, error) {
	tokenURL := serverutils.MustGetEnvVar("MYCAREHUB_INTROSPECT_URL")

	formData := url.Values{
		"token": []string{token},
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			helpers.ReportErrorToSentry(fmt.Errorf("Introspector() failed to close body:%w", err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to introspect token")
		return nil, err
	}

	var introspection IntrospectResponse

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bs, &introspection); err != nil {
		return nil, err
	}

	if !introspection.Active {
		err := fmt.Errorf("token is not active")
		return nil, err
	}

	return &introspection, nil
}

// AuthenticationMiddleware
func AuthenticationMiddleware(checkFunc IntrospectFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				bearerToken, err := firebasetools.ExtractBearerToken(r)
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				tokenInfo, err := checkFunc(context.Background(), bearerToken)
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusInternalServerError)
					return
				}

				if !tokenInfo.Active {
					err := fmt.Errorf("token is expired or invalid")
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				if tokenInfo.UserID == "" {
					err := fmt.Errorf("missing user ID")
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(r.Context(), firebasetools.AuthTokenContextKey, &auth.Token{UID: tokenInfo.UserID})

				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			},
		)
	}
}

// IUserProfile ...
type IUserProfile interface {
	GetUserProfile(ctx context.Context, userID string) (*domain.User, error)
}

// UserContextInjectionMiddleware retrieves a logged in user's organisation and program and sets it into context for the request
func UserContextInjectionMiddleware(us IUserProfile) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				uid, err := firebasetools.GetLoggedInUserUID(r.Context())
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				user, err := us.GetUserProfile(r.Context(), uid)
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				// put the organisation in the context
				ctx := context.WithValue(r.Context(), utils.OrganisationContextKey, user.CurrentOrganizationID)

				ctx = context.WithValue(ctx, utils.ProgramContextKey, user.CurrentProgramID)

				// and call the next with our new context
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)
			},
		)
	}
}
