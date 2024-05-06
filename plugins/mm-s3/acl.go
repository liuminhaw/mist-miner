package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type aclGrantee struct {
	displayName string
	granteeType string
	id          string
	permission  string
}

// getAclProperties retrieves the ACL properties of a bucket.
// Grantee value is combined with displayName, granteeType, id, and permission separated by `valueSeparator`.
func getAclProperties(client *s3.Client, bucket *types.Bucket) ([]shared.MinerProperty, error) {
	output, err := client.GetBucketAcl(
		context.Background(),
		&s3.GetBucketAclInput{
			Bucket: bucket.Name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("getAclProperties: %w", err)
	}

	var properties []shared.MinerProperty
	properties = append(properties, shared.MinerProperty{
		Type:  acl,
		Name:  "Owner",
		Value: *output.Owner.ID,
	})

	for _, grant := range output.Grants {

		grantee := aclGrantee{
			permission: string(grant.Permission),
		}
		if grant.Grantee.DisplayName != nil {
			grantee.displayName = *grant.Grantee.DisplayName
		}
		switch grant.Grantee.Type {
		case types.TypeCanonicalUser:
			grantee.granteeType = string(types.TypeCanonicalUser)
			grantee.id = *grant.Grantee.ID
		case types.TypeGroup:
			grantee.granteeType = string(types.TypeGroup)
			grantee.id = *grant.Grantee.URI
		case types.TypeAmazonCustomerByEmail:
			grantee.granteeType = string(types.TypeAmazonCustomerByEmail)
			grantee.id = *grant.Grantee.EmailAddress
		default:
			return nil, fmt.Errorf("getAclProperties: unknown grantee type: %s", grant.Grantee.Type)
		}

		properties = append(properties, shared.MinerProperty{
			Type: acl,
			Name: "Grantee",
			Value: strings.Join(
				[]string{grantee.displayName, grantee.granteeType, grantee.id, grantee.permission},
				valueSeparator,
			),
		})
	}

	return properties, nil
}
