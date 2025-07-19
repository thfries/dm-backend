package activities

import (
	"bytes"
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *DittoClient) ConfigureDevice(ctx context.Context, device models.Device, configParams map[string]string) error {
	url := fmt.Sprintf("http://%s/api/2/things/%s/attributes", c.Host, device.ThingId)
	payload, _ := json.Marshal(configParams)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/merge-patch+json")

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ditto HTTP call failed: %s", resp.Status)
	}
	return nil
}

// Wrap Activities.ConfigureDevice to use DittoClient
func (a *Activities) ConfigureDevice(ctx context.Context, device models.Device, configParams map[string]string) error {
	return a.DittoClient.ConfigureDevice(ctx, device, configParams)
}
