package sqs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
)

type Consumer struct {
	Client        *sqs.SQS
	QURL          string
	DeleteMessage chan string
	Redis         *redis.Client
}

func (c Consumer) Start() {
	go c.deleteLoop()
	for {
		result, err := c.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &c.QURL,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(3600), // 1 hours
			WaitTimeSeconds:     aws.Int64(5),
		})

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		if len(result.Messages) == 0 {
			continue
		}

		body := *result.Messages[0].Body
		fmt.Println(body)
		if v, err := c.Redis.GetSet(body, true).Result(); v != "" {
			fmt.Printf("Conflict! %s, %v\n", body, err)
		}
		c.Redis.Expire(body, 2*time.Second)
		c.DeleteMessage <- *result.Messages[0].ReceiptHandle
	}
}

func (c Consumer) deleteLoop() {
	for {
		receiptHandle := <-c.DeleteMessage
		c.deleteOne(receiptHandle)
	}
}

func (c Consumer) deleteOne(receiptHandle string) error {
	_, err := c.Client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &c.QURL,
		ReceiptHandle: &receiptHandle,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
