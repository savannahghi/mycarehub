package extension_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

func TestISCExtensionImpl(t *testing.T) {
	ex := extension.NewISCExtension()

	topanic := func() {
		ex.MakeRequest(http.MethodGet, "example.com", nil)
	}
	assert.Panics(t, topanic)
}
