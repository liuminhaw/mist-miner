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
			buffer := new(bytes.Buffer)
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(config); err != nil {
				return nil, fmt.Errorf("generate metrics: marshal config: %w", err)
			}
			configValue := buffer.Bytes()

			properties = append(properties, shared.MinerProperty{
				Type: metrics,
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

		if *mp.configurations.IsTruncated {
			mp.requestToken = *mp.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
