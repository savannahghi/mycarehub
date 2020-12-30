package utils

import (
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// ValidateUID checks that the uid supplied in the indicated request is valid
func ValidateUID(w http.ResponseWriter, r *http.Request) (*domain.BusinessPartnerUID, error) {
	p := &domain.BusinessPartnerUID{}
	base.DecodeJSONToTargetStruct(w, r, p)
	if p.UID == nil {
		err := fmt.Errorf("invalid credentials, expected a uid")
		return nil, err
	}
	if p == nil {
		err := fmt.Errorf(
			"nil business partner UID struct after decoding input")
		return nil, err
	}

	return p, nil
}
