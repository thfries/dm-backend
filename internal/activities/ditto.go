package activities

import (
	"bytes"
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

type CreateThingParams struct {
	Namespace          string
	UniqueAttributeKey string
	ThingData          map[string]interface{} // full JSON for the new thing
}

func (c *DittoClient) CreateThing(params CreateThingParams) (string, error) {
	// 1. Check that UniqueAttributeKey is present in ThingData["attributes"]
	attrs, ok := params.ThingData["attributes"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("thing data must contain an 'attributes' object")
	}
	val, exists := attrs[params.UniqueAttributeKey]
	if !exists {
		return "", fmt.Errorf("thing data is missing unique attribute key '%s' in attributes", params.UniqueAttributeKey)
	}

	// 2. Search for existing thing with the unique attribute value
	rql := fmt.Sprintf(`eq(attributes/%s,"%v")`, params.UniqueAttributeKey, val)
	searchURL := fmt.Sprintf("http://%s/api/2/search/things?filter=%s", c.Host, rql)
	searchReq, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create search request: %w", err)
	}
	searchReq.SetBasicAuth(c.Username, c.Password)
	searchResp, err := http.DefaultClient.Do(searchReq)
	if err != nil {
		return "", fmt.Errorf("search HTTP request failed: %w", err)
	}
	defer searchResp.Body.Close()
	if searchResp.StatusCode < 200 || searchResp.StatusCode >= 300 {
		body, _ := io.ReadAll(searchResp.Body)
		return "", fmt.Errorf("ditto search API returned status %d: %s", searchResp.StatusCode, string(body))
	}
	var searchResult struct {
		Items []struct {
			ThingID string `json:"thingId"`
		} `json:"items"`
	}
	if err := json.NewDecoder(searchResp.Body).Decode(&searchResult); err != nil {
		return "", fmt.Errorf("failed to decode search response: %w", err)
	}
	if len(searchResult.Items) > 0 {
		return "", fmt.Errorf("thing with %s=%v already exists", params.UniqueAttributeKey, val)
	}

	// 3. Create the new thing with the provided data
	url := fmt.Sprintf("http://%s/api/2/things?namespace=%s&requested-acks=twin-persisted,search-persisted&timeout=10", c.Host, params.Namespace)
	bodyBytes, err := json.Marshal(params.ThingData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal thing body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the returned thingId from Ditto's response
	var result struct {
		TwinPersisted struct {
			Payload struct {
				ThingID string `json:"thingId"`
			} `json:"payload"`
		} `json:"twin-persisted"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return result.TwinPersisted.Payload.ThingID, nil
}

func (a *Activities) CreateThing(ctx context.Context, params CreateThingParams) (string, error) {
	return a.DittoClient.CreateThing(params)
}
