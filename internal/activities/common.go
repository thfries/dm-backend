package activities

type DittoClient struct {
	Host     string
	Username string
	Password string
}

type Activities struct {
	DittoClient *DittoClient
}
