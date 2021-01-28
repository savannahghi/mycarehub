package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

func TestLoginClientMissingEnvs(t *testing.T) {
	username := "username"
	password := "password"

	// try login with environment variables. This should fail
	client, err := utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_CLIENT_ID
	os.Setenv("CORE_CLIENT_ID", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_CLIENT_SECRET
	os.Setenv("CORE_CLIENT_SECRET", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_USERNAME
	os.Setenv("CORE_USERNAME", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_PASSWORD
	os.Setenv("CORE_PASSWORD", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_GRANT_TYPE
	os.Setenv("CORE_GRANT_TYPE", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_API_SCHEME
	os.Setenv("CORE_API_SCHEME", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_TOKEN_URL
	os.Setenv("CORE_TOKEN_URL", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// set only CORE_TOKEN_URL
	os.Setenv("CORE_HOST", "variable")

	// try login again. This should fail
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	env1 := os.Getenv("CORE_CLIENT_ID")
	assert.Equal(t, "variable", env1)
	env2 := os.Getenv("CORE_CLIENT_SECRET")
	assert.Equal(t, "variable", env2)
	env3 := os.Getenv("CORE_USERNAME")
	assert.Equal(t, "variable", env3)
	env4 := os.Getenv("CORE_PASSWORD")
	assert.Equal(t, "variable", env4)
	env5 := os.Getenv("CORE_GRANT_TYPE")
	assert.Equal(t, "variable", env5)
	env6 := os.Getenv("CORE_API_SCHEME")
	assert.Equal(t, "variable", env6)
	env7 := os.Getenv("CORE_TOKEN_URL")
	assert.Equal(t, "variable", env7)
	env8 := os.Getenv("CORE_HOST")
	assert.Equal(t, "variable", env8)

	// try login again. This should fail because the environment variable are not correctly
	client, err = utils.LoginClient(username, password)
	assert.NotNil(t, err)
	assert.Nil(t, client)

}
