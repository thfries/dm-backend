package workflow

import (
	"dm-backend/internal/models"
	"time"

	"go.temporal.io/sdk/workflow"
)

type ConfigWorkflowParams struct {
	RQLQuery     string
	ConfigParams map[string]string
}

func MassDeviceConfigWorkflow(ctx workflow.Context, params ConfigWorkflowParams) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Fetch devices from Ditto
	var devices []models.Device
	// Use an activity that actually fetches devices, e.g., FetchDevicesActivity
	if err := workflow.ExecuteActivity(ctx, "FetchDevicesFromDitto", params.RQLQuery).Get(ctx, &devices); err != nil {
		return err
	}

	// Parallel config of devices
	futures := make([]workflow.Future, len(devices))
	for i, device := range devices {
		futures[i] = workflow.ExecuteActivity(ctx, "ConfigureDevice", device, params.ConfigParams)
	}
	for _, f := range futures {
		if err := f.Get(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}
