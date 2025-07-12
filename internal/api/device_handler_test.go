package api

import (
    "bytes"
    "encoding/json"
    "io"
    "mass-device-config/internal/workflow"
    "mass-device-config/pkg/models"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.temporal.io/sdk/client"
)

// Mocks Temporal client for testing
type mockTemporalClient struct {
    client.Client
    startedWorkflow bool
    workflowID      string
    runID           string
    describeCalled  bool
}

func (m *mockTemporalClient) ExecuteWorkflow(
    ctx interface{},
    options client.StartWorkflowOptions,
    workflowFunc interface{},
    args ...interface{},
) (client.WorkflowRun, error) {
    m.startedWorkflow = true
    m.workflowID = "test-workflow-id"
    m.runID = "test-run-id"
    return &mockWorkflowRun{id: m.workflowID, runID: m.runID}, nil
}

func (m *mockTemporalClient) DescribeWorkflowExecution(
    ctx interface{}, workflowID, runID string,
) (*client.DescribeWorkflowExecutionResponse, error) {
    m.describeCalled = true
    return &client.DescribeWorkflowExecutionResponse{
        WorkflowExecutionInfo: &client.WorkflowExecutionInfo{
            Execution: &client.WorkflowExecution{
                ID:    workflowID,
                RunID: runID,
            },
            Status: client.WorkflowExecutionStatusCompleted,
        },
    }, nil
}

type mockWorkflowRun struct {
    id    string
    runID string
}

func (m *mockWorkflowRun) GetID() string    { return m.id }
func (m *mockWorkflowRun) GetRunID() string { return m.runID }
func (m *mockWorkflowRun) Get(ctx interface{}, valuePtr interface{}) error {
    return nil
}

func TestStartMassDeviceConfigHandler(t *testing.T) {
    mockClient := &mockTemporalClient{}

    handler := StartMassDeviceConfigHandler(mockClient)

    reqBody := StartConfigRequest{
        Devices: []models.Device{
            {ID: "dev1", Name: "Device 1"},
        },
        ConfigParams: map[string]string{"mode": "test"},
    }
    body, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/config/start", bytes.NewReader(body))
    w := httptest.NewRecorder()

    handler(w, req)

    res := w.Result()
    assert.Equal(t, http.StatusOK, res.StatusCode)

    responseBody, _ := io.ReadAll(res.Body)
    var resp map[string]string
    err := json.Unmarshal(responseBody, &resp)
    assert.NoError(t, err)
    assert.Equal(t, "test-workflow-id", resp["workflowID"])
    assert.Equal(t, "test-run-id", resp["runID"])
    assert.True(t, mockClient.startedWorkflow)
}

func TestStartMassDeviceConfigHandler_BadRequest(t *testing.T) {
    mockClient := &mockTemporalClient{}

    handler := StartMassDeviceConfigHandler(mockClient)

    req := httptest.NewRequest("POST", "/api/config/start", bytes.NewReader([]byte("bad json")))
    w := httptest.NewRecorder()

    handler(w, req)

    res := w.Result()
    assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetWorkflowStatusHandler(t *testing.T) {
    mockClient := &mockTemporalClient{}

    handler := GetWorkflowStatusHandler(mockClient)

    req := httptest.NewRequest("GET", "/api/config/status?workflowID=test-id&runID=test-run", nil)
    w := httptest.NewRecorder()

    handler(w, req)

    res := w.Result()
    assert.Equal(t, http.StatusOK, res.StatusCode)

    responseBody, _ := io.ReadAll(res.Body)
    var resp client.DescribeWorkflowExecutionResponse
    err := json.Unmarshal(responseBody, &resp)
    assert.NoError(t, err)
    assert.Equal(t, "test-id", resp.WorkflowExecutionInfo.Execution.ID)
    assert.Equal(t, "test-run", resp.WorkflowExecutionInfo.Execution.RunID)
    assert.True(t, mockClient.describeCalled)
}

func TestGetWorkflowStatusHandler_BadRequest(t *testing.T) {
    mockClient := &mockTemporalClient{}

    handler := GetWorkflowStatusHandler(mockClient)

    req := httptest.NewRequest("GET", "/api/config/status", nil)
    w := httptest.NewRecorder()

    handler(w, req)

    res := w.Result()
    assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}