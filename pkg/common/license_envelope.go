package common

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type LicenseEnvelope struct {
	License   *License `json:"License"`
	Signature []byte   `json:"Signature"`
}

func NewLicenseEnvelope(license *License, signature []byte) *LicenseEnvelope {
	return &LicenseEnvelope{
		License:   license,
		Signature: signature,
	}
}

func (le *LicenseEnvelope) IsValid() bool {
	if le.License == nil || len(le.Signature) == 0 {
		return false
	}
	return true
}

func (le *LicenseEnvelope) IsExpired(currentTime time.Time) bool {
	if !le.IsValid() {
		return true
	}
	return le.License.IsExpiredAt(currentTime)
}

func (le *LicenseEnvelope) Bytes() ([]byte, error) {
	return json.Marshal(le)
}

func (le *LicenseEnvelope) String() string {
	jsonBytes, _ := le.Bytes()
	return string(jsonBytes)
}

func (le *LicenseEnvelope) EncodeBase64() string {
	jsonBytes, _ := le.Bytes()
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

func DecodeLicenseEnvelopeFromBytes(data []byte) (*LicenseEnvelope, error) {
	le := &LicenseEnvelope{}
	err := json.Unmarshal(data, le)
	if err != nil {
		return nil, err
	}
	return le, nil
}

func DecodeLicenseEnvelopeFromString(data string) (*LicenseEnvelope, error) {
	return DecodeLicenseEnvelopeFromBytes([]byte(data))
}

func DecodeLicenseEnvelopeFromBase64(data string) (*LicenseEnvelope, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return DecodeLicenseEnvelopeFromBytes(decoded)
}
