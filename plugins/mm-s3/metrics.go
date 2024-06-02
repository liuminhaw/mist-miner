package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type metricsProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.ListBucketMetricsConfigurationsOutput
	requestToken   string
}

func (mp *metricsProp) fetchConf() error {
	output, err := mp.client.ListBucketMetricsConfigurations(
		context.Background(),
		&s3.ListBucketMetricsConfigurationsInput{
			Bucket:            mp.bucket.Name,
			ContinuationToken: &mp.requestToken,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf metrics: %w", err)
	}

	mp.configurations = output
	return nil
}

func (mp *metricsProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	mp.requestToken = ""
	for {
		if err := mp.fetchConf(); err != nil {
			return nil, fmt.Errorf("generate metrics: %w", err)
		}

		for _, config := range mp.configurations.MetricsConfigurationList {
			property := shared.MinerProperty{
				Type: metrics,
				Label: shared.MinerPropertyLabel{
					Name:   *config.Id,
					Unique: true,
				},
				Content: shared.MinerPropertyContent{
					Format: shared.FormatJson,
				},
			}
			if err := property.FormatContentValue(config); err != nil {
				return nil, fmt.Errorf("generate metrics: %w", err)
			}

			properties = append(properties, property)
		}

		if *mp.configurations.IsTruncated {
			mp.requestToken = *mp.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
