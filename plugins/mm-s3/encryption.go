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
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(rule); err != nil {
			return nil, fmt.Errorf("generate encryptionProp: marshal rule: %w", err)
		}
		ruleValue := buffer.Bytes()

		properties = append(properties, shared.MinerProperty{
			Type: encryption,
			Label: shared.MinerPropertyLabel{
				Name:   "Rule",
				Unique: false,
			},
			Content: shared.MinerPropertyContent{
				Format: formatJson,
				Value:  string(ruleValue),
			},
		})
	}

	return properties, nil
}
