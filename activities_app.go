package cicdwf

import (
	"context"
	"fmt"
	"log"
)

func DeployApplication(ctx context.Context, details AppDetails) (AppResult, error) {
	commands := []string{
		"helm install myapp " + details.Chart + " --set dbConnectionString=" + details.DatabaseConnectionString + " --set image=" + details.Image,
	}

	for _, cmd := range commands {
		output, err := runCommand(ctx, cmd, ".")
		if err != nil {
			log.Printf("Command '%s' failed: %v\nOutput: %s", cmd, err, output)
			return AppResult{
				Message: fmt.Sprintf("command '%s' failed: %v\nOutput: %s", cmd, err, output),
			}, err
		}
	}

	// get the "live" app url...

	return AppResult{
		// AppUrl: "http://localhost:8080/myapp",
		AppUrl: "https://www.bing.com",
	}, nil
}

func RunApplicationTests(ctx context.Context, details AppDetails) (AppResult, error) {
	testsPassed := true

	commands := []string{
		"application tests command...",
	}

	for _, cmd := range commands {
		output, err := runCommand(ctx, cmd, ".")
		if err != nil {
			log.Printf("Command '%s' failed: %v\nOutput: %s", cmd, err, output)
			return AppResult{
				Message:     fmt.Sprintf("command '%s' failed: %v\nOutput: %s", cmd, err, output),
				TestsPassed: false,
			}, err
		}
	}

	return AppResult{
		TestsPassed: testsPassed,
		Message:     fmt.Sprintf("Tests %v", testsPassed),
	}, nil

}

func DestroyApplication(ctx context.Context, details AppDetails) (AppResult, error) {
	commands := []string{
		"helm uninstall myapp",
	}

	for _, cmd := range commands {
		output, err := runCommand(ctx, cmd, ".")
		if err != nil {
			log.Printf("Command '%s' failed: %v\nOutput: %s", cmd, err, output)
			return AppResult{
				Message: fmt.Sprintf("command '%s' failed: %v\nOutput: %s", cmd, err, output),
			}, err
		}
	}

	return AppResult{}, nil
}
