package sftp

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pkg/sftp"
)

type SFTPClient struct {
	Sc   *sftp.Client
	Conn *ssh.Client
}

func NewSFTPClient() *SFTPClient {
	return &SFTPClient{}
}

func InitSFTP() (*SFTPClient, error) {

	user := viper.GetString("SERVER.USER")
	pass := viper.GetString("SERVER.PASS")
	host := viper.GetString("SERVER.HOST")

	port := viper.GetInt("SERVER.PORT")

	var auths []ssh.AuthMethod

	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}

	if pass != "" {
		auths = append(auths, ssh.Password(pass))
	}

	config := ssh.ClientConfig{
		User: user,
		Auth: auths,
		// Uncomment to ignore host key check
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := ssh.Dial("tcp", addr, &config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	// defer conn.Close()

	sc, err := sftp.NewClient(conn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}

	sftpClient := NewSFTPClient()
	sftpClient.Sc = sc
	sftpClient.Conn = conn

	fmt.Fprintf(os.Stdout, "Connected to %s\n", host)

	return sftpClient, nil
}
