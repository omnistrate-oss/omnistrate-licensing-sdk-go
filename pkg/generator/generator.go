package generator

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/certificate"
	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/common"
)

type GeneratorInterface interface {
	GenerateLicense(orgId, productPlanUniqueID, instanceId, subscriptionId, description string, expirationDate time.Time) (envelope *common.LicenseEnvelope, err error)
	GenerateLicenseBase64(orgId, productPlanUniqueID, instanceId, subscriptionId, description string, expirationDate time.Time) (string, error)
	RenewLicense(envelope *common.LicenseEnvelope, expirationDate time.Time) (newEnvelope *common.LicenseEnvelope, err error)
	RenewLicenseBase64(envelopeBase64 string, expirationDate time.Time) (string, error)
	GetPublicCertificateBase64() string
}

type Manager struct {
	key     *rsa.PrivateKey
	certPEM []byte
}

func NewGenerator(key *rsa.PrivateKey, certPEM []byte) GeneratorInterface {
	return &Manager{
		key:     key,
		certPEM: certPEM,
	}
}

func NewGeneratorFromBytes(keyPEM []byte, certPEM []byte) (GeneratorInterface, error) {
	key, err := certificate.LoadPrivateKeyFromBytes(keyPEM)
	if err != nil {
		return nil, err
	}
	return NewGenerator(key, certPEM), nil
}

func NewGeneratorFromFiles(keyPath string, certPath string) (GeneratorInterface, error) {
	key, err := certificate.LoadPrivateKey(keyPath)
	if err != nil {
		return nil, err
	}
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	return NewGenerator(key, certPEM), nil
}

func NewGeneratorFromConfig(config *GeneratorConfig) (GeneratorInterface, error) {
	if config == nil {
		config = NewGeneratorConfigFromEnv()
	}
	return NewGeneratorFromFiles(config.KeyPath, config.CertPath)
}

func NewGeneratorFromEnv() (GeneratorInterface, error) {
	config := NewGeneratorConfigFromEnv()
	return NewGeneratorFromConfig(config)
}

func (m *Manager) GetPublicCertificateBase64() string {
	return base64.StdEncoding.EncodeToString(m.certPEM)
}

func (m *Manager) GenerateLicense(
	orgId string,
	productPlanUniqueID string,
	instanceId string,
	subscriptionId string,
	description string,
	expirationDate time.Time,
) (
	envelope *common.LicenseEnvelope,
	err error,
) {
	if m.key == nil {
		return nil, fmt.Errorf("licenseKey is required to generate a license")
	}

	now := time.Now().UTC()
	license := common.NewLicense(orgId, productPlanUniqueID, instanceId, subscriptionId, description, now, expirationDate)

	// Sign with private key
	licenseBytes, err := license.Bytes()
	if err != nil {
		return nil, err
	}
	signature, err := certificate.Sign(m.key, licenseBytes)
	if err != nil {
		return nil, err
	}

	// Create envelope
	envelope = common.NewLicenseEnvelope(license, signature)
	return
}

func (m *Manager) GenerateLicenseBase64(
	orgId string,
	productPlanUniqueID string,
	instanceId string,
	subscriptionId string,
	description string,
	expirationDate time.Time,
) (
	string,
	error,
) {
	envelope, err := m.GenerateLicense(orgId, productPlanUniqueID, instanceId, subscriptionId, description, expirationDate)
	if err != nil {
		return "", err
	}

	return envelope.EncodeBase64(), nil
}

func (m *Manager) RenewLicense(
	envelope *common.LicenseEnvelope,
	expirationDate time.Time,
) (
	newEnvelope *common.LicenseEnvelope,
	err error,
) {
	if envelope == nil {
		return nil, fmt.Errorf("envelope is required")
	}

	if !envelope.IsValid() {
		return nil, fmt.Errorf("envelope is invalid")
	}

	// Extract the license
	license := envelope.License
	license.Renew(expirationDate)

	// Sign with private key
	licenseBytes, err := license.Bytes()
	if err != nil {
		return nil, err
	}
	signature, err := certificate.Sign(m.key, licenseBytes)
	if err != nil {
		return nil, err
	}

	// Create envelope
	newEnvelope = common.NewLicenseEnvelope(license, signature)
	return
}

func (m *Manager) RenewLicenseBase64(envelopeBase64 string, expirationDate time.Time) (string, error) {
	// Decode the license envelope
	envelope, err := common.DecodeLicenseEnvelopeFromBase64(envelopeBase64)
	if err != nil {
		return "", err
	}

	// Renew the license
	newEnvelope, err := m.RenewLicense(envelope, expirationDate)
	if err != nil {
		return "", err
	}

	return newEnvelope.EncodeBase64(), nil
}
