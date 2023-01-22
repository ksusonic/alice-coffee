package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const REGION = "ru-central1"
const YaCloudApiUrl = "https://message-queue.api.cloud.yandex.net"

type MessageQueue struct {
	QueueUrl string
	Client   *sqs.Client
}

func NewMessageQueue(queueUrl string, config aws.Config) *MessageQueue {
	return &MessageQueue{
		QueueUrl: queueUrl,
		Client:   sqs.NewFromConfig(config),
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
