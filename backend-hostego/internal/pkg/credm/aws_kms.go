package credm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var client *awsCredClient

type awsCredClient struct {
	kmsClient *kms.Client
	kmsKeyID  string
}

func NewAWSCredClient(ctx context.Context, awsRegion, kmsKeyID string) (*awsCredClient, error) {
	if client == nil {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
		if err != nil {
			return nil, err
		}
		kmsClient := kms.NewFromConfig(cfg)
		//For local, use this
		//kmsClient := kms.NewFromConfig(cfg, func(o *kms.Options) {
		//	o.EndpointResolver = kms.EndpointResolverFromURL("http://localhost:4566")
		//})

		client = &awsCredClient{
			kmsClient: kmsClient,
			kmsKeyID:  kmsKeyID,
		}
	}
	return client, nil
}

func (s *awsCredClient) EncryptText(ctx context.Context, secret []byte) ([]byte, error) {
	encryptOutput, err := s.kmsClient.Encrypt(ctx, &kms.EncryptInput{
		KeyId:     aws.String(s.kmsKeyID),
		Plaintext: secret,
	})
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting secret: %w", err)
	}
	return encryptOutput.CiphertextBlob, nil
}

func (s *awsCredClient) DecryptText(ctx context.Context, encryptedSecret []byte) ([]byte, error) {
	decryptOutput, err := s.kmsClient.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: encryptedSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("Error while encrypting secret: %w", err)
	}
	return decryptOutput.Plaintext, nil
}
