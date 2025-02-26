package generator

import (
	_ "embed"
	"encoding/base64"
	"testing"
	"time"

	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed certificate-test-tls.crt
var certPEM []byte

//go:embed certificate-test-tls.key
var keyPEM []byte

func TestNewGenerator(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)
	require.NotNil(t, manager)
}

func TestGenerator_GenerateLicense(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	expiration := time.Now().UTC().Add(48 * time.Hour)
	envelope, err := manager.GenerateLicense("orgId", "SKU", "instance-1", "subs-1", "product a", expiration)
	assert.NoError(t, err)
	assert.NotNil(t, envelope)
	assert.NotEmpty(t, envelope.Signature)
	licenseExpiration, _ := envelope.License.GetExpirationTime()
	assert.Equal(t, expiration.Format(time.RFC3339), licenseExpiration.Format(time.RFC3339))
}

func TestGenerator_GenerateLicenseBase64(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	licenseBase64, err := manager.GenerateLicenseBase64("orgId", "SKU", "instance-1", "subs-1", "product a", time.Now().UTC().Add(48*time.Hour))
	assert.NoError(t, err)
	assert.NotEmpty(t, licenseBase64)
}
func TestGenerator_RenewLicense(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	envelope, err := manager.GenerateLicense("orgId", "SKU", "instance-1", "subs-1", "product a", time.Now().UTC().Add(48*time.Hour))
	assert.NoError(t, err)
	assert.NotNil(t, envelope)
	oldExpirationTime, _ := envelope.License.GetExpirationTime()

	newEnvelope, err := manager.RenewLicense(envelope, time.Now().UTC().Add(49*time.Hour))
	assert.NoError(t, err)
	assert.NotNil(t, newEnvelope)
	newExpirationTime, _ := newEnvelope.License.GetExpirationTime()

	assert.Equal(t, envelope.License.ID, newEnvelope.License.ID)
	assert.Greater(t, newExpirationTime.Unix(), oldExpirationTime.Unix())
}

func TestGenerator_RenewLicenseBase64(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	licenseBase64, err := manager.GenerateLicenseBase64("orgId", "SKU", "instance-1", "subs-1", "product a", time.Now().UTC().Add(48*time.Hour))
	assert.NoError(t, err)
	assert.NotEmpty(t, licenseBase64)

	renewLicenseBase64, err := manager.RenewLicenseBase64(licenseBase64, time.Now().UTC().Add(49*time.Hour))
	assert.NoError(t, err)
	assert.NotEmpty(t, renewLicenseBase64)

	decoded, err := common.DecodeLicenseEnvelopeFromBase64(licenseBase64)
	require.NoError(t, err)
	oldExpirationTime, _ := decoded.License.GetExpirationTime()

	decodedRenewed, err := common.DecodeLicenseEnvelopeFromBase64(renewLicenseBase64)
	require.NoError(t, err)
	newExpirationTime, _ := decodedRenewed.License.GetExpirationTime()

	assert.Equal(t, decoded.License.ID, decodedRenewed.License.ID)
	assert.Greater(t, newExpirationTime.Unix(), oldExpirationTime.Unix())
}

func TestGenerator_GetPublicCertificateBase64(t *testing.T) {
	t.Parallel()

	manager, err := NewGeneratorFromBytes(keyPEM, certPEM)
	require.NoError(t, err)

	certBase64 := manager.GetPublicCertificateBase64()
	assert.NotEmpty(t, certBase64)
	assert.Equal(t, base64.StdEncoding.EncodeToString(certPEM), certBase64)
	pem, err := base64.StdEncoding.DecodeString(certBase64)
	require.NoError(t, err)
	assert.Equal(t, certPEM, pem)
}
