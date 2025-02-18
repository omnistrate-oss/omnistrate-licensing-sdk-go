package validator

const (
	licenseCertPathEnv   = "SERVICE_PLAN_SUBSCRIPTION_LICENSE_CERT_PATH"
	licenseFilePathEnv   = "SERVICE_PLAN_SUBSCRIPTION_LICENSE_FILE_PATH"
	licenseInstanceIDEnv = "INSTANCE_ID"

	defaultValidatorCertPath    = "/var/subscription/license.crt"
	defaultValidatorLicensePath = "/var/subscription/license.lic"

	validDnsName = "licensing.omnistrate.cloud"
)
