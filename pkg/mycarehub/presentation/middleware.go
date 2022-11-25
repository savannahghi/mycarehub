package presentation

import (
	"context"
	"fmt"
	"net/http"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/serverutils"
)

// OrganizationID assign a default organisation to a type
var OrganizationID = serverutils.MustGetEnvVar(common.OrganizationID)

// OrganisationMiddleware retrieves a logged in user's organisation and sets it into context for the request
func OrganisationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				claims, err := firebasetools.GetLoggedInUserClaims(r.Context())
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				id, ok := claims["organisationID"].(string)
				if !ok {
					err := fmt.Errorf("expected user to have an organisation")
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				// put the organisation in the context
				ctx := context.WithValue(r.Context(), utils.OrganisationContextKey, id)
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)

			},
		)
	}
}
