package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	svc "symbiote/aws"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an instance",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := svc.NewAWSClients()
		if err != nil {
			fmt.Println(err)
		}

		instance := c.DefaultInstance()

		execCmd := exec.Command(
			"aws",
			"ec2-instance-connect",
			"ssh",
			"--instance-id",
			instance.InstanceID,
			"--connection-type",
			"eice",
			"--profile",
			svc.Profile,
			"--color",
			"on",
		)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err = execCmd.Run()
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

	},
}

func init() {
	aws.AddCommand(connectCmd)
}
