package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestClientErrors(t *testing.T) {
	err := ClientProfileNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ClientCCCIdentifierNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ClientHasUnresolvedPinResetRequestErr()
	assert.NotNil(t, err)

}
