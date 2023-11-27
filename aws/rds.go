package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func NewRDSClient() (*AWSClients, error) {
	cfg := cfg()
	rdsClient := rds.NewFromConfig(cfg)

	return &AWSClients{
		RDS: rdsClient,
	}, nil
}

type RDSInstanceInfo struct {
	InstanceID   string
	InstanceName string
	Status       string
	PrivateIP    string
	State        string
}

func (c *AWSClients) ListDBInstances() []RDSInstanceInfo {
	input := &rds.DescribeDBInstancesInput{}
	result, err := c.RDS.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instances, %v", err)
	}

	var instances []RDSInstanceInfo
	for _, db := range result.DBInstances {
		if db.Endpoint == nil {
			continue
		}

		instances = append(instances, RDSInstanceInfo{
			InstanceID:   *db.DBInstanceIdentifier,
			InstanceName: *db.DBInstanceIdentifier,
			Status:       *db.DBInstanceStatus,
			PrivateIP:    *db.Endpoint.Address,
			State:        *db.DBInstanceStatus,
		})

	}

	return instances

}
