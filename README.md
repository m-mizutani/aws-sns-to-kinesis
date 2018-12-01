# aws-sns-to-kinesis [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/aws-sns-to-kinesis)](https://goreportcard.com/report/github.com/m-mizutani/aws-sns-to-kinesis)

This is Serverless Application to forward a message from Amazon SNS (Simple Notification Service) to Amazon Kinesis Data Stream.

## Prerequisite

- Go >= 1.11
- aws-cli >= 1.14.40
- automake >= 3.81

## Usage

Create a config file named `stack.cfg`

```conf
StackName=your-sns-to-kinesis-stack
CodeS3Bucket=some-s3-bucket
CodeS3Prefix=functions

SnsTopicArn=arn:aws:sns:ap-northeast-1:1234567890:your-topic
KinesisStreamArn=arn:aws:kinesis:ap-northeast-1:1234567890:stream/your-stream
```

Parameter explanations are following.

- `StackName`: Stack name of CloudFormation
- `CodeS3Bucket`: S3 bucket name to save Lambda code
- `CodeS3Prefix`: S3 prefix to save Lambda code. NOTE: `/` is appended to tail of the parameter automatically
- `SnsTopicArn`: ARN of source SNS topic to receive messages
- `KinesisStreamArn`: ARN of destination Kinesis Stream
- `LambdaRoleArn` (optional): IAM Role ARN of Lambda function

After editing the config file, run deploy command.

```bash
$ make STACK_CONFIG=stack.cfg deploy
```

## Test

After deployment, you can run test.

```
env STACK_CONFIG=stack.cfg go test
```
