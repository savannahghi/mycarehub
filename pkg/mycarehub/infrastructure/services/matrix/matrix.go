package matrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

type Matrix interface {
	CreateCommunity(ctx context.Context, room *dto.CommunityInput) (string, error)
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
	HTTPClient *http.Client
}

// NewMatrixImpl initializes the service
func NewMatrixImpl(
	baseURL string,
) Matrix {
	client := http.Client{}
	return &ServiceImpl{
		BaseURL:    baseURL,
		HTTPClient: &client,
	}
}

// MakeRequest performs a http request and returns a response
func (m *ServiceImpl) MakeRequest(ctx context.Context, payload RequestHelperPayload) (*http.Response, error) {
	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return nil, err
	}

	p := bytes.NewBuffer(encoded)
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ") // TODO: Append Proper Matrix Access Token

	return m.HTTPClient.Do(req)
}

// CreateCommunity creates a room in Matrix homeserver
func (m *ServiceImpl) CreateCommunity(ctx context.Context, room *dto.CommunityInput) (string, error) {
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

	resp, err := m.MakeRequest(ctx, requestPayload)
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
