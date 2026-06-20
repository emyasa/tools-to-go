package loader

import "testing"

func TestLoadConfigFromGit(t *testing.T) {
	LoadFromGit(
		"/Users/emyasa/.ssh/id_ed25519",
		"git@github.com:emyasa/scratch-config",
		"main",
		"localhost:6379",
		"general",
	)
}
