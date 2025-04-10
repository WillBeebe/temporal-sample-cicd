package cicdwf

import "time"

const TaskQueueName = "ci-cd-wf"

// START: activity inputs
type InfraDetails struct {
	DeployDir string
	Chart     string
}

type AppDetails struct {
	DatabaseConnectionString string
	Chart                    string
	Image                    string
}

// END: activity inputs

type WorkflowResult struct {
	TestsPassed         bool
	ApplicationDeployed bool
	Message             string
}

// START: activity outputs
type InfraResult struct {
	ClusterId                string
	DatabaseConnectionString string
	Message                  string
}

type AppResult struct {
	TestsPassed bool
	AppImage    string
	AppUrl      string
	Message     string
}

// END: activity outputs

type PollDetails struct {
	PollInterval time.Time
	Endpoint     string
}
