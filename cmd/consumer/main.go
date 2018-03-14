package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-redis/redis"
	mimusqs "github.com/mimu0/sqs"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	consumer := mimusqs.Consumer{
		Client:        sqs.New(sess),
		QURL:          "https://sqs.us-west-2.amazonaws.com/508300286574/mimutest",
		DeleteMessage: make(chan string),
		Redis: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
	consumer.Start()
}
