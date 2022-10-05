package comms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
)

var (
	COMMS_BASE_URL        = serverutils.MustGetEnvVar("COMMS_BASE_URL")
	COMMS_EMAIL           = serverutils.MustGetEnvVar("COMMS_EMAIL")
	COMMS_PASSWORD        = serverutils.MustGetEnvVar("COMMS_PASSWORD")
	ACCESS_TOKEN_TIMEOUT  = 30 * time.Minute
	REFRESH_TOKEN_TIMEOUT = 24 * time.Hour
)

// SILCommsClient is the implementation
type SILCommsClient struct {
	client http.Client

	refreshToken       string
	refreshTokenTicker *time.Ticker

	accessToken       string
	accessTokenTicker *time.Ticker
}

func NewSILCommsClient() *SILCommsClient {
	s := &SILCommsClient{
		client:       http.Client{},
		accessToken:  "",
		refreshToken: "",
	}
	s.login()
	go s.background()

	return s
}

func (s *SILCommsClient) background() {
	for {
		select {
		case t := <-s.refreshTokenTicker.C:
			logrus.Println("SIL Comms Refresh Token updated at: ", t)
			s.login()

		case t := <-s.accessTokenTicker.C:
			logrus.Println("SIL Comms Access Token updated at: ", t)
			s.refreshAccessToken()

		}
	}
}

func (s *SILCommsClient) setAccessToken(token string) {
	s.accessToken = token
	if s.accessTokenTicker != nil {
		s.accessTokenTicker.Reset(ACCESS_TOKEN_TIMEOUT)
	} else {
		s.accessTokenTicker = time.NewTicker(ACCESS_TOKEN_TIMEOUT)
	}
}

func (s *SILCommsClient) setRefreshToken(token string) {
	s.refreshToken = token
	if s.refreshTokenTicker != nil {
		s.refreshTokenTicker.Reset(REFRESH_TOKEN_TIMEOUT)
	} else {
		s.refreshTokenTicker = time.NewTicker(REFRESH_TOKEN_TIMEOUT)
	}
}

func (s *SILCommsClient) login() {
	url := fmt.Sprintf("%s/auth/token/", COMMS_BASE_URL)
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    COMMS_EMAIL,
		Password: COMMS_PASSWORD,
	}

	response, err := s.MakeRequest(context.Background(), http.MethodPost, url, payload, false)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		panic("kimeumana")
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}

	var tokens TokenResponse
	err = mapstructure.Decode(resp.Data, &tokens)
	if err != nil {
		panic(err)
	}

	s.setRefreshToken(tokens.Refresh)
	s.setAccessToken(tokens.Access)

}

func (s *SILCommsClient) refreshAccessToken() {
	url := fmt.Sprintf("%s/auth/token/refresh/", COMMS_BASE_URL)
	payload := struct {
		Refresh string `json:"refresh"`
	}{
		Refresh: s.refreshToken,
	}

	response, err := s.MakeRequest(context.Background(), http.MethodPost, url, payload, false)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		panic("kimeumana")
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}

	var tokens TokenResponse
	err = mapstructure.Decode(resp.Data, &tokens)
	if err != nil {
		panic(err)
	}

	s.setAccessToken(tokens.Access)

}

func (s *SILCommsClient) MakeRequest(ctx context.Context, method, path string, body interface{}, authorised bool) (*http.Response, error) {
	var request *http.Request

	switch method {
	case http.MethodGet:
		req, err := http.NewRequestWithContext(ctx, method, path, nil)
		if err != nil {
			return nil, err
		}
		request = req
	case http.MethodPost:
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		payload := bytes.NewBuffer(encoded)

		req, err := http.NewRequestWithContext(ctx, method, path, payload)
		if err != nil {
			return nil, err
		}

		request = req
	default:
		return nil, fmt.Errorf("s.MakeRequest() unsupported http method: %s", method)

	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	if authorised {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	}

	return s.client.Do(request)
}
