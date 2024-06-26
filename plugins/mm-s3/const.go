package main

import "fmt"

const (
	accelerateConfig   = "AccelerateConfig"
	analyticsConfig    = "AnalyticsConfig"
	acl                = "Acl"
	cors               = "CORS"
	encryption         = "Encryption"
	intelligentTiering = "IntelligentTiering"
	inventory          = "Inventory"
	lifecycle          = "Lifecycle"
	location           = "Location"
	logging            = "Logging"
	metrics            = "Metrics"
	notification       = "Notification"
	ownershipControl   = "OwnershipControl"
	policy             = "Policy"
	policyStatus       = "PolicyStatus"
	replication        = "Replication"
	requestPayment     = "RequestPayment"
	tagging            = "Tag"
	versioning         = "Versioning"
	website            = "Website"

	noConfig       = "NoConfiguration"
	valueSeparator = "|"
)

var miningProperties = []string{
	accelerateConfig,
	analyticsConfig,
	acl,
	cors,
	encryption,
	intelligentTiering,
	inventory,
	lifecycle,
	logging,
	metrics,
	notification,
	ownershipControl,
	policy,
	policyStatus,
	replication,
	requestPayment,
	tagging,
	versioning,
	website,
}

type mmS3Error struct {
	category string
	code     string
}

func (e *mmS3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.category, e.code)
}
