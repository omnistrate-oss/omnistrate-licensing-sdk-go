package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidatorConfig(t *testing.T) {
	t.Parallel()

	config := NewValidatorConfig("", "", "")
	assert.Equal(t, defaultValidatorCertPath, config.CertPath)
	assert.Equal(t, defaultValidatorLicensePath, config.LicensePath)

	config = NewValidatorConfig("instance-id", "/custom/cert/path", "/custom/license/path")
	assert.Equal(t, "/custom/cert/path", config.CertPath)
	assert.Equal(t, "/custom/license/path", config.LicensePath)
	assert.Equal(t, "instance-id", config.InstanceID)
}

func TestNewValidatorConfigFromEnv(t *testing.T) {
	t.Setenv(licenseCertPathEnv, "/env/cert/path")
	t.Setenv(licenseFilePathEnv, "/env/license/path")
	t.Setenv(licenseInstanceIDEnv, "instance-id")

	config := NewValidatorConfigFromEnv()
	assert.Equal(t, "/env/cert/path", config.CertPath)
	assert.Equal(t, "/env/license/path", config.LicensePath)
	assert.Equal(t, "instance-id", config.InstanceID)
}

func TestValidatorConfigIsValid(t *testing.T) {
	t.Parallel()

	config := NewValidatorConfig("", "", "")
	assert.True(t, config.IsValid())

	config = NewValidatorConfig("", "/custom/key/path", "")
	assert.True(t, config.IsValid())

	config = NewValidatorConfig("/custom/cert/path", "", "")
	assert.True(t, config.IsValid())

	config = NewValidatorConfig("", "", "")
	config.CertPath = ""
	config.LicensePath = ""
	assert.False(t, config.IsValid())
}
