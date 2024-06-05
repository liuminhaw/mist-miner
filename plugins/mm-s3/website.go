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

type websiteProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketWebsiteOutput
}

func (wp *websiteProp) fetchConf() error {
	output, err := wp.client.GetBucketWebsite(
		context.Background(),
		&s3.GetBucketWebsiteInput{
			Bucket: wp.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchWebsiteConfiguration":
				return &mmS3Error{website, noConfig}
			default:
				return fmt.Errorf("fetchCont website: %w", err)
			}
		}
		return fmt.Errorf("fetchCont website: %w", err)
	}

	wp.configurations = output
	return nil
}

func (wp *websiteProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := wp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate website: %w", err)
	}

	property := shared.MinerProperty{
		Type: website,
		Label: shared.MinerPropertyLabel{
			Name:   "Website",
			Unique: true,
		},
		Content: shared.MinerPropertyContent{
			Format: shared.FormatJson,
		},
	}
	if err := property.FormatContentValue(wp.configurations); err != nil {
		return nil, fmt.Errorf("generate website: %w", err)
	}

	properties = append(properties, property)
	return properties, nil
}
