package validator

import "os"

type ValidatorConfig struct {
	CertPath    string `json:"certPath"`
	LicensePath string `json:"licensePath"`
	InstanceID  string `json:"instanceID"`
}

func NewValidatorConfig(instanceID, certPath, licensePath string) *ValidatorConfig {
	if certPath == "" {
		certPath = defaultValidatorCertPath
	}
	if licensePath == "" {
		licensePath = defaultValidatorLicensePath
	}
	return &ValidatorConfig{
		InstanceID:  instanceID,
		CertPath:    certPath,
		LicensePath: licensePath,
	}
}

func NewValidatorConfigFromEnv() *ValidatorConfig {
	return NewValidatorConfig(
		os.Getenv(licenseInstanceIDEnv),
		os.Getenv(licenseCertPathEnv),
		os.Getenv(licenseFilePathEnv),
	)
}

func (c *ValidatorConfig) IsValid() bool {
	if c.CertPath == "" || c.LicensePath == "" {
		return false
	}
	return true
}
