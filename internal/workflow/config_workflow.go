package workflow

import (
	"dm-backend/internal/activities"
	"dm-backend/pkg/models"
	"time"

	"go.temporal.io/sdk/workflow"
)

type ConfigWorkflowParams struct {
	Devices      []models.Device
	ConfigParams map[string]string
}

func MassDeviceConfigWorkflow(ctx workflow.Context, params ConfigWorkflowParams) error {
	// Parallel config of devices
	futures := make([]workflow.Future, len(params.Devices))
	for i, device := range params.Devices {
		activityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute,
		})
		futures[i] = workflow.ExecuteActivity(activityCtx, activities.ConfigureDevice, device, params.ConfigParams)
	}
	// Wait for all activities
	for _, f := range futures {
		if err := f.Get(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}
