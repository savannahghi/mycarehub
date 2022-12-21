package dto

import "time"

// AuthCredentials is the authentication credentials for a given user
type AuthCredentials struct {
	RefreshToken string `json:"refreshToken"`
	IDToken      string `json:"idToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// GetStreamToken is the token generated from getstream
type GetStreamToken struct {
	Token string
}

// User is the output dto for the user domain object
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Active   bool   `json:"active"`

	NextAllowedLogin time.Time `json:"-"`
	FailedLoginCount int       `json:"-"`
}

// Response models the response that will be returned after a user logs in
type Response struct {
	User            *User           `json:"userProfile"`
	AuthCredentials AuthCredentials `json:"credentials"`
	GetStreamToken  string          `json:"getStreamToken"`
}

// LoginResponse models the response to be returned on successful login
type LoginResponse struct {
	Response *Response `json:"response,omitempty"`
	Message  string    `json:"message,omitempty"`
	Code     int       `json:"code,omitempty"`

	RetryTime        float64 `json:"retryTime,omitempty"`
	Attempts         int     `json:"attempts,omitempty"`
	FailedLoginCount int     `json:"failed_login_count,omitempty"`
	IsCaregiver      bool    `json:"is_caregiver"`
	IsClient         bool    `json:"is_client"`
}

// ILoginResponse represents a login response getter and setter
type ILoginResponse interface {
	SetUserProfile(user *User)
	GetUserProfile() *User

	SetResponseCode(code int, message string)

	SetAuthCredentials(credentials AuthCredentials)

	SetStreamToken(token string)

	SetRetryTime(seconds float64)

	SetFailedLoginCount(count int)

	SetIsCaregiver()
	GetIsCaregiver() bool

	SetIsClient()
	GetIsClient() bool

	ClearProfiles()
}

// NewLoginResponse initializes a new login response
func NewLoginResponse() *LoginResponse {
	return &LoginResponse{
		Response: &Response{
			User: nil,
		},
	}
}

// SetUserProfile sets the user profile
func (l *LoginResponse) SetUserProfile(user *User) {
	l.Response.User = user
}

// GetUserProfile retrieves the user profile
func (l *LoginResponse) GetUserProfile() *User {
	return l.Response.User
}

// SetResponseCode sets the response message and code
func (l *LoginResponse) SetResponseCode(code int, message string) {
	l.Code = code
	l.Message = message
}

// SetAuthCredentials sets the auth credentials
func (l *LoginResponse) SetAuthCredentials(credentials AuthCredentials) {
	l.Response.AuthCredentials = credentials
}

// SetStreamToken sets the get-stream token
func (l *LoginResponse) SetStreamToken(token string) {
	l.Response.GetStreamToken = token
}

// SetRetryTime sets the next attempt
func (l *LoginResponse) SetRetryTime(seconds float64) {
	l.RetryTime = seconds
}

// SetFailedLoginCount sets the failed login count
func (l *LoginResponse) SetFailedLoginCount(count int) {
	l.FailedLoginCount = count
}

// ClearProfiles removes the response containing the user profiles and client/staff profiles
func (l *LoginResponse) ClearProfiles() {
	l.Response = nil
}

// SetIsCaregiver indicates whether the user is a caregiver
func (l *LoginResponse) SetIsCaregiver() {
	l.IsCaregiver = true
}

// GetIsCaregiver retrieves the is caregiver value
func (l *LoginResponse) GetIsCaregiver() bool {
	return l.IsCaregiver
}

// SetIsClient indicates whether the user is a client
func (l *LoginResponse) SetIsClient() {
	l.IsClient = true
}

// GetIsClient gets the is client value
func (l *LoginResponse) GetIsClient() bool {
	return l.IsClient
}
