package rest

import (
	"net/http"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest/html"
)

func (h *MyCareHubHandlersInterfacesImpl) AuthorizeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ar, err := h.provider.NewAuthorizeRequest(ctx, r)
		if err != nil {
			h.provider.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		/*
			---- Happy Path ----
			(TODO: This approach requires an actual stored session to complete i.e remembering state)
			1. Get the client. Determine if its PRO or CONSUMER Login. (Currently supporting PRO)
			2. Present the user with the username/PIN login form
			3. If successful, Initialize a session with the user ID
			4. Present user with program chooser page.
			5. Use selected program to retrieve staff profile and add to session
			6. Present user with facility chooser page.
			7. Add selected facility to session
			8. Return success
		*/

		// 1. Get the client. Determine if its PRO or CONSUMER Login. (Currently supporting PRO)
		client := ar.GetClient()
		flavour := feedlib.FlavourPro

		// 2. Check if username exists. If not present the user with the username/PIN login form
		loginInput := &dto.LoginInput{
			Username: r.FormValue("username"),
			PIN:      r.FormValue("pin"),
			Flavour:  flavour,
		}

		if loginInput.PIN == "" || loginInput.Username == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			params := html.LoginParams{
				Title: "Login",
			}

			_ = html.ServeLoginPage(w, params)

			return
		}

		loginResponse, successful := h.usecase.User.Login(ctx, loginInput)
		if !successful {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			params := html.LoginParams{
				Title:        "Login",
				HasError:     true,
				ErrorMessage: loginResponse.Message,
			}
			_ = html.ServeLoginPage(w, params)

			return
		}

		// 3. If successful, Initialize a session with the user profile
		user, err := h.usecase.User.GetUserProfile(ctx, loginResponse.Response.User.ID)
		if err != nil {
			h.provider.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		session := domain.NewSession(
			ctx,
			client.GetID(),
			*user.ID,
			user.Username,
			user.Name,
			map[string]interface{}{
				"user_id":         user.ID,
				"email":           *user.Email,
				"organisation_id": user.CurrentOrganizationID,
				"program_id":      user.CurrentProgramID,
				"gender":          user.Gender,
				"is_superuser":    user.IsSuperuser,
			},
		)

		response, err := h.provider.NewAuthorizeResponse(ctx, ar, session)
		if err != nil {
			h.provider.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		h.provider.WriteAuthorizeResponse(ctx, w, ar, response)

	}
}

func (h *MyCareHubHandlersInterfacesImpl) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ar, err := h.provider.NewAccessRequest(ctx, r, new(domain.Session))
		if err != nil {
			h.provider.WriteAccessError(ctx, w, ar, err)
			return
		}

		response, err := h.provider.NewAccessResponse(ctx, ar)
		if err != nil {
			h.provider.WriteAccessError(ctx, w, ar, err)
			return
		}

		h.provider.WriteAccessResponse(ctx, w, ar, response)
	}
}

func (h *MyCareHubHandlersInterfacesImpl) RevokeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := h.provider.NewRevocationRequest(ctx, r)
		if err != nil {
			h.provider.WriteRevocationResponse(ctx, w, err)
			return
		}

		h.provider.WriteRevocationResponse(ctx, w, nil)
	}
}

func (h *MyCareHubHandlersInterfacesImpl) IntrospectionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ir, err := h.provider.NewIntrospectionRequest(ctx, r, new(domain.Session))
		if err != nil {
			h.provider.WriteIntrospectionError(ctx, w, err)
			return
		}

		h.provider.WriteIntrospectionResponse(ctx, w, ir)
	}
}
