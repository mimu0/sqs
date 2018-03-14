package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Sender struct {
	Client *sqs.SQS
	QURL   string
	Count  int
	Batch  bool
}

func (s Sender) sendWithBatch() error {
	generatedCount := 0
	for generatedCount < s.Count {
		var partitionSize int
		if partitionSize = s.Count - generatedCount; s.Count-generatedCount > 10 {
			partitionSize = 10
		}
		sendMessages := make([]*sqs.SendMessageBatchRequestEntry, partitionSize)
		for p := 0; p < partitionSize; p++ {
			sendMessages[p] = &sqs.SendMessageBatchRequestEntry{
				Id:          aws.String(fmt.Sprintf("%d", generatedCount)),
				MessageBody: aws.String(fmt.Sprintf("Message %d", generatedCount)),
			}
			generatedCount++
		}
		_, err := s.Client.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  sendMessages,
			QueueUrl: &s.QURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Sender) sendEach() error {
	for i := 0; i < s.Count; i++ {
		_, err := s.Client.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(fmt.Sprintf("Message %d", i)),
			QueueUrl:    &s.QURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Sender) Send() error {
	if s.Batch {
		return s.sendWithBatch()
	}
	return s.sendEach()
}
