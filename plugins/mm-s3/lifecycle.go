package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/liuminhaw/mist-miner/shared"
)

type lifecycleProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketLifecycleConfigurationOutput
}

func (l *lifecycleProp) fetchConf() error {
	output, err := l.client.GetBucketLifecycleConfiguration(
		context.Background(),
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: l.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchLifecycleConfiguration":
				return &mmS3Error{lifecycle, noConfig}
			default:
				return fmt.Errorf("fetchConf lifecycle: %w", err)
			}
		}
		return fmt.Errorf("fetchConf lifecycle: %w", err)
	}

	l.configurations = output
	return nil
}

func (l *lifecycleProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := l.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate lifecycleProp: %w", err)
	}
	for _, rule := range l.configurations.Rules {
		property := shared.MinerProperty{
			Type: lifecycle,
			Label: shared.MinerPropertyLabel{
				Name:   *rule.ID,
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(rule); err != nil {
			return nil, fmt.Errorf("generate lifecycleProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
