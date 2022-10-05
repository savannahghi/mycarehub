package comms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/serverutils"
)

type ICommsClient interface {
	MakeRequest(ctx context.Context, method, path string, body interface{}, authorised bool) (*http.Response, error)
}

var (
	COMMS_SENDER_ID = serverutils.MustGetEnvVar("COMMS_SENDER_ID")
)

type SILCommsLib struct {
	Client ICommsClient
}

func NewCommsSILCommsLib(client ICommsClient) *SILCommsLib {
	lib := &SILCommsLib{
		Client: client,
	}

	return lib
}

// SendBulkSMS returns a 202 Accepted synchronous response while the API attempts to send the SMS in the background.
// An asynchronous call is made to the app's sms_callback URL with a notification that shows the Bulk SMS status.
// An asynchronous call is made to the app's sms_callback individually for each of the recipients with the SMS status.
func (l SILCommsLib) SendBulkSMS(ctx context.Context, message string, recipients []string) (*BulkSMSResponse, error) {
	url := fmt.Sprintf("%s/v1/sms/bulk/", COMMS_BASE_URL)
	payload := struct {
		Sender     string   `json:"sender"`
		Message    string   `json:"message"`
		Recipients []string `json:"recipients"`
	}{
		Sender:     COMMS_SENDER_ID,
		Message:    message,
		Recipients: recipients,
	}

	response, err := l.Client.MakeRequest(ctx, http.MethodPost, url, payload, true)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusAccepted {
		err := fmt.Errorf("invalid response code %d", response.StatusCode)
		return nil, err
	}

	var resp APIResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	var bulkSMS BulkSMSResponse
	err = mapstructure.Decode(resp.Data, &bulkSMS)
	if err != nil {
		return nil, err
	}

	return &bulkSMS, nil
}

func (l SILCommsLib) GetBulkSMS(ctx context.Context, guid string) error {
	return nil
}
