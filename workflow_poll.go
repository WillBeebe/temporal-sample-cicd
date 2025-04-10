package cicdwf

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func PollApplicationHealthy(ctx workflow.Context, input PollDetails) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Child workflow PollApplicationHealthy")

	options := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    5 * time.Second,
			MaximumInterval:    time.Minute,
			BackoffCoefficient: 2,
			MaximumAttempts:    10,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	a := &PollingActivities{}

	var pollResult string
	err := workflow.ExecuteActivity(ctx, a.DoPollActivity, input).Get(ctx, &pollResult)
	if err != nil {
		return fmt.Sprintf("application did not come up healthy: %s", err), fmt.Errorf("application did not come up healthy: %s", err)
	}

	return "", nil
}
