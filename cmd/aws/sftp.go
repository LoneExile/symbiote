package cmd

import (
	f "symbiote/cmd/aws/fn"

	"github.com/spf13/cobra"
)

var server bool

var sftp = &cobra.Command{
	Use:   "sftp",
	Short: "SFTP to an instance",
	Run: func(cmd *cobra.Command, args []string) {
		f.SFTP()
	},
}

func init() {
	sftp.Flags().BoolVarP(&server, "server", "s", false, "Run sftp server mode")
	aws.AddCommand(sftp)
}
