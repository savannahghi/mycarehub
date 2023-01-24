package pubsub_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/pubsub"
)

func TestServicePubSubImpl_ReceivePubSubPushMessages(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy Case: receive pubsub message",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakePubsub := pubsubMock.NewPubsubServiceMock()

			p := pubsub.NewUseCasePubSub(fakePubsub)

			p.ReceivePubSubPushMessages(tt.args.w, tt.args.r)
		})
	}
}
