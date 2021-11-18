package exceptions_test

import (
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/stretchr/testify/assert"
)

func TestModelHasCustomError(t *testing.T) {
	customerrors := exceptions.CustomError{}
	cr := customerrors.Error()
	assert.NotNil(t, cr)
}
