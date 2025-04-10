package cicdwf

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CICDWorkflow(ctx workflow.Context, input InfraDetails) (*WorkflowResult, error) {
	logger := workflow.GetLogger(ctx)

	var wfresult = &WorkflowResult{
		TestsPassed: false,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			MaximumInterval:    time.Minute,
			BackoffCoefficient: 2,
			MaximumAttempts:    2,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	// Activity: build image
	var buildResult AppResult
	err := workflow.ExecuteActivity(ctx, BuildImage, input).Get(ctx, &buildResult)
	if err != nil {
		return wfresult, fmt.Errorf("failed to build image: %s", err)
	}

	// Activity: push image
	var pushResult AppResult
	err = workflow.ExecuteActivity(ctx, PushImage, input).Get(ctx, &pushResult)
	if err != nil {
		return wfresult, fmt.Errorf("failed to push image: %s", err)
	}

	// Activity: deploy infrastructure
	var infraResult InfraResult
	err = workflow.ExecuteActivity(ctx, DeployInfrastructure, input).Get(ctx, &infraResult)
	if err != nil {
		// Activity: take down test environment
		var resultDestroy InfraResult
		err = workflow.ExecuteActivity(ctx, DestroyInfrastructure, input).Get(ctx, &resultDestroy)
		if err != nil {
			return wfresult, fmt.Errorf("failed to destroy infrastructure: %s", err)
		}
		return wfresult, fmt.Errorf("failed to deploy infrastructure: %s", err)
	}

	// Activity: deploy application
	appInput := AppDetails{
		DatabaseConnectionString: infraResult.DatabaseConnectionString,
		Chart:                    input.Chart,
		Image:                    pushResult.AppImage,
	}
	var appResult AppResult
	err = workflow.ExecuteActivity(ctx, DeployApplication, appInput).Get(ctx, &appResult)
	if err != nil {
		// Activity: destroy application resources
		var destroyResult AppResult
		err = workflow.ExecuteActivity(ctx, DestroyApplication, input).Get(ctx, &destroyResult)
		if err != nil {
			return wfresult, fmt.Errorf("failed to destroy application: %s", err)
		}
		return wfresult, fmt.Errorf("failed to deploy application: %s", err)
	}

	// Workflow: poll application healthy
	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "poll-application-healthy-" + workflow.GetInfo(ctx).WorkflowExecution.ID,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	pollInput := PollDetails{
		Endpoint: appResult.AppUrl,
	}
	var result string
	err = workflow.ExecuteChildWorkflow(ctx, PollApplicationHealthy, pollInput).Get(ctx, &result)
	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)

		// Activity: destroy application resources
		var destroyResult AppResult
		err = workflow.ExecuteActivity(ctx, DestroyApplication, input).Get(ctx, &destroyResult)
		if err != nil {
			return wfresult, fmt.Errorf("failed to destroy application: %s", err)
		}

		// Activity: take down test environment
		var resultDestroy InfraResult
		err = workflow.ExecuteActivity(ctx, DestroyInfrastructure, input).Get(ctx, &resultDestroy)
		if err != nil {
			return wfresult, fmt.Errorf("failed to destroy infrastructure: %s", err)
		}

		return wfresult, err
	}
	wfresult.ApplicationDeployed = true

	// Activity: Run Application Tests
	var testsResult AppResult
	err = workflow.ExecuteActivity(ctx, RunApplicationTests, input).Get(ctx, &testsResult)
	if err != nil {
		logger.Error("failed to run application tests.", "Error", err)
		// return wfresult, fmt.Errorf("failed to run application tests: %s", err)
	}
	wfresult.TestsPassed = testsResult.TestsPassed

	// Activity: destroy application resources
	var destroyResult AppResult
	err = workflow.ExecuteActivity(ctx, DestroyApplication, input).Get(ctx, &destroyResult)
	if err != nil {
		logger.Error("failed to destroy application.", "Error", err)
		// return wfresult, fmt.Errorf("failed to destroy application: %s", err)
	}

	// Activity: take down test environment
	var resultDestroy InfraResult
	err = workflow.ExecuteActivity(ctx, DestroyInfrastructure, input).Get(ctx, &resultDestroy)
	if err != nil {
		return wfresult, fmt.Errorf("failed to destroy infrastructure: %s", err)
	}

	return wfresult, nil
}
