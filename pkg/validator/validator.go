package validator

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/certificate"
	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/common"
	"github.com/pkg/errors"
)

type ValidatorInterface interface {
	ValidateLicense(envelope *common.LicenseEnvelope, sku, instanceID string, currentTime time.Time) error
	ValidateLicenseString(envelopeJson string, sku, instanceID string, currentTime time.Time) error
	ValidateLicenseBytes(envelopeBytes []byte, sku, instanceID string, currentTime time.Time) error
	ValidateLicenseBase64(envelopeBase64 string, sku, instanceID string, currentTime time.Time) error
	ValidateCertificate(currentTime time.Time) error
}

type Validator struct {
	cert *x509.Certificate
}

func NewValidator(cert *x509.Certificate) ValidatorInterface {
	return &Validator{
		cert: cert,
	}
}

func NewValidatorFromBytes(certPEM []byte) (ValidatorInterface, error) {
	cert, err := certificate.LoadCertificateFromBytes(certPEM)
	if err != nil {
		return nil, err
	}
	return NewValidator(cert), nil
}

func NewValidatorFromFiles(certPath string) (ValidatorInterface, error) {
	cert, err := certificate.LoadCertificate(certPath)
	if err != nil {
		return nil, err
	}
	return NewValidator(cert), nil
}

func NewValidatorFromConfig(config *ValidatorConfig) (ValidatorInterface, error) {
	return NewValidatorFromFiles(config.CertPath)
}

func (m *Validator) ValidateLicense(envelope *common.LicenseEnvelope, sku, instanceID string, currentTime time.Time) error {
	if m.cert == nil {
		return fmt.Errorf("signingCertificate is required to validate a license")
	}

	if envelope == nil {
		return fmt.Errorf("envelope is required")
	}

	if !envelope.IsValid() {
		return fmt.Errorf("envelope is invalid")
	}

	// Extract the signature
	signature := envelope.Signature

	// Extract the license
	license := envelope.License
	licenseBytes, err := license.Bytes()
	if err != nil {
		return err
	}

	// Check if the license is valid
	err = license.IsValid(sku, instanceID)
	if err != nil {
		return errors.Wrap(err, "license is invalid")
	}

	// Check if the license is expired
	if license.IsExpiredAt(currentTime) {
		return fmt.Errorf("license is expired")
	}

	// Verify the signature
	err = certificate.VerifySignature(m.cert, signature, licenseBytes)
	if err != nil {
		return errors.Wrap(err, "failed to verify signature")
	}

	// License is valid
	return nil
}

func (m *Validator) ValidateLicenseBase64(envelopeBase64 string, sku, instanceID string, currentTime time.Time) error {
	// Decode the license envelope
	envelope, err := common.DecodeLicenseEnvelopeFromBase64(envelopeBase64)
	if err != nil {
		return err
	}

	// Validate the license
	return m.ValidateLicense(envelope, sku, instanceID, currentTime)
}

func (m *Validator) ValidateLicenseString(envelopeJson string, sku, instanceID string, currentTime time.Time) error {
	// Decode the license envelope
	envelope, err := common.DecodeLicenseEnvelopeFromString(envelopeJson)
	if err != nil {
		return err
	}

	// Validate the license
	return m.ValidateLicense(envelope, sku, instanceID, currentTime)
}

func (m *Validator) ValidateLicenseBytes(envelopeBytes []byte, sku, instanceID string, currentTime time.Time) error {
	// Decode the license envelope
	envelope, err := common.DecodeLicenseEnvelopeFromBytes(envelopeBytes)
	if err != nil {
		return err
	}

	// Validate the license
	return m.ValidateLicense(envelope, sku, instanceID, currentTime)
}

func (m *Validator) ValidateCertificate(currentTime time.Time) error {
	if m.cert == nil {
		return fmt.Errorf("signingCertificate is required to validate a certificate")
	}

	// Validate the certificate
	return certificate.VerifyCertificate(m.cert, validDnsName, currentTime)
}
