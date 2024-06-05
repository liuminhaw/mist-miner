package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type versioningProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketVersioningOutput
}

func (vp *versioningProp) fetchConf() error {
	output, err := vp.client.GetBucketVersioning(
		context.Background(),
		&s3.GetBucketVersioningInput{
			Bucket: vp.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchCont versioning: %w", err)
	}

	vp.configurations = output
	return nil
}

func (vp *versioningProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := vp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate versioning: %w", err)
	}

	property := shared.MinerProperty{
		Type: versioning,
		Label: shared.MinerPropertyLabel{
			Name:   "Versioning",
			Unique: true,
		},
		Content: shared.MinerPropertyContent{
			Format: shared.FormatJson,
		},
	}
	if err := property.FormatContentValue(vp.configurations); err != nil {
		return nil, fmt.Errorf("generate versioning: %w", err)
	}

	properties = append(properties, property)
	return properties, nil
}
