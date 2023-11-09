package cmd

import (
	"fmt"

	f "symbiote/sftp"

	"github.com/spf13/cobra"
)

var list = &cobra.Command{

	Use:   "list",
	Short: "List all files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n list command called")
		fmt.Println("-------------------")

		client, err := f.InitSFTP()
		if err != nil {
			fmt.Println(err)
		}

		err = f.ListFiles(client.Sc)
		if err != nil {
			fmt.Println(err)
		}

		defer client.Conn.Close()
		defer client.Sc.Close()

		fmt.Println("-------------------")
	},
}

func init() {
	sftp.AddCommand(list)
}
