package presentation_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/tj/assert"
)

var (
	testOrgID = uuid.NewString()
)

func GenerateIDTokenWithOrg(t *testing.T) string {
	ctx := context.Background()
	user, err := firebasetools.GetOrCreateFirebaseUser(ctx, firebasetools.TestUserEmail)
	if err != nil {
		t.Errorf("unable to create Firebase user for email %v, error %v", firebasetools.TestUserEmail, err)
	}

	customToken, err := firebasetools.CreateFirebaseCustomTokenWithClaims(
		ctx,
		user.UID,
		map[string]interface{}{
			"organisationID": testOrgID,
		},
	)
	if err != nil {
		t.Errorf("unable to get custom token for %#v", user)
	}

	idTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		t.Errorf("unable to exchange custom token for ID tokens, error %s", err)
	}
	if idTokens.IDToken == "" {
		t.Errorf("got blank ID token")
	}

	return idTokens.IDToken
}

func TestAuthenticationMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fc := &firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, fa)

	mw := presentation.OrganisationMiddleware()
	h := firebasetools.AuthenticationMiddleware(fa)(mw(next))

	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))

	idToken := GenerateIDTokenWithOrg(t)
	authHeader := fmt.Sprintf("Bearer %s", idToken)

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Authorization", authHeader)

	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/", reader)

	h.ServeHTTP(rw1, req1)
}
