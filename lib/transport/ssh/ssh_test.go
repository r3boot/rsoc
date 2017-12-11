package ssh

import (
	"os"
	"strings"
	"testing"
)

const (
	TEST_UNAME = "Linux"
)

func TestRun(t *testing.T) {
	if os.Getenv("TRAVIS") == "" {
		args := []string{"localhost", "uname"}
		stdout, stderr, err := Run(args)
		if len(stdout) > 0 {
			if !strings.HasPrefix(stdout, TEST_UNAME) {
				t.Errorf("ssh.Run: uname called, expected %s, got %s", TEST_UNAME, stdout)
			}
		} else {
			t.Errorf("ssh.Run: uname called, but got no output")
		}

		if len(stderr) > 0 {
			t.Errorf("ssh.Run uname called, but got stderr output")
		}

		if err != nil {
			t.Errorf("ssh.Run uname called, err != nil: %v", err)
		}
	}

	args := []string{"--nonexistent", "--commandline", "--parameters"}
	stdout, stderr, err := Run(args)
	if len(stdout) > 0 {
		t.Errorf("ssh.Run: invalid ssh params, but got stdout output")
	}

	if len(stderr) == 0 {
		t.Errorf("ssh.Run: invalid ssh params, but got no stderr output")
	}

	if err != nil && !strings.Contains(err.Error(), "Run cmd.wait: exit status 255") {
		t.Errorf("ssh.Run invalid ssh params, err != nil but got invalid error: %v", err)
	}

	args = []string{"test"}
	stdout, stderr, err = Run(args)
	if len(stdout) > 0 {
		t.Errorf("ssh.Run: invalid exec.Command params, but got stdout output")
	}

	if len(stderr) > 0 {
		t.Errorf("ssh.Run: invalid exec.Command params, but got stderr output")
	}

	if err != nil && !strings.Contains(err.Error(), "cmd.Start: exec: already started") {
		t.Errorf("ssh.Run invalid exec.Command params, err != nil but got invalid error: %v", err)
	}
}
