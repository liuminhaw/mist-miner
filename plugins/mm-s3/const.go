package main

import "fmt"

const (
	accelerateConfig   = "AccelerateConfig"
	acl                = "Acl"
	cors               = "CORS"
	encryption         = "Encryption"
	intelligentTiering = "IntelligentTiering"
	inventory          = "Inventory"
	lifecycle          = "Lifecycle"
	tagging            = "Tag"
	noConfig           = "NoConfiguration"
	valueSeparator     = "|"

	formatJson = "json"
	formatText = "text"
)

var miningProperties = []string{
	accelerateConfig,
	acl,
    cors,
    encryption,
	intelligentTiering,
    inventory,
	lifecycle,
	tagging,
}

type mmS3Error struct {
	category string
	code     string
}

func (e *mmS3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.category, e.code)
}
