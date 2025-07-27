package activities

import (
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
)

// FetchDevicesFromDitto
func (c *DittoClient) FetchDevicesFromDitto(rqlQuery string) ([]models.Device, error) {
	url := fmt.Sprintf("http://%s/api/2/search/things?filter=%s", c.Host, rqlQuery)
	resp, respBody, err := c.doDittoRequest(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(respBody))
	}
	var res struct {
		Items []models.Device `json:"items"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quoted JSON string: %w", err)
	}
	return res.Items, nil
}

func (a *Activities) FetchDevicesFromDitto(ctx context.Context, rqlQuery string) ([]models.Device, error) {
	return a.DittoClient.FetchDevicesFromDitto(rqlQuery)
}
