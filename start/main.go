package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	"cicdwf"
)

func main() {
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := cicdwf.InfraDetails{
		DeployDir: "/tmp/run-8546/ci-cd-infra",
		Chart:     "./my-app-chart",
	}

	options := client.StartWorkflowOptions{
		ID:        "integration-test-" + uuid.New(),
		TaskQueue: cicdwf.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, cicdwf.CICDWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result cicdwf.WorkflowResult

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Printf("Workflow result: %+v", result)
}
