package sqs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func sendMsg(ctx context.Context, sess *session.Session, queueUrl string, payload string) (string, error) {
	svc := sqs.New(sess)

	sendMessageOutput, sendMessageErr := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: &payload,
		QueueUrl:    &queueUrl,
	})
	if sendMessageErr != nil {
		return "", sendMessageErr
	}
	return *sendMessageOutput.MessageId, nil
}

func PublishMessage(ctx context.Context, message string, awsRegion string, awsAccountId string, queueSuffix string) (string, error) {
	queueUrl := getQueueUrl(ctx, awsRegion, awsAccountId, queueSuffix)
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	return sendMsg(ctx, sess, queueUrl, message)
}

func getQueueUrl(ctx context.Context, awsRegion string, awsAccountId string, queueSuffix string) string {
	return fmt.Sprintf("https://sqs.%v.amazonaws.com/%v/%v", awsRegion, awsAccountId, queueSuffix)
}
