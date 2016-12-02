package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	expectedAddress := "1.1.1.1"
	os.Setenv("KUBE_SOURCE_ADDRESS", expectedAddress)

	Load()

	actualAddress := Get("address")

	if actualAddress != expectedAddress {
		t.Errorf("Config address failed: expected %s, got %s\n", expectedAddress,  actualAddress)
	}
}