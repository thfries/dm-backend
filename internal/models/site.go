// filepath: /Users/thomas/dev/dm-backend/internal/models/site.go
package models

type Site struct {
    SiteName    string `json:"siteName"`
    Host        string `json:"host"`
    Port        string `json:"port"`
    Username    string `json:"username"`
    Password    string `json:"password"`
    Description string `json:"description"`
}