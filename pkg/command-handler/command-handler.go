package commandhandler

import (
	"os"
	"os/exec"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func MustGetCommandOutput(programName, errorMsg string, args ...string) string {
	cmd := exec.Command(programName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.WriteErrorf("%s: %s\n", errorMsg, err)
	}

	return string(output)
}

func MustGetCommandOutputEvenIfExitError(programName, errorMsg string, args ...string) string {
	cmd := exec.Command(programName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return string(output)
		}

		logger.WriteErrorf("%s: %s\n", errorMsg, err)
	}

	return string(output)
}

func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

func MustGetUserConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logger.WriteErrorf("failed to get user config directory: %s\n", err)
	}

	return configDir
}
