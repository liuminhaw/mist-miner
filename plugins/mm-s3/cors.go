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

type corsProp struct {
    client         *s3.Client
    bucket         *types.Bucket
    configurations *s3.GetBucketCorsOutput
}

func (cp *corsProp) fetchConf() error {
	output, err := cp.client.GetBucketCors(
		context.Background(),
		&s3.GetBucketCorsInput{
			Bucket: cp.bucket.Name,
		},
	)
	if err != nil {
		var apiErr smithy.APIError
		if ok := errors.As(err, &apiErr); ok {
			switch apiErr.ErrorCode() {
			case "NoSuchCORSConfiguration":
				return &mmS3Error{cors, noConfig}
			default:
				return fmt.Errorf("fetchConf corsProp: %w", err)
			}
		}
		return fmt.Errorf("fetchConf corsProp: %w", err)
	}

    cp.configurations = output
    return nil
}

func (cp *corsProp) generate() ([]shared.MinerProperty, error) {
	var properties []shared.MinerProperty

    if err := cp.fetchConf(); err != nil {
        return nil, fmt.Errorf("generate corsProp: %w", err)
    }
	for _, rule := range cp.configurations.CORSRules {
		sortCorsRule(&rule)
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(rule); err != nil {
			return nil, fmt.Errorf("generate corsProp: marshal rule: %w", err)
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
