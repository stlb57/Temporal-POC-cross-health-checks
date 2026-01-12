package main

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MonitorTwinWorkflow(
	ctx workflow.Context,
	peerURL string,
	selfID string,
	port2 int,
) error {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 3,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	
	workflow.Sleep(ctx, time.Second*30)

	for {
		var healthy bool
		err := workflow.ExecuteActivity(
			ctx,
			CheckHealth,
			peerURL,
		).Get(ctx, &healthy)

		if err != nil || !healthy {
			workflow.ExecuteActivity(
				ctx,
				RestartTwin,
				selfID,
				port2,
			).Get(ctx, nil)
		}

		workflow.Sleep(ctx, time.Second*5)
	}
}
