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

type policyProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketPolicyOutput
}

func (p *policyProp) fetchConf() error {
	output, err := p.client.GetBucketPolicy(
		context.Background(),
		&s3.GetBucketPolicyInput{
			Bucket: p.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchBucketPolicy":
				return &mmS3Error{policy, noConfig}
			default:
				return fmt.Errorf("fetchConf policyProp: %w", err)
			}
		}
		return fmt.Errorf("fetchConf policyProp: %w", err)
	}

	p.configurations = output
	return nil
}

func (p *policyProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := p.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate policyProp: %w", err)
	}

	if p.configurations.Policy != nil {
		normalizedPolicy, err := shared.JsonNormalize(*p.configurations.Policy)
		if err != nil {
			return nil, fmt.Errorf("generate policyProp: %w", err)
		}

		property := shared.MinerProperty{
			Type: policy,
			Label: shared.MinerPropertyLabel{
				Name:   "Policy",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
				Value:  string(normalizedPolicy),
			},
		}

		properties = append(properties, property)
	}

	return properties, nil
}
