package sns2kinesis

import (
	"crypto/sha256"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	// "github.com/aws/aws-lambda-go/lambdacontext"
)

// Argument is for Handler
type Argument struct {
	Sns       events.SNSEvent
	StreamArn string
}

type ResultLog struct {
	Message string `json:"message"`
	Done    bool   `json:"done"`
	Error   string `json:"error"`
}
type Result struct {
	Logs []*ResultLog `json:"logs"`
}

func makeKinesisClient(arn string) (*kinesis.Kinesis, string) {
	arr := strings.Split(arn, ":")
	if len(arr) != 6 {
		log.WithField("invalid_arn", arn).Error("Invali ARN format")
	}
	nameParts := strings.Split(arr[5], "/")
	if len(nameParts) != 2 || nameParts[0] != "stream" {
		log.WithField("invalid_arn", arn).Error("Invali ARN format (must be 'stream/*'")
	}

	ssn := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(arr[3]),
	}))
	kinesisClient := kinesis.New(ssn)

	return kinesisClient, nameParts[1]
}

// Handler is a main function, should be called from test
func Handler(args Argument) (Result, error) {
	log.WithField("args", args).Info("Start")
	var result Result

	kinesisClient, streamName := makeKinesisClient(args.StreamArn)
	input := kinesis.PutRecordsInput{
		StreamName: aws.String(streamName),
	}

	for _, record := range args.Sns.Records {
		data := record.SNS.Message

		reslog := ResultLog{Message: data}
		result.Logs = append(result.Logs, &reslog)

		entry := kinesis.PutRecordsRequestEntry{
			Data:         []byte(data),
			PartitionKey: aws.String(fmt.Sprintf("%x", sha256.Sum256([]byte(data)))),
		}

		input.Records = append(input.Records, &entry)
		reslog.Done = true
	}

	_, err := kinesisClient.PutRecords(&input)

	if err != nil {
		log.WithField("input", input).Error("Kinesis PutRecords error")
	}

	log.WithField("result", result).Info("Exit normally")
	return result, nil
}
