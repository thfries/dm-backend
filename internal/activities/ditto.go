package activities

import (
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *DittoClient) FetchDevicesFromDitto(rqlQuery string) ([]models.Device, error) {
	url := fmt.Sprintf("http://%s/api/2/search/things?filter=%s", c.Host, rqlQuery)
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(c.Username, c.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, _ := io.ReadAll(resp.Body)

	type Result struct {
		Items []models.Device `json:"items"`
	}
	var res Result

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal quoted JSON string: %w", err)
	}

	return res.Items, nil
}

func (a *Activities) FetchDevicesFromDitto(ctx context.Context, rqlQuery string) ([]models.Device, error) {
	return a.DittoClient.FetchDevicesFromDitto(rqlQuery)
}
