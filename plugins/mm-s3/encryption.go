package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type encryptionProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketEncryptionOutput
}

func (ep *encryptionProp) fetchConf() error {
	output, err := ep.client.GetBucketEncryption(context.Background(), &s3.GetBucketEncryptionInput{
		Bucket: ep.bucket.Name,
	})
	if err != nil {
		return fmt.Errorf("fetchConf encryption: %w", err)
	}

	ep.configurations = output
	return nil
}

func (ep *encryptionProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := ep.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate encryptionProp: %w", err)
	}

	for _, rule := range ep.configurations.ServerSideEncryptionConfiguration.Rules {
		property := shared.MinerProperty{
			Type: encryption,
			Label: shared.MinerPropertyLabel{
				Name:   "Rule",
				Unique: false,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(rule); err != nil {
			return nil, fmt.Errorf("generate encryptionProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
