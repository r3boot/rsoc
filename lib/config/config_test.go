package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/r3boot/rsoc/lib/logger"
)

const (
	TMP_DIR    = "/tmp"
	TMP_PREFIX = "rsoc_lib_config_test"
)

func CreateTempFile(t *testing.T) string {
	tmpFile, err := ioutil.TempFile(TMP_DIR, TMP_PREFIX)
	if err != nil {
		t.Errorf("Failed to create tempfile: %v", err)
	}
	fname := tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		t.Errorf("Failed to close tempfile: %v", err)
	}
	if err := os.Remove(fname); err != nil {
		t.Errorf("Failed to remove tempfile: %v", err)
	}

	return fname
}

func CleanupTempFile(t *testing.T, fname string) {
	if err := os.Remove(fname); err != nil {
		t.Errorf("Failed to remove tempfile: %v", err)
	}
}

func TestNewConfig(t *testing.T) {
	tmpFile := CreateTempFile(t)
	defer CleanupTempFile(t, tmpFile)

	log := logger.NewLogger(false, true)

	config, err := NewConfig(log, tmpFile)
	if err != nil {
		t.Errorf("NewConfig err != nil: %v", err)
	}

	if config == nil {
		t.Errorf("NewConfig config == nil")
	}
}
