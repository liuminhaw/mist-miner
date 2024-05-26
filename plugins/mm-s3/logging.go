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

type loggingProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketLoggingOutput
}

func (lp *loggingProp) fetchConf() error {
	output, err := lp.client.GetBucketLogging(
		context.Background(),
		&s3.GetBucketLoggingInput{
			Bucket: lp.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf loggingProp: %w", err)
	}

	lp.configurations = output
	return nil
}

func (lp *loggingProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := lp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate loggingProp: %w", err)
	}

	if lp.configurations.LoggingEnabled != nil {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(lp.configurations.LoggingEnabled); err != nil {
			return nil, fmt.Errorf("generate loggingProp: marshal logging: %w", err)
		}
		loggingValue := buffer.Bytes()

		properties = append(properties, shared.MinerProperty{
			Type: logging,
			Label: shared.MinerPropertyLabel{
				Name:   "Logging",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: formatJson,
				Value:  string(loggingValue),
			},
		})
	}

	return properties, nil
}
