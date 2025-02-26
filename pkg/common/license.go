package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type License struct {
	ID                  string `json:"ID,omitempty"`
	CreationTime        string `json:"CreationTime,omitempty"`
	ExpirationTime      string `json:"ExpirationTime,omitempty"`
	Description         string `json:"Description,omitempty"`
	InstanceID          string `json:"InstanceID,omitempty"`
	SubscriptionID      string `json:"SubscriptionID,omitempty"`
	ProductPlanUniqueID string `json:"ProductPlanUniqueID,omitempty"`
	OrganizationID      string `json:"OrganizationID,omitempty"`
	Version             uint64 `json:"Version,omitempty"`
}

func NewLicense(orgID, productPlanUniqueID, instanceID, subscriptionID, description string, creationTime, expirationTime time.Time) *License {
	return &License{
		ID:                  uuid.NewString(),
		CreationTime:        creationTime.UTC().Format(time.RFC3339),
		ExpirationTime:      expirationTime.UTC().Format(time.RFC3339),
		InstanceID:          instanceID,
		SubscriptionID:      subscriptionID,
		Description:         description,
		ProductPlanUniqueID: productPlanUniqueID,
		OrganizationID:      orgID,
		Version:             1,
	}
}

func (l *License) GetExpirationTime() (time.Time, error) {
	return time.Parse(time.RFC3339, l.ExpirationTime)
}

func (l *License) GetCreationTime() (time.Time, error) {
	return time.Parse(time.RFC3339, l.CreationTime)
}

func (l *License) IsValid(orgID, productPlanUniqueID, instanceID string) error {
	if l.ID == "" || l.CreationTime == "" || l.ExpirationTime == "" {
		return errors.New("missing required fields")
	}
	if orgID != "" && l.OrganizationID != orgID {
		return fmt.Errorf("invalid organization id %s, expected %s", orgID, l.OrganizationID)
	}
	if productPlanUniqueID != "" && l.ProductPlanUniqueID != productPlanUniqueID {
		return fmt.Errorf("invalid product unique id %s, expected %s", productPlanUniqueID, l.ProductPlanUniqueID)
	}
	if instanceID != "" && l.InstanceID != instanceID {
		return fmt.Errorf("invalid instance id: %s, expected %s", instanceID, l.InstanceID)
	}
	if _, err := l.GetCreationTime(); err != nil {
		return errors.Wrap(err, "invalid creation time")
	}
	if _, err := l.GetExpirationTime(); err != nil {
		return errors.Wrap(err, "invalid expiration time")
	}
	return nil
}

func (l *License) IsExpired() bool {
	currentTime := time.Now().UTC()
	return l.IsExpiredAt(currentTime.UTC())
}

func (l *License) IsExpiredAt(t time.Time) bool {
	expirationTime, _ := l.GetExpirationTime()
	return t.UTC().After(expirationTime)
}

func (l *License) Renew(expirationTime time.Time) {
	l.CreationTime = time.Now().UTC().Format(time.RFC3339)
	l.ExpirationTime = expirationTime.UTC().Format(time.RFC3339)
	l.Version++
}

func (l *License) Bytes() ([]byte, error) {
	return json.Marshal(l)
}

func (l *License) String() string {
	jsonBytes, _ := l.Bytes()
	return string(jsonBytes)
}

func (l *License) FromString(s string) error {
	return json.Unmarshal([]byte(s), l)
}
