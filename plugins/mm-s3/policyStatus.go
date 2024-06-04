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

type policyStatusProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketPolicyStatusOutput
}

func (ps *policyStatusProp) fetchConf() error {
	output, err := ps.client.GetBucketPolicyStatus(
		context.Background(),
		&s3.GetBucketPolicyStatusInput{
			Bucket: ps.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchBucketPolicy":
				return &mmS3Error{policyStatus, noConfig}
			default:
				return fmt.Errorf("fetchConf policyStatusProp: %w", err)
			}
		}
		return fmt.Errorf("fetchConf policyStatusProp: %w", err)
	}

	ps.configurations = output
	return nil
}

func (ps *policyStatusProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := ps.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate policyStatusProp: %w", err)
	}

	if ps.configurations.PolicyStatus != nil {
		property := shared.MinerProperty{
			Type: policyStatus,
			Label: shared.MinerPropertyLabel{
				Name:   "PolicyStatus",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(ps.configurations.PolicyStatus); err != nil {
			return nil, fmt.Errorf("generate policyStatusProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
