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

type taggingProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketTaggingOutput
}

func (t *taggingProp) fetchConf() error {
	output, err := t.client.GetBucketTagging(
		context.Background(),
		&s3.GetBucketTaggingInput{
			Bucket: t.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchTagSet":
				return &mmS3Error{tagging, noConfig}
			default:
				return fmt.Errorf("fetchConf taggingProp: %w", err)
			}
		} else {
			return fmt.Errorf("fetchConf taggingProp: %w", err)
		}
	}

	t.configurations = output
	return nil
}

func (t *taggingProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := t.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate taggingProp: %w", err)
	}
	for _, tag := range t.configurations.TagSet {
		property := shared.MinerProperty{
			Type: tagging,
			Label: shared.MinerPropertyLabel{
				Name:   *tag.Key,
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatText,
			},
		}
		if err := property.FormatContentValue(*tag.Value); err != nil {
			return nil, fmt.Errorf("generate taggingProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
