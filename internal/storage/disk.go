//go:build save_disk

package storage

import (
	// standard
	"bytes"
	"os"

	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/logger"
)

func storageInit() {}

func storageSave(output string, data bytes.Buffer) {
	os.WriteFile(output, data.Bytes(), 0644)
	Log.Info("file created", "path", output)
}
