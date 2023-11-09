package cmd

import (
	"symbiote/cmd"

	"github.com/spf13/cobra"
)

var sftp = &cobra.Command{

	Use:   "sftp",
	Short: "SFTP related commands",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cmd.RootCmd.AddCommand(sftp)
}
