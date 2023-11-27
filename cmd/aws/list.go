package cmd

import (
	f "symbiote/cmd/aws/fn"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List instances",
	Run: func(cmd *cobra.Command, args []string) {
		f.ListInstances()
	},
}

func init() {
	aws.AddCommand(listCmd)
}
