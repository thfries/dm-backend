package activities

import (
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type DittoClient struct {
	Host           string
	Username       string
	Password       string
	DevopsUsername string // Optional, if needed for devops operations
	DevopsPassword string // Optional, if needed for devops operations

	wsConn  *websocket.Conn
	wsMutex sync.Mutex
}

// ConnectWebSocket establishes a websocket connection to Ditto
func (c *DittoClient) ConnectWebSocket(wsURL string) error {
	c.wsMutex.Lock()
	defer c.wsMutex.Unlock()
	if c.wsConn != nil {
		return nil // already connected
	}
	dialer := websocket.DefaultDialer

	header := http.Header{}

	header.Set("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	conn, _, err := dialer.Dial(wsURL, header)
	if err != nil {
		return err
	}
	c.wsConn = conn
	return nil
}

// basicAuth returns the base64 encoded basic auth string for username and password
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return encodeBase64(auth)
}

// encodeBase64 encodes a string to base64
func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// CloseWebSocket closes the websocket connection
func (c *DittoClient) CloseWebSocket() error {
	c.wsMutex.Lock()
	defer c.wsMutex.Unlock()
	if c.wsConn != nil {
		err := c.wsConn.Close()
		c.wsConn = nil
		return err
	}
	return nil
}

type Activities struct {
	DittoClient *DittoClient
}
