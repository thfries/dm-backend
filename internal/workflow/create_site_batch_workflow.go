package workflow

import (
	"dm-backend/internal/models"
	"time"

	"go.temporal.io/sdk/workflow"
)

// CreateSiteBatchWorkflowParams defines the input for the batch workflow
type CreateSiteBatchWorkflowParams struct {
	Sites []models.Site
}

// CreateSiteBatchWorkflow starts a CreateSiteWorkflow for each site in the batch
func CreateSiteBatchWorkflow(ctx workflow.Context, params CreateSiteBatchWorkflowParams) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	futures := make([]workflow.Future, len(params.Sites))
	for i, site := range params.Sites {
		siteParams := CreateSiteParams{Site: site}
		futures[i] = workflow.ExecuteChildWorkflow(ctx, CreateSiteWorkflow, siteParams)
	}

	// Wait for all child workflows to complete
	for _, f := range futures {
		if err := f.Get(ctx, nil); err != nil {
			return err
		}
	}
	return nil
}
