package workflow

import (
    "testing"
    "mass-device-config/pkg/models"
    "go.temporal.io/sdk/testsuite"
    "github.com/stretchr/testify/require"
)

func TestMassDeviceConfigWorkflow(t *testing.T) {
    var suite testsuite.WorkflowTestSuite
    env := suite.NewTestWorkflowEnvironment()

    devices := []models.Device{
        {ID: "dev1", Name: "Device 1"},
        {ID: "dev2", Name: "Device 2"},
    }
    params := ConfigWorkflowParams{
        Devices:      devices,
        ConfigParams: map[string]string{"param1": "value1"},
    }

    env.ExecuteWorkflow(MassDeviceConfigWorkflow, params)
    require.True(t, env.IsWorkflowCompleted())
    require.NoError(t, env.GetWorkflowError())
}