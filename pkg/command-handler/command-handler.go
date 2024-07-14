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

func MustRunCommand(programName, errorMsg string, args ...string) {
	cmd := exec.Command(programName, args...)
	err := cmd.Run()
	if err != nil {
		logger.WriteErrorf("%s: %s\n", errorMsg, err)
	}
}

func MustChangeDirectoryTo(path string) {
	err := os.Chdir(path)

	if err != nil {
		logger.WriteErrorf("failed to change directory to %q: %s\n", path, err)
	}
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
