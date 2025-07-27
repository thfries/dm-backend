package activities

import (
	"os"
	"testing"
)

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
