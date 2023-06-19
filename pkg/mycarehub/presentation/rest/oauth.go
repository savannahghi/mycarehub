package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/hashicorp/go-multierror"
	"github.com/ory/fosite"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest/html"
)

// AuthorizationSession is an object that represents a session that helps in decision making when logging in with OAuth
type AuthorizationSession struct {
	Page           string
	User           domain.User
	Program        domain.Program
	Facility       domain.Facility
	QueryParams    url.Values
	LoggedInUserID string
}

// putSession is a helper function that saves data in a session
func putSession(ctx context.Context, session AuthorizationSession, store SessionManager) error {
	bs, err := json.Marshal(session)
	if err != nil {
		return err
	}
	store.Put(ctx, "session", bs)

	return nil
}

// handleWriteAuthorizeError destroys the session and writes an authorization error
func (h *MyCareHubHandlersInterfacesImpl) handleWriteAuthorizeError(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, err error) {
	ctx := r.Context()
	sessionErr := h.sessionManager.Destroy(ctx)
	h.provider.WriteAuthorizeError(ctx, w, ar, multierror.Append(sessionErr, err))
}

// handleLoginPage helps with the rendering of the login page and on
// submission with the correct details, it also renders the next page, choose program
func (h *MyCareHubHandlersInterfacesImpl) handleLoginPage(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, authorizationSession *AuthorizationSession) {
	ctx := r.Context()

	loginInput := &dto.LoginInput{
		Username: r.FormValue("username"),
		PIN:      r.FormValue("pin"),
		Flavour:  feedlib.FlavourPro,
	}
	if loginInput.PIN == "" || loginInput.Username == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		params := html.LoginParams{
			Title: "Login",
		}

		err := html.ServeLoginPage(w, params)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		return
	}

	loginResponse, successful := h.usecase.User.Login(context.Background(), loginInput)
	if !successful {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		errorMessage := loginResponse.Message

		messageSegments := strings.Split(loginResponse.Message, ":")

		if len(messageSegments) > 1 {
			errorMessage = strings.TrimSpace(messageSegments[1])
		}

		if strings.Contains(errorMessage, "please try again after") {
			errorMessage = "Please use your mobile device to reset your pin."
		}

		params := html.LoginParams{
			Title:        "Login",
			HasError:     true,
			ErrorMessage: errorMessage,
		}
		err := html.ServeLoginPage(w, params)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		return
	}

	user, err := h.usecase.User.GetUserProfile(context.Background(), loginResponse.Response.User.ID)
	if err != nil {
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}
	}

	authorizationSession.User = *user

	authorizationSession.LoggedInUserID = loginResponse.Response.User.ID

	authorizationSession.Page = "chooseProgram"

	// save user in a session
	err = putSession(ctx, *authorizationSession, h.sessionManager)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	// render programs page
	userProgramsObject, err := h.usecase.Programs.ListUserPrograms(ctx, *authorizationSession.User.ID, feedlib.FlavourPro)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	params := html.ProgramChooserParams{
		Title:             "Programs",
		AvailablePrograms: userProgramsObject.Programs,
	}

	err = html.ServeProgramChooserPage(w, params)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

}

// handleChooseProgramPage helps with the rendering of the choose program page.
// If a program is selected successfully, it renders the next page, choose facility
func (h *MyCareHubHandlersInterfacesImpl) handleChooseProgramPage(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, authorizationSession *AuthorizationSession) {
	ctx := r.Context()

	// save the selected program in the session
	programID := r.FormValue("program")

	if programID == "" {
		// render programs page
		userProgramsObject, err := h.usecase.Programs.ListUserPrograms(ctx, *authorizationSession.User.ID, feedlib.FlavourPro)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		params := html.ProgramChooserParams{
			Title:             "Programs",
			HasError:          true,
			ErrorMessage:      "Please select a program",
			AvailablePrograms: userProgramsObject.Programs,
		}

		err = html.ServeProgramChooserPage(w, params)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		return
	}

	selectedProgram, err := h.usecase.Programs.GetProgramByID(ctx, programID)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	authorizationSession.Program = *selectedProgram

	authorizationSession.Page = "chooseFacility"

	// save program in a session
	err = putSession(ctx, *authorizationSession, h.sessionManager)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	ctx = context.WithValue(ctx, firebasetools.AuthTokenContextKey, &auth.Token{UID: authorizationSession.LoggedInUserID})

	staffProfile, err := h.usecase.User.GetStaffProfile(context.Background(), authorizationSession.LoggedInUserID, authorizationSession.Program.ID)
	if err != nil {
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}
	}

	userFacilitiesObject, err := h.usecase.User.GetStaffFacilities(ctx, *staffProfile.ID, dto.PaginationsInput{Limit: 10, CurrentPage: 1})
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	// render facility page
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	params := html.FacilityChooserParams{
		Title:               "Facilities",
		AvailableFacilities: userFacilitiesObject.Facilities,
	}

	err = html.ServeFacilityChooserPage(w, params)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

}

// handleChooseFacilityPage helps with the rendering of the choose facility page.
// If a program is selected successfully, it renders the next page, homepage
func (h *MyCareHubHandlersInterfacesImpl) handleChooseFacilityPage(w http.ResponseWriter, r *http.Request, ar fosite.AuthorizeRequester, authorizationSession *AuthorizationSession) {
	ctx := r.Context()
	facilityID := r.FormValue("facility")

	if facilityID == "" {
		ctx = context.WithValue(ctx, firebasetools.AuthTokenContextKey, &auth.Token{UID: authorizationSession.LoggedInUserID})
		userFacilitiesObject, err := h.usecase.Facility.ListProgramFacilities(ctx, nil, nil, nil, &dto.PaginationsInput{Limit: 10, CurrentPage: 1})
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		// render facility page
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		params := html.FacilityChooserParams{
			Title: "Facilities",
			// add error message
			AvailableFacilities: userFacilitiesObject.Facilities,
		}

		err = html.ServeFacilityChooserPage(w, params)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		return
	}

	selectedFacility, err := h.usecase.Facility.RetrieveFacility(ctx, &facilityID, true)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}

	authorizationSession.Facility = *selectedFacility

	// save facility in a session
	err = putSession(ctx, *authorizationSession, h.sessionManager)
	if err != nil {
		h.handleWriteAuthorizeError(w, r, ar, err)
		return
	}
}

// AuthorizeHandler contains the OAuth logic for authorization. it follows the following steps:
// 0. get or create http session
// 1. Get the client. Determine if its PRO or CONSUMER Login. (Currently supporting PRO)
// 2. Present the user with the username/PIN login page
// 3. Present user with program chooser page.
// 4. Use selected program to retrieve staff profile and add to session
// 5. Present user with facility chooser page.
// 6. Add selected facility to session
// 7. If successful, Initialize an oauth session with the user ID
// 8. Return success
func (h *MyCareHubHandlersInterfacesImpl) AuthorizeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ar, err := h.provider.NewAuthorizeRequest(ctx, r)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		// 1. Get the client. Determine if its PRO or CONSUMER Login. (Currently supporting PRO)
		client := ar.GetClient()

		var authorizationSession = AuthorizationSession{}

		exists := h.sessionManager.Exists(r.Context(), "session")
		if exists {
			authorizationSessionBs := h.sessionManager.GetBytes(ctx, "session")
			_ = json.Unmarshal(authorizationSessionBs, &authorizationSession)
		} else {

			queryParams := r.URL.Query()
			authorizationSession = AuthorizationSession{
				Page:        "login",
				QueryParams: queryParams,
			}

			err = putSession(ctx, authorizationSession, h.sessionManager)
			if err != nil {
				h.handleWriteAuthorizeError(w, r, ar, err)
				return
			}
		}

		switch authorizationSession.Page {
		case "login":
			h.handleLoginPage(w, r, ar, &authorizationSession)
			return

		case "chooseProgram":
			h.handleChooseProgramPage(w, r, ar, &authorizationSession)
			return

		case "chooseFacility":
			h.handleChooseFacilityPage(w, r, ar, &authorizationSession)

		}

		r.URL.RawQuery = authorizationSession.QueryParams.Encode()

		ar, err = h.provider.NewAuthorizeRequest(ctx, r)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		user := authorizationSession.User
		program := authorizationSession.Program
		facility := authorizationSession.Facility
		ctx = context.WithValue(ctx, firebasetools.AuthTokenContextKey, &auth.Token{UID: authorizationSession.LoggedInUserID})

		_, err = h.usecase.Programs.SetCurrentProgram(ctx, program.ID)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		staffProfile, err := h.usecase.User.GetStaffProfile(ctx, *user.ID, program.ID)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		_, err = h.usecase.User.SetStaffDefaultFacility(ctx, *staffProfile.ID, *facility.ID)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		extraDetails := map[string]interface{}{
			"user_id":               user.ID,
			"organisation_id":       program.Organisation.ID,
			"program_id":            program.ID,
			"gender":                user.Gender,
			"is_superuser":          user.IsSuperuser,
			"staff_id":              *staffProfile.ID,
			"facility_id":           *facility.ID,
			"is_organisation_admin": false,
			"is_program_admin":      false,
		}

		if user.Email != nil {
			extraDetails["email"] = *user.Email
		}

		session := domain.NewSession(
			ctx,
			client.GetID(),
			*user.ID,
			user.Username,
			user.Name,
			extraDetails,
		)

		response, err := h.provider.NewAuthorizeResponse(ctx, ar, session)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		// cleanup http sessions
		err = h.sessionManager.Destroy(ctx)
		if err != nil {
			h.handleWriteAuthorizeError(w, r, ar, err)
			return
		}

		h.provider.WriteAuthorizeResponse(ctx, w, ar, response)

	}
}

// TokenHandler generates a token for the given client
func (h *MyCareHubHandlersInterfacesImpl) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ar, err := h.provider.NewAccessRequest(ctx, r, new(domain.Session))
		if err != nil {
			h.provider.WriteAccessError(ctx, w, ar, err)
			return
		}

		// if client credentials grant, add the client to the session
		if ar.GetGrantTypes().ExactOne("client_credentials") {
			client := ar.GetClient()
			session := ar.GetSession().(*domain.Session)

			session.ClientID = client.GetID()
			ar.SetSession(session)
		}

		response, err := h.provider.NewAccessResponse(ctx, ar)
		if err != nil {
			h.provider.WriteAccessError(ctx, w, ar, err)
			return
		}

		h.provider.WriteAccessResponse(ctx, w, ar, response)

	}
}

// RevokeHandler expires an existing token which will then no longer be valid
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

// IntrospectionHandler returns the content of a given token
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
