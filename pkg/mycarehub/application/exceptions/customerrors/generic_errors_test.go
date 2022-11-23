package customerrors

import (
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestGenericErrors(t *testing.T) {
	err := EmptyInputErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = NotActiveErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = FailedToUpdateItemErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = ItemNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = InputValidationErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = FailedToSaveItemErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

}
