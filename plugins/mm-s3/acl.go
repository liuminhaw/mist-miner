package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type aclProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketAclOutput
}

func (aclP *aclProp) fetchConf() error {
	output, err := aclP.client.GetBucketAcl(
		context.Background(),
		&s3.GetBucketAclInput{
			Bucket: aclP.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("getAclProperties: %w", err)
	}

	aclP.configurations = output
	return nil
}

func (aclP *aclProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := aclP.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate aclProp: %w", err)
	}

	properties = append(properties, shared.MinerProperty{
		Type: acl,
		Label: shared.MinerPropertyLabel{
			Name:   "Owner",
			Unique: true,
		},
		Content: shared.MinerPropertyContent{
			Format: formatText,
			Value:  *aclP.configurations.Owner.ID,
		},
	})
	for _, grant := range aclP.configurations.Grants {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(grant); err != nil {
			return nil, fmt.Errorf("generate aclProp: marshal grant: %w", err)
		}
		grantValue := buffer.Bytes()

		properties = append(properties, shared.MinerProperty{
			Type: acl,
			Label: shared.MinerPropertyLabel{
				Name:   "Grantee",
				Unique: false,
			},
			Content: shared.MinerPropertyContent{
				Format: formatJson,
				Value:  string(grantValue),
			},
		})
	}

	return properties, nil
}
