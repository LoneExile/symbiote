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
		ListFiles()
	},
}

func init() {
	sftp.AddCommand(list)
}

func ListFiles() {
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
}
