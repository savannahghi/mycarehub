package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// Matrix defines the methods to be used in making various matrix requests
type Matrix interface {
	CreateCommunity(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error)
	RegisterUser(ctx context.Context, username string, password string) (*dto.MatrixUserRegistrationOutput, error)
}

// RequestHelperPayload is the payload that is used to make requests to matrix client
type RequestHelperPayload struct {
	Method string
	Path   string
	Body   interface{}
}

// ServiceImpl implements the Matrix's client
type ServiceImpl struct {
	BaseURL    string
	HTTPClient http.Client
}

// NewMatrixImpl initializes the service
func NewMatrixImpl(
	baseURL string,
) Matrix {
	client := http.Client{}
	return &ServiceImpl{
		BaseURL:    baseURL,
		HTTPClient: client,
	}
}

// UserRegistration defines the structure of the input to be used when registering a Matrix user
type UserRegistration struct {
	Auth     Auth   `json:"auth"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auth is defines the type of authentication to be used when registering a new user
type Auth struct {
	Type string `json:"type"`
}

// Identifier represents the matrix identifier to be used while logging in
type Identifier struct {
	Type string `json:"type"`
	User string `json:"user"`
}

// Login is used to authenticate matrix user
func (m *ServiceImpl) Login(ctx context.Context, username string, password string) (string, error) {
	loginPayload := struct {
		Identifier *Identifier `json:"identifier"`
		Type       string      `json:"type"`
		Password   string      `json:"password"`
	}{
		Identifier: &Identifier{
			Type: "m.id.user",
			User: username,
		},
		Type:     "m.login.password",
		Password: password,
	}

	matrixRoomURL := fmt.Sprintf("%s/_matrix/client/v3/login", m.BaseURL)
	payload := RequestHelperPayload{
		Method: http.MethodPost,
		Path:   matrixRoomURL,
		Body:   loginPayload,
	}

	resp, err := m.MakeRequest(ctx, "", payload)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to login in with status code %v", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}

// MakeRequest performs a http request and returns a response
func (m *ServiceImpl) MakeRequest(ctx context.Context, token string, payload RequestHelperPayload) (*http.Response, error) {
	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return nil, err
	}

	p := bytes.NewBuffer(encoded)
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return m.HTTPClient.Do(req)
}

// RegisterUser registers a user in our Matrix homeserver
func (m *ServiceImpl) RegisterUser(ctx context.Context, username string, password string) (*dto.MatrixUserRegistrationOutput, error) {
	matrixUser := &UserRegistration{
		Auth: Auth{
			Type: "m.login.dummy",
		},
		Username: username,
		Password: password,
	}

	matrixRoomURL := fmt.Sprintf("%s/_matrix/client/v3/register", m.BaseURL)
	payload := RequestHelperPayload{
		Method: http.MethodPost,
		Path:   matrixRoomURL,
		Body:   matrixUser,
	}

	resp, err := m.MakeRequest(ctx, "", payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to register user with status code %v", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userResponse dto.MatrixUserRegistrationOutput
	err = json.Unmarshal(respBytes, &userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse, nil
}

// CreateCommunity creates a room in Matrix homeserver
func (m *ServiceImpl) CreateCommunity(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error) {
	payload := struct {
		Name       string `json:"name"`
		Topic      string `json:"topic"`
		Visibility string `json:"visibility"`
		Preset     string `json:"preset"`
	}{
		Name:       room.Name,
		Topic:      room.Topic,
		Visibility: room.Visibility.String(),
		Preset:     room.Preset.String(),
	}

	matrixRoomURL := fmt.Sprintf("%s/_matrix/client/v3/createRoom", m.BaseURL)

	requestPayload := RequestHelperPayload{
		Method: http.MethodPost,
		Path:   matrixRoomURL,
		Body:   payload,
	}

	accessToken, err := m.Login(ctx, auth.Username, auth.Password)
	if err != nil {
		return "", err
	}

	resp, err := m.MakeRequest(ctx, accessToken, requestPayload)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to create room with status code %v", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := struct {
		RoomID string `json:"room_id"`
	}{}
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", err
	}

	return data.RoomID, nil
}
