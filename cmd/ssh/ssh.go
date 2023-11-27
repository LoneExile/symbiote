package ssh

import (
	"log"
	"os"

	"symbiote/cmd"

	"golang.org/x/crypto/ssh"

	"golang.org/x/term"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sshx = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into an instance",
	Run: func(cmd *cobra.Command, args []string) {
		SSH()
	},
}

func init() {
	cmd.RootCmd.AddCommand(sshx)
}

func SSH() {

	var (
		sshServerHost = viper.GetString("SERVER.HOST")
		sshServerPort = viper.GetString("SERVER.SERVER_PORT")
		sshUser       = viper.GetString("SERVER.USER")
		sshPassword   = viper.GetString("SERVER.PASS")
	)

	// ----------------------------------------------------------------------------
	// https://stackoverflow.com/questions/28921409/how-can-i-send-terminal-escape-sequences-through-ssh-with-go

	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // NOTE: This is not recommended for production code
	}

	conn, err := ssh.Dial("tcp", sshServerHost+":"+sshServerPort, config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fileDescriptor := int(os.Stdin.Fd())

	if term.IsTerminal(fileDescriptor) {
		originalState, err := term.MakeRaw(fileDescriptor)
		if err != nil {
			panic(err)
		}
		defer term.Restore(fileDescriptor, originalState)
		if err != nil {
			panic(err)
		}
		termWidth, termHeight, err := term.GetSize(fileDescriptor)
		if err != nil {
			panic(err)
		}
		err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			panic(err)
		}
	}

	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	// Run commands from the terminal
	if err := session.Wait(); err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			log.Fatalf("Remote command exited with status: %d", exitErr.ExitStatus())
		} else {
			log.Fatalf("Failed to wait for session: %s", err)
		}
	}
}
