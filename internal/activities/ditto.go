package activities

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *DittoClient) doDittoRequest(ctx context.Context, method, url string, body []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create %s request: %w", method, err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	return resp, respBody, nil
}
