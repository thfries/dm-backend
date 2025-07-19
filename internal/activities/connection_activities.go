package activities

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

//go:embed connection_template_mqtt5.json
var connectionTemplateMQTT5 string

// go:embed connection_template_another.json
// var connectionTemplateAnother string

var connectionTemplates = map[string]string{
	"mqtt5": connectionTemplateMQTT5,
	// "another": connectionTemplateAnother,
}

type CreateConnectionParams struct {
	ConnectionName string
	TemplateName   string            // e.g. "mqtt5"
	Placeholders   map[string]string // key: placeholder, value: replacement
}

func (c *DittoClient) CreateConnection(ctx context.Context, params CreateConnectionParams) (string, error) {
	tmplStr, ok := connectionTemplates[params.TemplateName]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", params.TemplateName)
	}
	tmpl, err := template.New("connection").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params.Placeholders); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Create Ditto connection via HTTP API
	url := fmt.Sprintf("http://%s/api/2/connections", c.Host)
	req, err := http.NewRequest("POST", url, &buf)
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
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("ditto API returned status %d: %s", resp.StatusCode, string(body))
	}

	var newConnection struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(body, &newConnection); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}
	return newConnection.ID, nil
}

// Register this activity with your Activities struct:
func (a *Activities) CreateConnection(ctx context.Context, params CreateConnectionParams) (string, error) {
	return a.DittoClient.CreateConnection(ctx, params)
}
