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

// getIntelligentTieringProperties returns the intelligent tiering properties of a bucket
func getIntelligentTieringProperties(
	client *s3.Client,
	bucket *types.Bucket,
) ([]shared.MinerProperty, error) {
	propsCrawler, err := New(client, bucket, intelligentTiering)
	if err != nil {
		return nil, fmt.Errorf("getIntelligentTieringProperties: %w", err)
	}
	// Check if returned props type is intelligentTieringProp
	if _, ok := propsCrawler.(*intelligentTieringProp); !ok {
		return nil, fmt.Errorf("getIntelligentTieringProperties: unexpected crawler type")
	}

	return propsCrawler.generate()
}

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
			buffer := new(bytes.Buffer)
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(config); err != nil {
				return nil, fmt.Errorf("generate intelligentTiering: marshal config: %w", err)
			}
			configValue := buffer.Bytes()

			properties = append(properties, shared.MinerProperty{
				Type: intelligentTiering,
				Label: shared.MinerPropertyLabel{
					Name:   *config.Id,
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: formatJson,
					Value:  string(configValue),
				},
			})
		}

		if *it.configurations.IsTruncated {
			it.requestToken = *it.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
