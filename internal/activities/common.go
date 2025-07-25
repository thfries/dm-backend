package activities

type DittoClient struct {
	Host           string
	Username       string
	Password       string
	DevopsUsername string // Optional, if needed for devops operations
	DevopsPassword string // Optional, if needed for devops operations
}

type Activities struct {
	DittoClient *DittoClient
}
