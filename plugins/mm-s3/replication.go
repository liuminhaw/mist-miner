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

type replicationProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketReplicationOutput
}

func (rp *replicationProp) fetchConf() error {
	output, err := rp.client.GetBucketReplication(
		context.Background(),
		&s3.GetBucketReplicationInput{
			Bucket: rp.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "ReplicationConfigurationNotFoundError":
				return &mmS3Error{replication, noConfig}
			default:
				return fmt.Errorf("fetchConf replicationProp: %w", err)
			}
		}
		return fmt.Errorf("fetchConf replicationProp: %w", err)
	}

	rp.configurations = output
	return nil
}

func (rp *replicationProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := rp.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate replicationProp: %w", err)
	}

	if rp.configurations.ReplicationConfiguration != nil {
		property := shared.MinerProperty{
			Type: replication,
			Label: shared.MinerPropertyLabel{
				Name:   "Replication",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
		if err := property.FormatContentValue(rp.configurations.ReplicationConfiguration); err != nil {
			return nil, fmt.Errorf("generate replicationProp: %w", err)
		}

		properties = append(properties, property)
	}

	return properties, nil
}
