package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestOrganizationErrors(t *testing.T) {
	err := CreateOrganisationErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NonExistentOrganizationErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
