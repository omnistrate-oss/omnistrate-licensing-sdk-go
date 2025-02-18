package generator

import "os"

type GeneratorConfig struct {
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

func NewGeneratorConfig(certPath, keyPath string) *GeneratorConfig {
	if certPath == "" {
		certPath = defaultGeneratorCertPath
	}
	if keyPath == "" {
		keyPath = defaultGeneratorKeyPath
	}
	return &GeneratorConfig{
		CertPath: certPath,
		KeyPath:  keyPath,
	}
}

func NewGeneratorConfigFromEnv() *GeneratorConfig {
	return NewGeneratorConfig(os.Getenv(licenseCertPathEnv), os.Getenv(licenseKeyPathEnv))
}

func (c *GeneratorConfig) IsValid() bool {
	if c.CertPath == "" || c.KeyPath == "" {
		return false
	}
	return true
}
