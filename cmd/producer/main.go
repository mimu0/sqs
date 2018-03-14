package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	mimusqs "github.com/mimu0/sqs"
)

func main() {
	var count = flag.Int("count", 1, "Message count to send")
	var useBatch = flag.Bool("batch", false, "Use batch to send")
	flag.Parse()

	batchType := map[bool]string{
		true:  "with",
		false: "without",
	}
	fmt.Printf("Send %d messages, %s batch\n", *count, batchType[*useBatch])

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sender := mimusqs.Sender{
		Client: sqs.New(sess),
		QURL:   "https://sqs.us-west-2.amazonaws.com/508300286574/mimutest",
		Count:  *count,
		Batch:  *useBatch,
	}
	if err := sender.Send(); err != nil {
		panic(err)
	}
}
