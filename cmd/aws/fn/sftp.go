package fn

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	svc "symbiote/aws"

	"github.com/creack/pty"
)

func EicSFTPCmd() *exec.Cmd {
	c, err := svc.NewEC2Client()
	if err != nil {
		fmt.Println(err)
	}
	instance := c.DefaultInstance()
	eiceID := c.DefaultEndpoint()
	profileConfig := svc.GetProfileConfig()
	// fmt.Println("Private IP:", instance.PrivateIP)
	// fmt.Println("Port:", profileConfig.Port)

	tunnelCmd := exec.Command(
		"aws", "ec2-instance-connect", "open-tunnel",
		"--instance-connect-endpoint-id", eiceID,
		"--private-ip-address", instance.PrivateIP,
		"--local-port", profileConfig.Port,
		"--remote-port", "22",
		"--profile", profileConfig.Name,
	)
	return tunnelCmd
}

func SFTPConnectCmd() *exec.Cmd {
	profileConfig := svc.GetProfileConfig()
	sftpCmd := exec.Command(
		"sftp",
		"-P", profileConfig.Port,
		"-i", profileConfig.PemKeyPath,
		"ec2-user@localhost",
	)
	return sftpCmd
}

func SFTP(server bool) {
	// profileConfig := svc.GetProfileConfig()
	tunnelCmd := EicSFTPCmd()
	ptmx, err := pty.Start(tunnelCmd)
	if err != nil {
		fmt.Printf("Error starting command with pty: %v\n", err)
		return
	}
	fmt.Printf("Process started with PID: %d\n", tunnelCmd.Process.Pid)
	defer func() { _ = ptmx.Close() }()
	defer func() { _ = tunnelCmd.Process.Kill() }()

	// Copy the pty's output to the stdout (terminal)
	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	scanner := bufio.NewScanner(ptmx)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Listening") {
			// fmt.Printf("Listening for connections on port %s.\n", profileConfig.Port)
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from pty: %v\n", err)
	}

	if !found {
		fmt.Println("The word was not found in the output.")
	}

	if !server {
		fmt.Println("------------------")
		sftpCmd := SFTPConnectCmd()

		sftpCmd.Stdin = os.Stdin
		sftpCmd.Stdout = os.Stdout
		sftpCmd.Stderr = os.Stderr

		if err := sftpCmd.Run(); err != nil {
			fmt.Printf("Error running sftp command: %v\n", err)
		}
	} else {
		if err := tunnelCmd.Wait(); err != nil {
			fmt.Printf("Command returned error: %v\n", err)
		}
	}
}
