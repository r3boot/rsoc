package cliutils

import (
	"fmt"
	"testing"

	"strings"

	"github.com/r3boot/rsoc/lib/logger"
)

func TestExpandTilde(t *testing.T) {
	testLogger := logger.NewLogger(false, true)
	Setup(testLogger)

	testEmptyPath := ""
	result, err := ExpandTilde(testEmptyPath)
	if result != testEmptyPath {
		t.Errorf("ExpandTilde empty path: returned path != offered path")
	}
	if err != nil {
		t.Errorf("ExpandTilde empty path: err != nil: %v", err)
	}

	testAbsPath := "/nonexisting/path/to/somewhere.yaml"
	result, err = ExpandTilde(testAbsPath)
	if result != testAbsPath {
		fmt.Printf("result: %v", result)
		t.Errorf("ExpandTilde abs path: returned path != offered path")
	}
	if err != nil {
		t.Errorf("ExpandTilde abs path: err != nil: %v", err)
	}

	testTildePath := "~/.config/rsoc/config.yaml"
	result, err = ExpandTilde(testTildePath)
	if !strings.HasPrefix(result, "/") {
		t.Errorf("ExpandTilde tilde path: returned path does not start with /")
	}
	if err != nil {
		t.Errorf("ExpandTilde tilde path: err != nil: %v", err)
	}
}
