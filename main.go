package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sns2kinesis "github.com/m-mizutani/aws-sns-to-kinesis/lib"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
		log.WithField("log_level", os.Getenv("LOG_LEVEL")).Info("no properly log level setting, set InfoLevel as default")
	}

	lambda.Start(func(ctx context.Context, ev events.SNSEvent) (sns2kinesis.Result, error) {
		pusherEvent := sns2kinesis.Argument{ev, os.Getenv("DST_STREAM")}
		return sns2kinesis.Handler(pusherEvent)
	})
}
