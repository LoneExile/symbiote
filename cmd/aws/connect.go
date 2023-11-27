package cmd

import (
	f "symbiote/cmd/aws/fn"

	"github.com/spf13/cobra"
)

// aws configure list-profiles
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an instance",
	Run: func(cmd *cobra.Command, args []string) {
		f.ConnectToInstance()
	},
}

func init() {
	aws.AddCommand(connectCmd)
}
