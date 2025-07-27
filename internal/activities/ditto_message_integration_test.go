package activities

import (
	"dm-backend/internal/models"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestSendDittoProtocolMessage_Integration(t *testing.T) {
	if os.Getenv("DITTO_INTEGRATION") != "1" {
		t.Skip("set DITTO_INTEGRATION=1 to run integration tests")
	}

	activitiesImpl := &Activities{
		DittoClient: &DittoClient{
			Host:     dittoHost,
			Username: dittoUsername,
			Password: dittoPassword,
		},
	}

	// 1. Create Thing using the existing method
	thingData := map[string]interface{}{
		"attributes": map[string]interface{}{
			"serial": "ws-test-serial",
			"foo":    "bar",
		},
	}
	createParams := CreateThingParams{
		Namespace:          dittoNamespace,
		UniqueAttributeKey: "serial",
		ThingData:          thingData,
	}
	t.Logf("Creating thing with params: %+v", createParams)
	thingID, err := activitiesImpl.DittoClient.CreateThing(createParams)
	if err != nil {
		t.Fatalf("CreateThing failed: %v", err)
	}
	t.Logf("Created thingID: %s", thingID)

	// 2. Create Feature
	featureMsg := models.DittoProtocolMessage{
		Topic: "<namespace>/<name>/things/twin/commands/modify",
		Path:  "/features/MyFeature",
		Value: map[string]interface{}{
			"properties": map[string]interface{}{
				"status": "active",
			},
		},
	}
	featureParams := SendDittoProtocolMessageParams{Message: featureMsg, ThingId: thingID}
	t.Logf("Sending create feature message: %+v, %s", featureMsg, thingID)
	if err := activitiesImpl.SendDittoProtocolMessage(t.Context(), featureParams); err != nil {
		t.Fatalf("failed to send create feature message: %v", err)
	}

	// 3. Modify Feature
	modifyMsg := models.DittoProtocolMessage{
		Topic: "<namespace>/<name>/things/twin/commands/modify",
		Path:  "/features/MyFeature/properties/status",
		Value: "inactive",
	}
	modifyParams := SendDittoProtocolMessageParams{Message: modifyMsg, ThingId: thingID}
	t.Logf("Sending modify feature message: %+v, %s", modifyMsg, thingID)
	if err := activitiesImpl.SendDittoProtocolMessage(t.Context(), modifyParams); err != nil {
		t.Fatalf("failed to send modify feature message: %v", err)
	}

	// 4. Delete Thing
	deleteMsg := models.DittoProtocolMessage{
		Topic: "<namespace>/<name>/things/twin/commands/delete",
		Path:  "/",
	}
	deleteParams := SendDittoProtocolMessageParams{Message: deleteMsg, ThingId: thingID}
	t.Logf("Sending delete thing message: %+v, %s", deleteMsg, thingID)
	if err := activitiesImpl.SendDittoProtocolMessage(t.Context(), deleteParams); err != nil {
		t.Fatalf("failed to send delete thing message: %v", err)
	}

	// Verify thing is deleted
	url := fmt.Sprintf("http://%s/api/2/things/%s", dittoHost, thingID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.SetBasicAuth(dittoUsername, dittoPassword)
	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("failed to GET thing after deletion: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 Not Found after deletion, got status: %d", resp.StatusCode)
	}
}
