package aws

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func NewSSMClient() (*AWSClients, error) {
	cfg := cfg()
	ssmClient := ssm.NewFromConfig(cfg)

	return &AWSClients{
		SSM: ssmClient,
	}, nil
}
