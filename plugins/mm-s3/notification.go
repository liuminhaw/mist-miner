package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type notificationProp struct {
	client         *s3.Client
	bucket         *types.Bucket
	configurations *s3.GetBucketNotificationConfigurationOutput
}

func (np *notificationProp) fetchConf() error {
	output, err := np.client.GetBucketNotificationConfiguration(
		context.Background(),
		&s3.GetBucketNotificationConfigurationInput{
			Bucket: np.bucket.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("fetchConf notification: %w", err)
	}

	np.configurations = output
	return nil
}

func (np *notificationProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

	if err := np.fetchConf(); err != nil {
		return nil, fmt.Errorf("generate notificationProp: %w", err)
	}

	if np.notificationIsEmpty() {
		log.Println("No notification configuration found")
	} else {
		property := shared.MinerProperty{
			Type: notification,
			Label: shared.MinerPropertyLabel{
				Name:   "Notification",
				Unique: true,
			},
			Content: shared.MinerPropertyContent{
				Format: shared.FormatJson,
			},
		}
        if err := property.FormatContentValue(np.configurations); err != nil {
            return nil, fmt.Errorf("generate notificationProp: %w", err)
        }

		properties = append(properties, property)
	}

	return properties, nil
}

func (np *notificationProp) notificationIsEmpty() bool {
	return np.configurations.EventBridgeConfiguration == nil &&
		len(np.configurations.LambdaFunctionConfigurations) == 0 &&
		len(np.configurations.QueueConfigurations) == 0 &&
		len(np.configurations.TopicConfigurations) == 0
}
