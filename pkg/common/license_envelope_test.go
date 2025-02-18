package common

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLicenseEnvelope(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	assert.NotNil(t, le)
	assert.Equal(t, license, le.License)
	assert.Equal(t, signature, le.Signature)
}

func TestLicenseEnvelope_IsValid(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	assert.True(t, le.IsValid())

	le.Signature = []byte{}
	assert.False(t, le.IsValid())

	le.License = nil
	assert.False(t, le.IsValid())
}

func TestLicenseEnvelope_IsExpired(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	assert.False(t, le.IsExpired(time.Now().UTC()))

	license.ExpirationTime = time.Now().AddDate(0, 0, -1).UTC().Format(time.RFC3339)
	assert.True(t, le.IsExpired(time.Now().UTC()))
}

func TestLicenseEnvelope_Bytes(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	bytes, err := le.Bytes()
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
}

func TestLicenseEnvelope_String(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	str := le.String()
	assert.NotEmpty(t, str)
}

func TestLicenseEnvelope_EncodeBase64(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	encoded := le.EncodeBase64()
	assert.NotEmpty(t, encoded)

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
}

func TestDecodeLicenseEnvelopeFromBytes(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	bytes, err := le.Bytes()
	assert.NoError(t, err)

	decodedLe, err := DecodeLicenseEnvelopeFromBytes(bytes)
	assert.NoError(t, err)
	assert.NotNil(t, decodedLe)
	assert.Equal(t, le.License, decodedLe.License)
	assert.Equal(t, le.Signature, decodedLe.Signature)
}

func TestDecodeLicenseEnvelopeFromString(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	str := le.String()

	decodedLe, err := DecodeLicenseEnvelopeFromString(str)
	assert.NoError(t, err)
	assert.NotNil(t, decodedLe)
	assert.Equal(t, le.License, decodedLe.License)
	assert.Equal(t, le.Signature, decodedLe.Signature)
}

func TestDecodeLicenseEnvelopeFromBase64(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("r\x7f\xe0E5\xbdx\x9d0ãz?\xaaU\t\xfd\xe1\xbf-Y\x82\xed\xfeI3\x99\x8f\xb0ɡj\xfb\x93\xbf\xc62\x1d\xaf;\xf8N\xd0B\x95\xd2\x12\x84n@\xe8\xd5\x15ac&\xa8HCQ\x8f\x1f\xe4\t")
	le := NewLicenseEnvelope(license, signature)

	encoded := le.EncodeBase64()

	decodedLe, err := DecodeLicenseEnvelopeFromBase64(encoded)
	assert.NoError(t, err)
	assert.NotNil(t, decodedLe)
	assert.Equal(t, le.License, decodedLe.License)
	assert.Equal(t, le.Signature, decodedLe.Signature)
}

func TestGetSignature(t *testing.T) {
	t.Parallel()

	license := &License{
		ID:             "12345",
		CreationTime:   time.Now().UTC().Format(time.RFC3339),
		ExpirationTime: time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339),
	}
	signature := []byte("test-signature")
	le := NewLicenseEnvelope(license, signature)

	assert.Equal(t, signature, le.Signature)
}
