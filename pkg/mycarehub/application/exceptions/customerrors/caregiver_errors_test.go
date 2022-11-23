package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestCaregiverErrors(t *testing.T) {
	err := UpdateClientCaregiverErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = CreateClientCaregiverErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
