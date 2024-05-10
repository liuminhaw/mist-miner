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
	var properties []shared.MinerProperty

	contToken := ""
	for {
		output, err := client.ListBucketIntelligentTieringConfigurations(
			context.Background(),
			&s3.ListBucketIntelligentTieringConfigurationsInput{
				Bucket:            bucket.Name,
				ContinuationToken: &contToken,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("getIntelligentTieringProperties: %w", err)
		}

		for _, config := range output.IntelligentTieringConfigurationList {
			buffer := new(bytes.Buffer)
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(config); err != nil {
				return nil, fmt.Errorf(
					"getIntelligentTieringProperties: marshal IntelligentTiering config: %w",
					err,
				)
			}
			configValue := buffer.Bytes()

			properties = append(properties, shared.MinerProperty{
				Type: intelligentTiering,
				Label: shared.MinerPropertyLabel{
					Name:   *config.Id,
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: "json",
					Value:  string(configValue),
				},
			})
		}

		if *output.IsTruncated {
			contToken = *output.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
