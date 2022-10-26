package presentation

import (
	"context"
	"net/http"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/serverutils"
)

// OrganisationMiddleware retrieves a logged in user's organisation and sets it into context for the request
func OrganisationMiddleware(db infrastructure.Query) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				uid, err := firebasetools.GetLoggedInUserUID(r.Context())
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				user, err := db.GetUserProfileByUserID(r.Context(), uid)
				if err != nil {
					serverutils.WriteJSONResponse(w, err, http.StatusUnauthorized)
					return
				}

				// put the organisation in the context
				ctx := context.WithValue(r.Context(), utils.OrganisationContextKey, user.OrganizationID)
				// and call the next with our new context
				r = r.WithContext(ctx)

				next.ServeHTTP(w, r)

			},
		)
	}
}
