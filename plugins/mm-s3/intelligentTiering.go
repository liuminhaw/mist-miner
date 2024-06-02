package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

// intelligentTieringProp is a crawler for fetching s3 IntelligentTiering properties
type intelligentTieringProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.ListBucketIntelligentTieringConfigurationsOutput
	requestToken   string
}

// fetchConf fetches the IntelligentTiering configurations for the bucket
func (it *intelligentTieringProp) fetchConf() error {
	output, err := it.client.ListBucketIntelligentTieringConfigurations(
		context.Background(),
		&s3.ListBucketIntelligentTieringConfigurationsInput{
			Bucket:            it.bucket.Name,
			ContinuationToken: &it.requestToken,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf intelligentTiering: %w", err)
	}

	it.configurations = output
	return nil
}

// generate generates the IntelligentTiering properties in MinerProperty format
// to be returned to the main miner
func (it *intelligentTieringProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	it.requestToken = ""
	for {
		if err := it.fetchConf(); err != nil {
			return nil, fmt.Errorf("generate intelligentTiering: %w", err)
		}

		for _, config := range it.configurations.IntelligentTieringConfigurationList {
			property := shared.MinerProperty{
				Type: intelligentTiering,
				Label: shared.MinerPropertyLabel{
					Name:   *config.Id,
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: shared.FormatJson,
				},
			}
            if err := property.FormatContentValue(config); err != nil {
                return nil, fmt.Errorf("generate intelligentTiering: %w", err)
            }
			properties = append(properties, property)
		}

		if *it.configurations.IsTruncated {
			it.requestToken = *it.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
