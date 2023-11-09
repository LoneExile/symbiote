package main

import (
	"symbiote/cmd"
	_ "symbiote/cmd/aws"
	_ "symbiote/cmd/sftp"
)

func main() {
	cmd.Execute()
}
