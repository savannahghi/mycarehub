package domain

// Response models the response that will be returned after a user logs in
type Response struct {
	User            *User           `json:"-"`
	Client          *ClientProfile  `json:"clientProfile"`
	Staff           *StaffProfile   `json:"staffProfile"`
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
}

// ILoginResponse represents a login response getter and setter
type ILoginResponse interface {
	SetUserProfile(user *User)
	GetUserProfile() *User

	SetClientProfile(client *ClientProfile)
	GetClientProfile() *ClientProfile

	SetStaffProfile(staff *StaffProfile)
	GetStaffProfile() *StaffProfile

	SetResponseCode(code int, message string)

	SetAuthCredentials(credentials AuthCredentials)

	SetStreamToken(token string)

	SetRetryTime(seconds float64)

	SetFailedLoginCount(count int)

	SetRoles(roles []*AuthorityRole)

	SetPermissions(permissions []*AuthorityPermission)

	ClearProfiles()
}

// NewLoginResponse initializes a new login response
func NewLoginResponse() *LoginResponse {
	return &LoginResponse{
		Response: &Response{
			User:   nil,
			Client: nil,
			Staff:  nil,
		},
	}
}

// SetUserProfile sets the user profile
func (l *LoginResponse) SetUserProfile(user *User) {
	l.Response.User = user

	if l.Response.Client != nil {
		l.Response.Client.UserID = *user.ID
		l.Response.Client.User = user
	}

	if l.Response.Staff != nil {
		l.Response.Staff.UserID = *user.ID
		l.Response.Staff.User = user
	}

	l.SetFailedLoginCount(user.FailedLoginCount)
}

// GetUserProfile retrieves the user profile
func (l *LoginResponse) GetUserProfile() *User {
	return l.Response.User
}

// SetClientProfile sets the client profile
func (l *LoginResponse) SetClientProfile(client *ClientProfile) {
	l.Response.Client = client
}

// GetClientProfile retrieves the client profile
func (l *LoginResponse) GetClientProfile() *ClientProfile {
	return l.Response.Client
}

// SetStaffProfile sets the staff profile
func (l *LoginResponse) SetStaffProfile(staff *StaffProfile) {
	l.Response.Staff = staff
}

// GetStaffProfile retrieves the staff profile
func (l *LoginResponse) GetStaffProfile() *StaffProfile {
	return l.Response.Staff
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

// SetRoles sets the roles
func (l *LoginResponse) SetRoles(roles []*AuthorityRole) {
	user := l.GetUserProfile()
	user.Roles = roles
	l.SetUserProfile(user)
}

// SetPermissions sets the permissions
func (l *LoginResponse) SetPermissions(permissions []*AuthorityPermission) {
	user := l.GetUserProfile()
	user.Permissions = permissions
	l.SetUserProfile(user)
}

// ClearProfiles removes the response containing the user profiles and client/staff profiles
func (l *LoginResponse) ClearProfiles() {
	l.Response = nil
}
