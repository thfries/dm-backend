package activities

import (
	"context"
	"dm-backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DittoClient struct {
	Host string
}

func (c *DittoClient) FetchDevicesFromDitto(rqlQuery string) ([]models.Device, error) {
	url := fmt.Sprintf("https://%s/api/2/search/things?filter=%s", c.Host, rqlQuery)
	req, _ := http.NewRequest("GET", url, nil)
	// Add authentication headers if needed
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Things []models.Device `json:"things"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Things, nil
}

type Activities struct {
	Ditto *DittoClient
}

func (a *Activities) FetchDevicesFromDitto(ctx context.Context, rqlQuery string) ([]models.Device, error) {
	return a.Ditto.FetchDevicesFromDitto(rqlQuery)
}

// Similarly, wrap ConfigureDevice if it needs DittoClient
