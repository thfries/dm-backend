package activities

import (
    "testing"
    "mass-device-config/pkg/models"
    "github.com/stretchr/testify/require"
)

func TestConfigureDevice(t *testing.T) {
    device := models.Device{ID: "dev1", Name: "TestDevice"}
    err := ConfigureDevice(nil, device, map[string]string{"foo": "bar"})
    require.NoError(t, err)
}