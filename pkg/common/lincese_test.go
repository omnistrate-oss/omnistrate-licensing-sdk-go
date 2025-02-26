package common

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLicense(t *testing.T) {
	t.Parallel()

	creationTime := time.Now().UTC()
	expirationTime := creationTime.Add(24 * time.Hour)
	license := NewLicense("SKU", "instance-id", "sub-id", "desc", creationTime, expirationTime)

	assert.NotEmpty(t, license.ID)
	assert.Equal(t, "SKU", license.ProductPlanUniqueId)
	assert.Equal(t, "desc", license.Description)
	assert.Equal(t, "sub-id", license.SubscriptionID)
	assert.Equal(t, "instance-id", license.InstanceID)
	assert.Equal(t, creationTime.Format(time.RFC3339), license.CreationTime)
	assert.Equal(t, expirationTime.Format(time.RFC3339), license.ExpirationTime)
	assert.Equal(t, uint64(1), license.Version)
}

func TestGetExpirationTime(t *testing.T) {
	t.Parallel()

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}

	expirationTime, err := license.GetExpirationTime()
	require.NoError(t, err)
	assert.Equal(t, "2022-01-01T00:00:00Z", expirationTime.Format(time.RFC3339))
}

func TestRenew(t *testing.T) {
	t.Parallel()

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	expirationTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	license.Renew(expirationTime)

	assert.Equal(t, expirationTime.Format(time.RFC3339), license.ExpirationTime)
}

func TestGetCreationTime(t *testing.T) {
	t.Parallel()

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}

	creationTime, err := license.GetCreationTime()
	require.NoError(t, err)
	assert.Equal(t, "2021-01-01T00:00:00Z", creationTime.Format(time.RFC3339))
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	validLicense := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	invalidLicense := License{
		ID:             "",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}

	err := validLicense.IsValid("", "")
	assert.NoError(t, err)

	err = invalidLicense.IsValid("", "")
	assert.Error(t, err)
}

func TestIsValidWithSKU(t *testing.T) {
	t.Parallel()

	validLicense := License{
		ID:                  "1234",
		ProductPlanUniqueId: "SKU",
		CreationTime:        "2021-01-01T00:00:00Z",
		ExpirationTime:      "2022-01-01T00:00:00Z",
	}
	invalidLicense := License{
		ID:                  "",
		ProductPlanUniqueId: "SKU",
		CreationTime:        "2021-01-01T00:00:00Z",
		ExpirationTime:      "2022-01-01T00:00:00Z",
	}

	err := validLicense.IsValid("", "")
	assert.NoError(t, err)

	err = validLicense.IsValid("SKU", "")
	assert.NoError(t, err)

	err = validLicense.IsValid("INVALID", "")
	assert.Error(t, err)

	err = invalidLicense.IsValid("SKU", "")
	assert.Error(t, err)
}

func TestIsValidWithInstanceID(t *testing.T) {
	t.Parallel()

	validLicense := License{
		ID:                  "1234",
		ProductPlanUniqueId: "SKU",
		InstanceID:          "instance-id",
		CreationTime:        "2021-01-01T00:00:00Z",
		ExpirationTime:      "2022-01-01T00:00:00Z",
	}
	invalidLicense := License{
		ID:                  "",
		ProductPlanUniqueId: "SKU",
		InstanceID:          "instance-id",
		CreationTime:        "2021-01-01T00:00:00Z",
		ExpirationTime:      "2022-01-01T00:00:00Z",
	}

	err := validLicense.IsValid("", "")
	assert.NoError(t, err)

	err = validLicense.IsValid("SKU", "instance-id")
	assert.NoError(t, err)

	err = invalidLicense.IsValid("SKU", "INVALID")
	assert.Error(t, err)

	err = invalidLicense.IsValid("SKU", "instance-id")
	assert.Error(t, err)
}

func TestIsExpired(t *testing.T) {
	t.Parallel()

	expiredLicense := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2021-01-01T00:00:00Z",
	}
	validLicense := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2100-01-01T00:00:00Z",
	}

	assert.True(t, expiredLicense.IsExpired())
	assert.False(t, validLicense.IsExpired())
}

func TestBytes(t *testing.T) {
	t.Parallel()

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	jsonBytes, err := license.Bytes()
	require.NoError(t, err)

	expectedBytes, err := json.Marshal(license)
	require.NoError(t, err)
	assert.Equal(t, expectedBytes, jsonBytes)
}

func TestString(t *testing.T) {
	t.Parallel()

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	jsonString := license.String()

	expectedString, err := json.Marshal(license)
	require.NoError(t, err)
	assert.Equal(t, string(expectedString), jsonString)
}

func TestFromString(t *testing.T) {
	t.Parallel()

	license := License{}
	jsonString := `{"ID":"1234","creationTime":"2021-01-01T00:00:00Z","expirationTime":"2022-01-01T00:00:00Z"}`
	err := license.FromString(jsonString)
	require.NoError(t, err)

	expectedLicense := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	assert.Equal(t, expectedLicense, license)
}

func TestLicenseSerialization(t *testing.T) {
	t.Parallel()

	require := require.New(t)
	assert := assert.New(t)

	license := License{
		ID:             "1234",
		CreationTime:   "2021-01-01T00:00:00Z",
		ExpirationTime: "2022-01-01T00:00:00Z",
	}
	jsonBytes, err := json.Marshal(license)
	require.NoError(err, "failed to marshal license")

	license2 := License{}
	err = json.Unmarshal(jsonBytes, &license2)
	require.NoError(err, "failed to unmarshal license")

	assert.Equal(license, license2, "license should be the same after serialization")
}
