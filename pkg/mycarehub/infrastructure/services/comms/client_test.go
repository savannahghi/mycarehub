package comms

import (
	"net/http"
	"testing"
)

func TestSILComms_Login(t *testing.T) {
	type fields struct {
		client http.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "happy case: successful login",
			fields: fields{
				client: http.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SILCommsClient{
				client: tt.fields.client,
			}
			s.login()
		})
	}
}
