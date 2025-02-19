package example

import (
	"testing"
	"time"

	"github.com/omnistrate-oss/omnistrate-licensing-sdk-go/pkg/validator"
	"github.com/stretchr/testify/require"
)

// TestValidateExample is a test function that demonstrates how to test the Validation function
// from the validator package.

// The test function uses options to ensure the certificate validation is skipped and old time is set
// In production code, ValidateLicense or ValidateLicenseForProduct can be used instead of ValidateLicenseWithOptions
func TestValidateExample(t *testing.T) {
	require := require.New(t)

	err := validator.ValidateLicenseWithOptions(validator.ValidationOptions{
		CertificateDomain: "licensing.omnistrate.dev", // test certificate
		CurrentTime:       time.Date(2025, 2, 19, 0, 0, 0, 0, time.UTC),
		CertPath:          "license.crt",
		LicensePath:       "license.lic",
		InstanceID:        "instance-jzxo986k2",
	})

	require.NoError(err)
}
