package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/liuminhaw/mist-miner/shared"
)

func getLifecycleProperties(
	client *s3.Client,
	bucket *types.Bucket,
) ([]shared.MinerProperty, error) {
	output, err := client.GetBucketLifecycleConfiguration(
		context.Background(),
		&s3.GetBucketLifecycleConfigurationInput{
			Bucket: bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchLifecycleConfiguration":
				return nil, &mmS3Error{lifecycle, noConfig}
			default:
				return nil, fmt.Errorf("getLifeCycleProperties: %w", err)
			}
		}
		return nil, fmt.Errorf("getLifeCycleProperties: %w", err)
	}

	var properties []shared.MinerProperty
	for _, rule := range output.Rules {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(rule); err != nil {
			return nil, fmt.Errorf(
				"getLifeCycleProperties: marshal LifeCycle rule: %w",
				err,
			)
		}
		ruleValue := buffer.Bytes()

		properties = append(properties, shared.MinerProperty{
			Type: lifecycle,
			Label: shared.MinerPropertyLabel{
				Name:   *rule.ID,
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: "json",
				Value:  string(ruleValue),
			},
		})
	}

	return properties, nil
}
