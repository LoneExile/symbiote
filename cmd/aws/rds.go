package cmd

import (
	f "symbiote/cmd/aws/fn"

	"github.com/spf13/cobra"
)

var port string

var db = &cobra.Command{
	Use:   "rds",
	Short: "List instances",
	Run: func(cmd *cobra.Command, args []string) {
		f.RDSTunnel(port)
	},
}

func init() {
	db.Flags().StringVarP(&port, "port", "p", "5434:5432", "Port number local:remote")
	aws.AddCommand(db)
}
