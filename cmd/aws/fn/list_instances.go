package fn

import (
	"fmt"
	svc "symbiote/aws"
)

func ListInstances() {
	c, err := svc.NewEC2Client()
	if err != nil {
		fmt.Println(err)
	}
	instances := c.ListInstances()

	instancesLen := len(instances)
	if instancesLen == 0 {
		fmt.Println("\nNo instances found.")
		return
	}

	fmt.Println("EC2 Instances:")
	for key, val := range instances {
		fmt.Printf(
			"%d: %s (%s) %s %s\n",
			key,
			val.InstanceName,
			val.Status,
			val.PrivateIP,
			val.State,
		)
	}
}
