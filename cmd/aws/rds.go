package cmd

import (
	"fmt"
	svc "symbiote/aws"

	// f "symbiote/cmd/aws/fn"

	"github.com/spf13/cobra"
)

var db = &cobra.Command{
	Use:   "rds",
	Short: "List instances",
	Run: func(cmd *cobra.Command, args []string) {
		RDSTunnel()
	},
}

// aws ssm start-session --target i-08524dd18ebc04037
// --document-name AWS-StartPortForwardingSessionToRemoteHost
// --parameters '{"portNumber":["5432"],"localPortNumber":["5434"],"host":["demo-data-database-1.cus5tynd mnve.ap-southeast-1.rds.amazonaws.com"]}'
// --profile demo2

func RDSTunnel() {
	r, err := svc.NewRDSClient()
	if err != nil {
		fmt.Println(err)
	}

	e, err := svc.NewEC2Client()
	if err != nil {
		fmt.Println(err)
	}

	instance := e.DefaultInstance()
	fmt.Println("Private IP:", instance.PrivateIP)
	// profileConfig := svc.GetProfileConfig()
	list := r.ListDBInstances()

	for _, db := range list {
		fmt.Println(db.InstanceID)
	}

}

func init() {
	aws.AddCommand(db)
}
