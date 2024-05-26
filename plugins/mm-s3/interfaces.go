package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/liuminhaw/mist-miner/shared"
)

type crawler interface {
	fetchConf() error
	generate() ([]shared.MinerProperty, error)
}

type propConstructor func(client *s3.Client, bucket *types.Bucket) crawler

var propConstructors = map[string]propConstructor{
    accelerateConfig: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &accelerateProp{
            client: client,
            bucket: bucket,
        }
    },
    acl: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &aclProp{
            client: client,
            bucket: bucket,
        }
    },
    cors: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &corsProp{
            client: client,
            bucket: bucket,
        }
    },
    encryption: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &encryptionProp{
            client: client,
            bucket: bucket,
        }
    },
    intelligentTiering: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &intelligentTieringProp{
            client: client,
            bucket: bucket,
        }
    },
    inventory: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &inventoryProp{
            client: client,
            bucket: bucket,
        }
    },
    lifecycle: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &lifecycleProp{
            client: client,
            bucket: bucket,
        }
    },
    tagging: func(client *s3.Client, bucket *types.Bucket) crawler {
        return &taggingProp{
            client: client,
            bucket: bucket,
        }
    },
}

func New(client *s3.Client, bucket *types.Bucket, propType string) (crawler, error) {
    constructor, ok := propConstructors[propType]
    if !ok {
        return nil, fmt.Errorf("New crawler: unknown property type: %s", propType)
    }
    return constructor(client, bucket), nil
}
