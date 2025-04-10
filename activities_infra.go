package cicdwf

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

func DeployInfrastructure(ctx context.Context, details InfraDetails) (InfraResult, error) {
	commands := []string{
		"terraform init",
		"terraform plan",
		"terraform apply -auto-approve",
	}

	for _, cmd := range commands {
		output, err := runCommand(ctx, cmd, details.DeployDir)
		if err != nil {
			activity.GetLogger(ctx).Info(fmt.Sprintf("Command '%s' failed: %v\nOutput: %s", cmd, err, output))
			return InfraResult{
				Message: fmt.Sprintf("command '%s' failed: %v\nOutput: %s", cmd, err, output),
			}, err
		}
	}

	// get values from terraform output...
	return InfraResult{
		ClusterId:                "gke_project_zone_cluster-123456",
		DatabaseConnectionString: "postgresql://user:password@localhost:5432/mydb",
		Message:                  fmt.Sprintf("Deployed %s", details.DeployDir),
	}, nil
}

func DestroyInfrastructure(ctx context.Context, details InfraDetails) (InfraResult, error) {
	commands := []string{
		"terraform init",
		"terraform plan",
		"terraform destroy -auto-approve",
	}

	for _, cmd := range commands {
		output, err := runCommand(ctx, cmd, details.DeployDir)
		if err != nil {
			return InfraResult{
				Message: fmt.Sprintf("command '%s' failed: %v\nOutput: %s", cmd, err, output),
			}, err
		}
	}

	return InfraResult{
		Message: fmt.Sprintf("Destroyed %s", details.DeployDir),
	}, nil
}

func runCommand(ctx context.Context, command, dir string) (string, error) {
	activity.GetLogger(ctx).Info(fmt.Sprintf("Running command: %s in directory: %s", command, dir))
	// cmd := exec.Command("sh", "-c", command)
	// cmd.Dir = dir

	// Simulate command execution time with random sleep
	time.Sleep(time.Duration(500+rand.Intn(2500)) * time.Millisecond)

	// output, err := cmd.CombinedOutput()
	output := "example output"
	// if err != nil {
	// 	log.Printf("Command failed: %v", err)
	// 	return "", err
	// }

	activity.GetLogger(ctx).Info("Command executed successfully")
	return output, nil
}
