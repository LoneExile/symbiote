package cmd

import (
	"symbiote/cmd"

	svc "symbiote/aws"

	"github.com/spf13/cobra"
)

var aws = &cobra.Command{
	Use:   "aws",
	Short: "AWS related commands",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	aws.PersistentFlags().StringVarP(&svc.Profile, "profile", "P", "", "AWS Profile")
	aws.PersistentFlags().StringVarP(&svc.Region, "region", "R", "", "AWS Region")

	cmd.RootCmd.AddCommand(aws)
}
