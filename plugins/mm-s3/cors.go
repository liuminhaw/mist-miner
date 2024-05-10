package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/liuminhaw/mist-miner/shared"
)

// getEncryptionProperties returns the encryption properties of a bucket
func getCorsProperties(client *s3.Client, bucket *types.Bucket) ([]shared.MinerProperty, error) {
	output, err := client.GetBucketCors(
		context.Background(),
		&s3.GetBucketCorsInput{
			Bucket: bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchCORSConfiguration":
				return nil, &mmS3Error{cors, noConfig}
			default:
				return nil, fmt.Errorf("getCorsProperties: %w", err)
			}
		}
		return nil, fmt.Errorf("getCorsProperties: %w", err)
	}

	var properties []shared.MinerProperty
	for _, rule := range output.CORSRules {
		sortCorsRule(&rule)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(rule); err != nil {
			return nil, fmt.Errorf("getCorsProperties: marshal CORS rule: %w", err)
		}

		corsValue := buffer.Bytes()
		h := md5.New()
		h.Write(corsValue)

		properties = append(properties, shared.MinerProperty{
			Type: cors,
			Label: shared.MinerPropertyLabel{
				Name:   fmt.Sprintf("%x", h.Sum(nil)),
				Unique: false,
			},
			Content: shared.MinerPropertyContent{
				Format: "json",
				Value:  string(corsValue),
			},
		})
	}

	return properties, nil
}

func sortCorsRule(rule *types.CORSRule) {
	sort.Strings(rule.AllowedMethods)
	sort.Strings(rule.AllowedOrigins)
	if rule.AllowedHeaders != nil {
		sort.Strings(rule.AllowedHeaders)
	}
	if rule.ExposeHeaders != nil {
		sort.Strings(rule.ExposeHeaders)
	}
}
