package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type requestPaymentProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketRequestPaymentOutput
}

func (rp *requestPaymentProp) fetchConf() error {
	output, err := rp.client.GetBucketRequestPayment(
		context.Background(),
		&s3.GetBucketRequestPaymentInput{
			Bucket: rp.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf requestPayment: %w", err)
	}

	rp.configurations = output
	return nil
}

func (rp *requestPaymentProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := rp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate requestPaymentProp: %w", err)
	}

	property := shared.MinerProperty{
		Type: requestPayment,
		Label: shared.MinerPropertyLabel{
			Name:   "RequestPayment",
			Unique: true,
		},
		Content: shared.MinerPropertyContent{
			Format: shared.FormatText,
		},
	}
	if err := property.FormatContentValue(rp.configurations.Payer); err != nil {
		return nil, fmt.Errorf("generate requestPaymentProp: %w", err)
	}

	properties = append(properties, property)
	return properties, nil
}
