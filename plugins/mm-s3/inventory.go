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

type inventoryProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.ListBucketInventoryConfigurationsOutput
	requestToken   string
}

func (ip *inventoryProp) fetchConf() error {
	output, err := ip.client.ListBucketInventoryConfigurations(
		context.Background(),
		&s3.ListBucketInventoryConfigurationsInput{
			Bucket:            ip.bucket.Name,
			ContinuationToken: &ip.requestToken,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf inventory: %w", err)
	}

	ip.configurations = output
	return nil
}

func (ip *inventoryProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	ip.requestToken = ""
	for {
		if err := ip.fetchConf(); err != nil {
			return nil, fmt.Errorf("generate invnetory: %w", err)
		}

		for _, config := range ip.configurations.InventoryConfigurationList {
			buffer := new(bytes.Buffer)
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			if err := encoder.Encode(config); err != nil {
				return nil, fmt.Errorf("generate inventory: marshal config: %w", err)
			}
			configValue := buffer.Bytes()

			properties = append(properties, shared.MinerProperty{
				Type: inventory,
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

		if *ip.configurations.IsTruncated {
			ip.requestToken = *ip.configurations.NextContinuationToken
		} else {
			break
		}
	}

	return properties, nil
}
