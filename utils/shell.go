package utils

import (
	"os/exec"
)

func ExecuteCLI(head string, args ...string) ([]byte, error) {

	cmd := exec.Command(head, args...)
	outputBytes, err := cmd.Output()

	if err != nil {
		return outputBytes, err
	}

	return outputBytes, nil
}
