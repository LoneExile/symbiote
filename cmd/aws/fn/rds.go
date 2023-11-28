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

func RDSTunnel(port string) {
	r, err := svc.NewRDSClient()
	if err != nil {
		fmt.Println(err)
	}

	e, err := svc.NewEC2Client()
	if err != nil {
		fmt.Println(err)
	}

	instance := e.DefaultInstance()
	dbEndpoint := r.ListDBInstances()[0].PrivateIP
	p := strings.Split(port, ":")
	localPortNumber := p[0]
	portNumber := p[1]

	fmt.Println("DB Endpoint:", dbEndpoint)
	fmt.Println("Instance ID:", instance.InstanceID)
	fmt.Println("Instance Name:", instance.InstanceName)
	fmt.Println("Port (local:remote):", port)

	params := fmt.Sprintf(
		"{\"portNumber\":[\"%s\"],\"localPortNumber\":[\"%s\"],\"host\":[\"%s\"]}",
		portNumber,
		localPortNumber,
		dbEndpoint,
	)
	ssmCmd := exec.Command(
		"aws", "ssm", "start-session",
		"--target", instance.InstanceID,
		"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters", params,
		"--profile", svc.Profile,
	)

	ptmx, err := pty.Start(ssmCmd)
	if err != nil {
		fmt.Printf("Error starting command with pty: %v\n", err)
		return
	}
	fmt.Printf("Process started with PID: %d\n", ssmCmd.Process.Pid)
	defer func() { _ = ptmx.Close() }()
	defer func() { _ = ssmCmd.Process.Kill() }()

	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	scanner := bufio.NewScanner(ptmx)
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Waiting") {
			fmt.Printf("Waiting for connections on port %s.\n", portNumber)
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to start session on port %s.\n", portNumber)
	}

	if !found {
		fmt.Printf("Failed to start session on port %s.\n", portNumber)
	}

	if err := ssmCmd.Wait(); err != nil {
		fmt.Printf("Command returned error: %v\n", err)
	}
}
