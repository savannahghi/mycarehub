package surveys_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	originalSurveysemailEnv := os.Getenv("SURVEYS_SYSTEM_EMAIL")
	originalSurveysPasswordEnv := os.Getenv("SURVEYS_SYSTEM_PASSWORD")
	originalSurveysBaseURLEnv := os.Getenv("SURVEYS_BASE_URL")

	os.Setenv("SURVEYS_SYSTEM_EMAIL", "test@user.com")
	os.Setenv("SURVEYS_SYSTEM_PASSWORD", "testpassword")
	os.Setenv("SURVEYS_BASE_URL", "https://example.com")

	run := m.Run()

	os.Setenv("SURVEYS_SYSTEM_EMAIL", originalSurveysemailEnv)
	os.Setenv("SURVEYS_SYSTEM_PASSWORD", originalSurveysPasswordEnv)
	os.Setenv("SURVEYS_BASE_URL", originalSurveysBaseURLEnv)

	os.Exit(run)
}
