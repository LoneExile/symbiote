package aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/spf13/viper"
)

type Ec2InstanceInfo struct {
	InstanceID   string
	InstanceName string
	Status       string
	PrivateIP    string
	State        string
}

type AWSClients struct {
	EC2 *ec2.Client
	RDS *rds.Client
}

// TODO: add default, check config file => ENV => ~/.aws/config
// - check if region is valid
var (
	Region  string
	Profile string
)

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

	fmt.Println("Profile:", Profile)
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

func NewEC2Client() (*AWSClients, error) {
	cfg := cfg()
	ec2Client := ec2.NewFromConfig(cfg)

	return &AWSClients{
		EC2: ec2Client,
	}, nil
}

func (c *AWSClients) ListInstances() []Ec2InstanceInfo {
	input := &ec2.DescribeInstancesInput{}
	result, err := c.EC2.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instances, %v", err)
	}

	var instances []Ec2InstanceInfo
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PrivateIpAddress == nil {
				continue
			}

			instances = append(instances, Ec2InstanceInfo{
				InstanceID:   *instance.InstanceId,
				InstanceName: selectTagVal(instance.Tags, "Name"),
				Status:       string(instance.State.Name),
				PrivateIP:    string(*instance.PrivateIpAddress),
				State:        string(instance.State.Name),
			})

		}
	}
	return instances
}

func (c *AWSClients) ListRegions() {
	regionsOutput, err := c.EC2.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		log.Fatalf("failed to describe regions, %v", err)
	}
	fmt.Println("Available AWS Regions:")
	for _, region := range regionsOutput.Regions {
		fmt.Println(*region.RegionName)
	}

}

func (c *AWSClients) GetEndPoint() *ec2.DescribeInstanceConnectEndpointsOutput {
	result, err := c.EC2.DescribeInstanceConnectEndpoints(
		context.TODO(),
		&ec2.DescribeInstanceConnectEndpointsInput{},
	)
	if err != nil {
		log.Fatalf("failed to describe VPC endpoint services, %v", err)
	}

	return result
}

func (c *AWSClients) DefaultEndpoint() string {
	eice := c.GetEndPoint()

	instanceCECount := len(eice.InstanceConnectEndpoints)
	eiceID := ""

	if instanceCECount == 0 {
		fmt.Println("No VPC Endpoint Service Found")
		os.Exit(1)
		return eiceID
	}

	if instanceCECount > 1 {
		fmt.Println("VPC Endpoint Service Count:", instanceCECount)
		fmt.Println("More than one VPC Endpoint Service Found, using the first one")
	}

	eiceID = *eice.InstanceConnectEndpoints[0].InstanceConnectEndpointId
	fmt.Printf("VPC Endpoint Service Names: %s\n", eiceID)

	return eiceID
}

func (c *AWSClients) DefaultInstance() Ec2InstanceInfo {
	instances := c.ListInstances()
	instancesName := viper.GetString("GLOBAL.AWS.INSTANCE_TAG_NAME")
	if instancesName == "" {
		fmt.Println("Please set the INSTANCE_TAG_NAME in config.yaml")
		os.Exit(1)
	}

	var instance Ec2InstanceInfo
	for _, val := range instances {
		if strings.Contains(val.InstanceName, instancesName) {
			instance = val
			break
		}
	}
	return instance
}
