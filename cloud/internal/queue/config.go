package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type CredentialsProvider struct {
	accessKey, secretKey string
}

func NewCredentialsProvider(accessKey, secretKey string) *CredentialsProvider {
	return &CredentialsProvider{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

func (c CredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.accessKey,
		SecretAccessKey: c.secretKey,
	}, nil
}
