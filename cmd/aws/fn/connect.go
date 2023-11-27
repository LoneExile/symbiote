package fn

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	svc "symbiote/aws"
)

func getInstances() svc.Ec2InstanceInfo {
	c, err := svc.NewEC2Client()
	if err != nil {
		fmt.Println(err)
	}
	instance := c.DefaultInstance()
	return instance
}

func ConnectCmd() *exec.Cmd {
	instanceID := getInstances().InstanceID
	execCmd := exec.Command(
		"aws",
		"ec2-instance-connect",
		"ssh",
		"--instance-id",
		instanceID,
		"--connection-type",
		"eice",
		"--profile",
		svc.Profile,
		"--color",
		"on",
	)
	return execCmd
}

func ConnectToInstance() {
	execCmd := ConnectCmd()

	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	err := execCmd.Run()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	defer func() { _ = execCmd.Process.Kill() }()
}
