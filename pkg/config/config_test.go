package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {

	// Create temporary YAML config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "openqueue.yml")

	yamlContent := `
queues:
  - name: queue1
    db: test1
  - name: queue2
    db: test2
`

	err := os.WriteFile(configFile, []byte(yamlContent), 0644)
	require.NoError(t, err)

	ctx := context.Background()

	err = loadConfigFile(ctx, configFile)
	cfg := GetConfig()

	t.Log(cfg)
	require.NoError(t, err)
	assert.Equal(t, len(cfg.Queues), 2)
}
