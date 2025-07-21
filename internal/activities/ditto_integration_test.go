package activities

import (
	"bytes"
	"dm-backend/internal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	appConfig      = config.LoadConfig()
	dittoHost      = appConfig.DittoHost
	dittoUsername  = appConfig.DittoUsername
	dittoPassword  = appConfig.DittoPassword
	dittoNamespace = appConfig.DittoNamespace
)

func setupDittoTestData(t *testing.T) {
	thing := map[string]interface{}{
		"thingId":    "org.example:test-thing",
		"attributes": map[string]interface{}{"foo": "bar"},
	}
	payload, _ := json.Marshal(thing)
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/api/2/things/%s", dittoHost, "org.example:test-thing"),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetBasicAuth(dittoUsername, dittoPassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to PUT test thing: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected status from Ditto: %s", resp.Status)
	}
}

func cleanupDittoThing(t *testing.T, thingID string) {
	url := fmt.Sprintf("http://%s/api/2/things/%s", dittoHost, thingID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}
	req.SetBasicAuth(dittoUsername, dittoPassword)
	client := &http.Client{Timeout: 5 * time.Second}

	t.Logf("Cleanup - Deleting thing thingID: %s", thingID)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE thing: %v", err)
	}
	defer resp.Body.Close()
	// Accept 2xx and 404 (already deleted)
}

func TestCreateThing_Integration(t *testing.T) {
	if os.Getenv("DITTO_INTEGRATION") != "1" {
		t.Skip("set DITTO_INTEGRATION=1 to run integration tests")
	}

	client := &DittoClient{
		Host:     dittoHost,
		Username: dittoUsername,
		Password: dittoPassword,
	}

	namespace := dittoNamespace
	serialValue := fmt.Sprintf("test-%d", time.Now().UnixNano())
	thingData := map[string]interface{}{
		"attributes": map[string]interface{}{
			"serial": serialValue,
			"foo":    "bar", // second attribute
		},
	}
	params := CreateThingParams{
		Namespace:          namespace,
		UniqueAttributeKey: "serial",
		ThingData:          thingData,
	}

	// Test successful creation
	t.Logf("Creating thing with params: %+v", params)
	thingID, err := client.CreateThing(params)
	if err != nil {
		t.Fatalf("CreateThing failed: %v", err)
	}
	t.Logf("Created thingID: %s", thingID)
	defer cleanupDittoThing(t, thingID)

	// Check that the created thingID starts with the namespace
	if got, want := thingID, params.Namespace+":"; len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("thingID %q does not start with namespace %q", thingID, params.Namespace)
	}

	// Test duplicate creation is avoided
	t.Logf("Try to create duplicate thing with params: %+v", params)
	_, err = client.CreateThing(params)
	if err == nil {
		t.Error("expected error when creating duplicate thing, got nil")
	} else {
		t.Logf("duplicate creation correctly failed: %v", err)
	}
}

func TestFetchDevicesFromDitto_Integration(t *testing.T) {
	if os.Getenv("DITTO_INTEGRATION") != "1" {
		t.Skip("set DITTO_INTEGRATION=1 to run integration tests")
	}

	setupDittoTestData(t)

	client := &DittoClient{
		Host:     dittoHost,
		Username: dittoUsername,
		Password: dittoPassword,
	}

	rql := `eq(thingId,"org.example:test-thing")`
	devices, err := client.FetchDevicesFromDitto(rql)
	if err != nil {
		t.Fatalf("FetchDevicesFromDitto failed: %v", err)
	}
	if len(devices) == 0 {
		t.Fatal("expected at least one device, got none")
	}
	found := false
	for _, d := range devices {
		if d.ThingId == "org.example:test-thing" {
			found = true
			break
		}
	}
	if !found {
		t.Error("test device not found in Ditto response")
	}
}
