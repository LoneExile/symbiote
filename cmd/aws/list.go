package cmd

import (
	"fmt"
	svc "symbiote/aws"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List instances",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := svc.NewAWSClients()
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
	},
}

func init() {
	aws.AddCommand(listCmd)
}
