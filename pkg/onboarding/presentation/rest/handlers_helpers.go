package rest

import (
	"fmt"
	"net/http"

	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func decodePhoneNumberPayload(
	w http.ResponseWriter,
	r *http.Request,
	span trace.Span,
) (*dto.PhoneNumberPayload, error) {
	payload := &dto.PhoneNumberPayload{}
	serverutils.DecodeJSONToTargetStruct(w, r, payload)

	span.AddEvent("decode json payload to struct", trace.WithAttributes(
		attribute.Any("payload", payload),
	))

	if payload.PhoneNumber == nil {
		return nil, fmt.Errorf(
			"expected a phone number to be given but it was not supplied",
		)
	}

	return payload, nil
}
