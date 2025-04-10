package main

import (
	"cicdwf"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, cicdwf.TaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(cicdwf.CICDWorkflow)
	w.RegisterWorkflow(cicdwf.PollApplicationHealthy)
	w.RegisterActivity(cicdwf.BuildImage)
	w.RegisterActivity(cicdwf.PushImage)
	w.RegisterActivity(cicdwf.DeployInfrastructure)
	w.RegisterActivity(cicdwf.DestroyInfrastructure)
	w.RegisterActivity(cicdwf.DeployApplication)
	w.RegisterActivity(cicdwf.RunApplicationTests)
	w.RegisterActivity(cicdwf.DestroyApplication)
	// todo: use poll interval
	activities := &cicdwf.PollingActivities{
		TestService: cicdwf.NewTestService(),
		// todo: use poll interval
		PollInterval: 5 * time.Second,
	}
	w.RegisterActivity(activities)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
