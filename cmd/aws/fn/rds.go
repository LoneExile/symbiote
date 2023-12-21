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

func RdsCmd(port string, profile string) *exec.Cmd {
	svc.Profile = profile
	r, err := svc.NewRDSClient()
	if err != nil {
		fmt.Println(err)
	}
	if len(r.ListDBInstances()) == 0 {
		fmt.Println("No RDS instances found.")
		return nil
	} else if len(r.ListDBInstances()) > 1 {
		fmt.Println("More than one RDS instance found.")
		return nil
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

	params := fmt.Sprintf(
		"{\"portNumber\":[\"%s\"],\"localPortNumber\":[\"%s\"],\"host\":[\"%s\"]}",
		portNumber,
		localPortNumber,
		dbEndpoint,
	)

	fmt.Println("DB Endpoint:", dbEndpoint)
	fmt.Println("Instance ID:", instance.InstanceID)
	fmt.Println("Instance Name:", instance.InstanceName)
	fmt.Println("Port (local:remote):", port)

	ssmCmd := exec.Command(
		"aws", "ssm", "start-session",
		"--target", instance.InstanceID,
		"--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
		"--parameters", params,
		"--profile", svc.Profile,
	)
	return ssmCmd
}

func RDSTunnel(port string, profile string) {

	ssmCmd := RdsCmd(port, profile)
	if ssmCmd == nil {
		return
	}

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
			fmt.Printf("Waiting for connections on port %s.\n", port)
			found = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to start session on port %s.\n", port)
	}

	if !found {
		fmt.Printf("Failed to start session on port %s.\n", port)
	}

	if err := ssmCmd.Wait(); err != nil {
		fmt.Printf("Command returned error: %v\n", err)
	}
}
