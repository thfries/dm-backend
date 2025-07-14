package activities

import (
	"bytes"
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func ConfigureDevice(ctx context.Context, device models.Device, configParams map[string]string) error {
	url := fmt.Sprintf("https://<ditto-host>/api/2/things/%s", device.ID)
	payload, _ := json.Marshal(map[string]interface{}{
		"attributes": configParams,
	})
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// Add authentication headers if needed
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ditto PUT failed: %s", resp.Status)
	}
	return nil
}
