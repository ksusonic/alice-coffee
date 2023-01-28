package queue

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const REGION = "ru-central1"

type MessageQueue struct {
	QueueUrl string
	Client   *sqs.Client
}

func NewMessageQueue(queueUrl, accessKey, secretKey string) *MessageQueue {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               queueUrl,
			SigningRegion:     REGION,
			HostnameImmutable: true,
		}, nil
	})

	return &MessageQueue{
		QueueUrl: queueUrl,
		Client: sqs.NewFromConfig(aws.Config{
			Region:                      REGION,
			Credentials:                 NewCredentialsProvider(accessKey, secretKey),
			EndpointResolverWithOptions: resolver,
		}),
	}
}

func (m *MessageQueue) SendMessage(data []byte) (*string, error) {
	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(data)),
		QueueUrl:    &m.QueueUrl,
	}

	resp, err := m.Client.SendMessage(context.TODO(), sMInput)
	if err != nil {
		return nil, fmt.Errorf("Got an error sending the message: %v", err)
	}

	return resp.MessageId, nil
}
