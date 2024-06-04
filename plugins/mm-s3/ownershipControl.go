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

type ownershipControlProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketOwnershipControlsOutput
}

func (oc *ownershipControlProp) fetchConf() error {
	output, err := oc.client.GetBucketOwnershipControls(
		context.Background(),
		&s3.GetBucketOwnershipControlsInput{
			Bucket: oc.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "OwnershipControlsNotFoundError":
				return &mmS3Error{ownershipControl, noConfig}
			default:
				return fmt.Errorf("fetchConf ownershipControl: %w", err)
			}
		}
		return fmt.Errorf("fetchConf ownershipControl: %w", err)
	}

	oc.configurations = output
	return nil
}

func (oc *ownershipControlProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := oc.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate ownershipControlProp: %w", err)
	}

	if oc.configurations.OwnershipControls != nil {
		property := shared.MinerProperty{
			Type: ownershipControl,
			Label: shared.MinerPropertyLabel{
				Name:   "OwnershipControls",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(oc.configurations.OwnershipControls); err != nil {
			return nil, fmt.Errorf("generate ownershipControlProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
