package activities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	testDittoHost           = "localhost:8080"
	testDittoUsername       = "ditto"
	testDittoPassword       = "ditto"
	testDittoDevopsUsername = "devops"
	testDittoDevopsPassword = "foobar"
)

func setupDittoTestData(t *testing.T) {
	thing := map[string]interface{}{
		"thingId":    "org.example:test-thing",
		"attributes": map[string]interface{}{"foo": "bar"},
	}
	payload, _ := json.Marshal(thing)
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("http://%s/api/2/things/%s", testDittoHost, "org.example:test-thing"),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetBasicAuth(testDittoUsername, testDittoPassword)
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

func TestFetchDevicesFromDitto_Integration(t *testing.T) {
	if os.Getenv("DITTO_INTEGRATION") != "1" {
		t.Skip("set DITTO_INTEGRATION=1 to run integration tests")
	}

	setupDittoTestData(t)

	client := &DittoClient{
		Host:     testDittoHost,
		Username: testDittoUsername,
		Password: testDittoPassword,
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
