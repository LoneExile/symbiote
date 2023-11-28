package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// TODO: add default, check config file => ENV => ~/.aws/config
// - check if region is valid
var (
	Region  string
	Profile string
)

type AWSClients struct {
	EC2 *ec2.Client
	RDS *rds.Client
	SSM *ssm.Client
}

func cfg() aws.Config {
	if Profile == "" {
		Profile = GetProfile()
	}
	if Region == "" {
		Region = GetRegion()
	}

	optFns := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(Profile),
	}

	// fmt.Println("Profile:", Profile)
	if Region != "" {
		fmt.Println("Region:", Region)
		optFns = append(optFns, config.WithRegion(Region))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		optFns...,
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	return cfg
}
