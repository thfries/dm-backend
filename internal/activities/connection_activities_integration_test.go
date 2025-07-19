package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	testConnectionName = "test-mqtt5-connection"
	testTemplateName   = "mqtt5"
	testMQTTHost       = "localhost"
	testMQTTPort       = "1883"
)

func cleanupDittoConnection(t *testing.T, dittoHost, username, password, connectionName string) {

	// Fetch all connections from Ditto and find the connectionID by name
	url := fmt.Sprintf("http://%s/api/2/connections?fields=id,name", dittoHost)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request for connections: %v", err)
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to GET connections: %v", err)
	}
	defer resp.Body.Close()

	type connection struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	var result []connection

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode connections response: %v", err)
	}

	var connectionID string
	for _, c := range result {
		if c.Name == connectionName {
			connectionID = c.ID
			break
		}
	}
	if connectionID == "" {
		// If not found, nothing to clean up
		return
	}

	url = fmt.Sprintf("http://%s/api/2/connections/%s", dittoHost, connectionID)
	req, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}
	req.SetBasicAuth(username, password)
	client = &http.Client{Timeout: 5 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE connection: %v", err)
	}
	defer resp.Body.Close()
	// Accept 2xx and 404 (already deleted)
	if resp.StatusCode >= 300 && resp.StatusCode != 404 {
		t.Fatalf("unexpected status from Ditto DELETE: %s", resp.Status)
	}
}

func TestCreateConnection_Integration(t *testing.T) {
	if os.Getenv("DITTO_INTEGRATION") != "1" {
		t.Skip("set DITTO_INTEGRATION=1 to run integration tests")
	}

	dittoHost := testDittoHost
	username := testDittoDevopsUsername
	password := testDittoDevopsPassword

	// Clean up before and after test
	cleanupDittoConnection(t, dittoHost, username, password, testConnectionName)
	defer cleanupDittoConnection(t, dittoHost, username, password, testConnectionName)

	client := &DittoClient{
		Host:     dittoHost,
		Username: username,
		Password: password,
	}

	params := CreateConnectionParams{
		ConnectionName: testConnectionName,
		TemplateName:   testTemplateName,
		Placeholders: map[string]string{
			"ConnectionName": testConnectionName,
			"MQTTHost":       testMQTTHost,
			"MQTTPort":       testMQTTPort,
		},
	}

	newConnectionId, err := client.CreateConnection(context.Background(), params)
	if err != nil {
		t.Fatalf("CreateConnection failed: %v", err)
	}

	// Optionally, verify the connection exists
	url := fmt.Sprintf("http://%s/api/2/connections/%s", dittoHost, newConnectionId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.SetBasicAuth(username, password)
	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("failed to GET connection: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from Ditto GET, got: %s", resp.Status)
	}
}
