package validator

import (
	"os"
	"time"
)

type ValidationOptions struct {
	SkipCertificateValidation bool
	CertificateDomain         string
	CurrentTime               time.Time
	CertPath                  string
	LicensePath               string
	OrganizationID            string
	ProductPlanUniqueID       string
	InstanceID                string
}

func ValidateLicense(orgId, sku string) (err error) {
	return ValidateLicenseWithOptions(ValidationOptions{
		OrganizationID:      orgId,
		ProductPlanUniqueID: sku,
	})
}

func ValidateLicenseWithOptions(
	options ValidationOptions,
) (err error) {
	config := NewValidatorConfigFromEnv()

	if options.CertPath != "" {
		config.CertPath = options.CertPath
	}

	if options.LicensePath != "" {
		config.LicensePath = options.LicensePath
	}

	licenseBytes, err := os.ReadFile(config.LicensePath)
	if err != nil {
		return err
	}

	validator, err := NewValidatorFromConfig(config)
	if err != nil {
		return err
	}

	currentTime := time.Now().UTC()
	if !options.CurrentTime.IsZero() {
		currentTime = options.CurrentTime
	}

	if !options.SkipCertificateValidation {
		certificateDomain := options.CertificateDomain
		if certificateDomain == "" {
			certificateDomain = signingCertificateValidDnsName
		}

		err = validator.ValidateCertificate(certificateDomain, currentTime)
		if err != nil {
			return
		}
	}

	instanceID := options.InstanceID
	if options.InstanceID == "" {
		instanceID = config.InstanceID
	}

	err = validator.ValidateLicenseBytes(licenseBytes, options.OrganizationID, options.ProductPlanUniqueID, instanceID, currentTime)
	if err != nil {
		return err
	}

	return nil
}
