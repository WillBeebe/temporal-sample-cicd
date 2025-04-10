package cicdwf

import (
	"context"
	"time"
)

type PollingActivities struct {
	TestService  *TestService
	PollInterval time.Duration
}

func (a *PollingActivities) DoPollActivity(ctx context.Context, input PollDetails) (string, error) {
	return a.TestService.GetServiceResult(ctx, input.Endpoint)
}
