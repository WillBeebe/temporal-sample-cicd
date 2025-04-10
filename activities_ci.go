package cicdwf

import (
	"context"
	"fmt"
	"log"
)

func BuildImage(ctx context.Context, details AppDetails) (AppResult, error) {
	commands := []string{
		"docker build " + details.Image,
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

func PushImage(ctx context.Context, details AppDetails) (AppResult, error) {
	commands := []string{
		"docker push",
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
		AppImage: "us-docker.pkg.dev/my-project/my-repo/my-app:4433aa",
	}, nil

}
