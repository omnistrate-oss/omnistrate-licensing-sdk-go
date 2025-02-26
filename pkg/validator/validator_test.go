package validator

import (
	_ "embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/common"
	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed certificate-test-tls.crt
var certPEM []byte

//go:embed certificate-test-tls.key
var keyPEM []byte

func TestNewValidator(t *testing.T) {
	t.Parallel()

	validator, err := NewValidatorFromBytes(certPEM)
	require.NoError(t, err)
	require.NotNil(t, validator)
}

func TestManager_ValidateLicense(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	manager, err := generator.NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	envelope, err := manager.GenerateLicense("orgId", "SKU", "instance-1", "subs-1", "product a", now.Add(48*time.Hour))
	assert.NoError(t, err)
	assert.NotNil(t, envelope)

	validator, err := NewValidatorFromBytes(certPEM)
	require.NoError(t, err)

	err = validator.ValidateLicense(envelope, "orgId", "SKU", "instance-1", now)
	assert.NoError(t, err)

	err = validator.ValidateLicense(envelope, "INVALID", "SKU", "instance-1", now)
	assert.Error(t, err)

	err = validator.ValidateLicense(envelope, "orgId", "INVALID", "instance-1", now)
	assert.Error(t, err)

	err = validator.ValidateLicense(envelope, "orgId", "SKU", "INVALID", now)
	assert.Error(t, err)
}

func TestManager_ValidateLicenseBase64(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	manager, err := generator.NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	license, err := manager.GenerateLicense("orgId", "SKU", "instance-1", "subs-1", "product a", now.Add(48*time.Hour))
	require.NoError(t, err)

	licenseBase64 := license.String()
	require.NoError(t, err)
	assert.NotEmpty(t, licenseBase64)

	decoded, err := common.DecodeLicenseEnvelopeFromString(licenseBase64)
	require.NoError(t, err)

	assert.Equal(t, license.License, decoded.License)
	assert.Equal(t, license.Signature, decoded.Signature)

	validator, err := NewValidatorFromBytes(certPEM)
	require.NoError(t, err)

	err = validator.ValidateLicense(decoded, "orgId", "SKU", "instance-1", now)
	assert.NoError(t, err)

	err = validator.ValidateLicense(decoded, "INVALID", "SKU", "instance-1", now)
	assert.Error(t, err)

	err = validator.ValidateLicense(decoded, "orgId", "INVALID", "instance-1", now)
	assert.Error(t, err)

	err = validator.ValidateLicense(decoded, "orgId", "SKU", "INVALID", now)
	assert.Error(t, err)
}

func TestManager_ValidateLicense_Invalid(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	validator, err := NewValidatorFromBytes(certPEM)
	require.NoError(t, err)

	invalidEnvelope := &common.LicenseEnvelope{
		License: &common.License{
			ID:             uuid.NewString(),
			CreationTime:   time.Now().UTC().Format(time.RFC3339),
			ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
		},
		Signature: []byte("invalid-signature"),
	}

	err = validator.ValidateLicense(invalidEnvelope, "orgId", "", "", now)
	assert.Error(t, err)
}
