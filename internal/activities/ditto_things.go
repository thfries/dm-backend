package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateThingParams struct {
	Namespace          string
	UniqueAttributeKey string
	ThingData          map[string]interface{} // full JSON for the new thing
}

// CreateThing
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
	searchResp, searchBody, err := c.doDittoRequest(context.Background(), "GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("search HTTP request failed: %w", err)
	}
	defer searchResp.Body.Close()
	if searchResp.StatusCode < 200 || searchResp.StatusCode >= 300 {
		return "", fmt.Errorf("ditto search API returned status %d: %s", searchResp.StatusCode, string(searchBody))
	}
	var searchResult struct {
		Items []struct {
			ThingID string `json:"thingId"`
		} `json:"items"`
	}
	if err := json.Unmarshal(searchBody, &searchResult); err != nil {
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
	resp, respBody, err := c.doDittoRequest(context.Background(), "POST", url, bodyBytes)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(respBody))
	}
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

type DeleteThingParams struct {
	ThingID string
}

// DeleteThing
func (c *DittoClient) DeleteThing(ctx context.Context, params DeleteThingParams) error {
	url := fmt.Sprintf("http://%s/api/2/things/%s", c.Host, params.ThingID)
	resp, respBody, err := c.doDittoRequest(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Accept 2xx and 404 (already deleted)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		// Thing already deleted, treat as success
		return nil
	}
	return fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(respBody))
}

func (a *Activities) DeleteThing(ctx context.Context, params DeleteThingParams) error {
	return a.DittoClient.DeleteThing(ctx, params)
}
