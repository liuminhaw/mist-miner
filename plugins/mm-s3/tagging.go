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

// getTaggingProperties retrieves the tagging properties of a bucket.
// With property type "Tag", property name as the key, and property value as the value.
func getTaggingProperties(client *s3.Client, bucket *types.Bucket) ([]shared.MinerProperty, error) {
	output, err := client.GetBucketTagging(
		context.Background(),
		&s3.GetBucketTaggingInput{
			Bucket: bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchTagSet":
				return nil, &mmS3Error{tagging, noConfig}
			default:
				return nil, fmt.Errorf("getTaggingProperties: %w", err)
			}
		} else {
			return nil, fmt.Errorf("getTaggingProperties: %w", err)
		}
	}

	var properties []shared.MinerProperty
	for _, tag := range output.TagSet {
		properties = append(properties, shared.MinerProperty{
			Type: tagging,
			Label: shared.MinerPropertyLabel{
				Name:   *tag.Key,
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: "string",
				Value:  *tag.Value,
			},
		})
	}

	return properties, nil
}
