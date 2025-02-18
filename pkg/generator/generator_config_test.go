package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGeneratorConfig(t *testing.T) {
	t.Parallel()

	config := NewGeneratorConfig("", "")
	assert.Equal(t, defaultGeneratorCertPath, config.CertPath)
	assert.Equal(t, defaultGeneratorKeyPath, config.KeyPath)

	config = NewGeneratorConfig("/custom/cert/path", "/custom/key/path")
	assert.Equal(t, "/custom/cert/path", config.CertPath)
	assert.Equal(t, "/custom/key/path", config.KeyPath)
}

func TestNewGeneratorConfigFromEnv(t *testing.T) {
	t.Setenv(licenseCertPathEnv, "/env/cert/path")
	t.Setenv(licenseKeyPathEnv, "/env/key/path")

	config := NewGeneratorConfigFromEnv()
	assert.Equal(t, "/env/cert/path", config.CertPath)
	assert.Equal(t, "/env/key/path", config.KeyPath)
}

func TestGeneratorConfigIsValid(t *testing.T) {
	t.Parallel()

	config := NewGeneratorConfig("", "")
	assert.True(t, config.IsValid())

	config = NewGeneratorConfig("", "/custom/key/path")
	assert.True(t, config.IsValid())

	config = NewGeneratorConfig("/custom/cert/path", "")
	assert.True(t, config.IsValid())

	config = NewGeneratorConfig("", "")
	config.CertPath = ""
	config.KeyPath = ""
	assert.False(t, config.IsValid())
}
