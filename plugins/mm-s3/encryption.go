package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type encryptionConfig struct {
	sse       string
	kmsId     string
	bucketKey string
}

// Value: ServerSideEncryption|KMSMasterKeyID|BucketKeyEnabled
func getEncryptionProperties(
	client *s3.Client,
	bucket *types.Bucket,
) ([]shared.MinerProperty, error) {
	output, err := client.GetBucketEncryption(context.Background(), &s3.GetBucketEncryptionInput{
		Bucket: bucket.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("getEncryptionProperties: %w", err)
	}

	var properties []shared.MinerProperty
	for _, rule := range output.ServerSideEncryptionConfiguration.Rules {
		config := encryptionConfig{}
		// Check nil
		if rule.BucketKeyEnabled == nil {
			log.Printf("BucketKeyEnabled is nil")
			config.bucketKey = ""
		} else if *rule.BucketKeyEnabled {
			config.bucketKey = "true"
		} else {
			config.bucketKey = "false"
		}

		if rule.ApplyServerSideEncryptionByDefault == nil {
			log.Printf("ApplyServerSideEncryptionByDefault is nil")
		} else if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == types.ServerSideEncryptionAes256 {
			config.sse = string(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
		} else {
			config.sse = string(rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm)
			config.kmsId = *rule.ApplyServerSideEncryptionByDefault.KMSMasterKeyID
		}

		properties = append(properties, shared.MinerProperty{
			Type: encryption,
			Label: shared.MinerPropertyLabel{
				Name:   "Rule",
				Unique: false,
			},
			Content: shared.MinerPropertyContent{
				Format: "string",
				Value: strings.Join(
					[]string{config.sse, config.kmsId, config.bucketKey},
					valueSeparator,
				),
			},
		})

	}

	return properties, nil
}
