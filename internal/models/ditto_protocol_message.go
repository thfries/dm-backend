package models

// DittoProtocolMessage represents a message in Ditto protocol format.
type DittoProtocolMessage struct {
	Topic   string                 `json:"topic"`
	Path    string                 `json:"path,omitempty"`
	Value   interface{}            `json:"value,omitempty"`
	Headers map[string]interface{} `json:"headers,omitempty"`
}
