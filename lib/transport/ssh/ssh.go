package ssh

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Run(args []string) (string, string, error) {
	var stdoutBuff, stderrBuff bytes.Buffer

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Start()
	if err != nil {
		return "", "", fmt.Errorf("Run cmd.Start: %v", err)
	}

	err = cmd.Wait()
	stdout := ""
	if stdoutBuff.Len() > 0 {
		stdout = stdoutBuff.String()
	}

	stderr := ""
	if stderrBuff.Len() > 0 {
		stderr = stderrBuff.String()
	}

	return stdout, stderr, nil
}
