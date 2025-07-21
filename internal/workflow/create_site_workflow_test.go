package workflow_test

import (
	"errors"
	"testing"

	"dm-backend/internal/activities"
	"dm-backend/internal/workflow"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

type MockActivities struct {
	FailCreateThing      bool
	FailCreateConnection bool
}

func (m *MockActivities) CreateThing(_ interface{}, _ activities.CreateThingParams) (string, error) {
	if m.FailCreateThing {
		return "", errors.New("CreateThing failed")
	}
	return "gateway:site-thing-id", nil
}

func (m *MockActivities) CreateConnection(_ interface{}, _ activities.CreateConnectionParams) (string, error) {
	if m.FailCreateConnection {
		return "", errors.New("CreateConnection failed")
	}
	return "site-conn-id", nil
}

func (m *MockActivities) DeleteThing(_ interface{}, _ string) error {
	return nil
}

func TestCreateSiteWorkflow_Success(t *testing.T) {
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	mockActs := &MockActivities{}
	env.RegisterActivity(mockActs.CreateThing)
	env.RegisterActivity(mockActs.CreateConnection)
	env.RegisterActivity(mockActs.DeleteThing)

	params := workflow.CreateSiteParams{
		SiteName:    "site1",
		Host:        "localhost",
		Port:        "1883",
		Username:    "user",
		Password:    "pass",
		Description: "Test site",
	}

	env.ExecuteWorkflow(workflow.CreateSiteWorkflow, params)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func TestCreateSiteWorkflow_CreateThingFails(t *testing.T) {
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	mockActs := &MockActivities{FailCreateThing: true}
	env.RegisterActivity(mockActs.CreateThing)
	env.RegisterActivity(mockActs.CreateConnection)
	env.RegisterActivity(mockActs.DeleteThing)

	params := workflow.CreateSiteParams{
		SiteName:    "site1",
		Host:        "localhost",
		Port:        "1883",
		Username:    "user",
		Password:    "pass",
		Description: "Test site",
	}

	env.ExecuteWorkflow(workflow.CreateSiteWorkflow, params)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
	require.Contains(t, env.GetWorkflowError().Error(), "CreateThing failed")
}

func TestCreateSiteWorkflow_CreateConnectionFailsWithCompensation(t *testing.T) {
	ts := testsuite.WorkflowTestSuite{}
	env := ts.NewTestWorkflowEnvironment()

	mockActs := &MockActivities{FailCreateConnection: true}
	env.RegisterActivity(mockActs.CreateThing)
	env.RegisterActivity(mockActs.CreateConnection)
	env.RegisterActivity(mockActs.DeleteThing)

	params := workflow.CreateSiteParams{
		SiteName:    "site1",
		Host:        "localhost",
		Port:        "1883",
		Username:    "user",
		Password:    "pass",
		Description: "Test site",
	}

	env.ExecuteWorkflow(workflow.CreateSiteWorkflow, params)
	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
	require.Contains(t, env.GetWorkflowError().Error(), "createGatewayConnection failed after retries")
}
