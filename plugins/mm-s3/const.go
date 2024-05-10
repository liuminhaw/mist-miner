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
)

type mmS3Error struct {
	category string
	code     string
}

func (e *mmS3Error) Error() string {
	return fmt.Sprintf("%s: %s", e.category, e.code)
}
