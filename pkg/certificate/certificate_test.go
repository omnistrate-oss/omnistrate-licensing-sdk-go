package certificate

import (
	_ "embed"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed certificate-test-tls.crt
var certPEM []byte

//go:embed certificate-test-tls.key
var keyPEM []byte

func TestLoadPublicCertificateFromBytes(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)

	cert, err := LoadCertificateFromBytes(certPEM)
	require.NoError(err)
	require.NotNil(cert)

	assert.Equal("Let's Encrypt", cert.Issuer.Organization[0])
	assert.Equal("licensing-test.omnistrate.dev", cert.DNSNames[0])
	assert.False(cert.IsCA)
	assert.NotNil(cert.PublicKey)

	// Test valid certificate
	err = VerifyCertificate(cert, "licensing-test.omnistrate.dev", cert.NotBefore.Add(cert.NotAfter.Sub(cert.NotBefore)/2))
	require.NoError(err)

	// Test expired certificate
	err = VerifyCertificate(cert, "licensing-test.omnistrate.dev", cert.NotAfter.Add(time.Hour))
	require.Error(err)

	// Test invalid DNS name
	err = VerifyCertificate(cert, "licensing.omnistrate.dev.invalid", cert.NotBefore.Add(cert.NotAfter.Sub(cert.NotBefore)/2))
	require.Error(err)
}

func TestLoadPrivateCertificateFromBytes(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	cert, err := LoadPrivateKeyFromBytes(keyPEM)
	require.NoError(err)
	require.NotNil(cert)
}

func TestSignAndVerify(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	cert, err := LoadCertificateFromBytes(certPEM)
	require.NoError(err)
	require.NotNil(cert)

	key, err := LoadPrivateKeyFromBytes(keyPEM)
	require.NoError(err)
	require.NotNil(key)

	signature, err := Sign(key, []byte("test"))
	require.NoError(err)

	signature2, err := Sign(key, []byte("test"))
	require.NoError(err)
	assert.Equal(signature, signature2)

	err = VerifySignature(cert, signature, []byte("test"))
	require.NoError(err)

	err = VerifySignature(cert, signature, []byte("test "))
	require.Error(err)

	err = VerifySignature(cert, signature, []byte("tesy"))
	require.Error(err)
}
