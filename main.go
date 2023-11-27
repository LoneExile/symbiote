package main

import (
	"symbiote/cmd"
	_ "symbiote/cmd/aws"
	_ "symbiote/cmd/sftp"
	_ "symbiote/cmd/ssh"
)

func main() {
	cmd.Execute()
}
