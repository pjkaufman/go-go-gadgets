package commandhandler

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
)

func MustGetCommandOutput(programName, errorMsg string, args ...string) string {
	cmd := exec.Command(programName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.WriteError(fmt.Sprintf(`%s: %s`, errorMsg, err))
	}

	return string(output)
}

func MustRunCommand(programName, errorMsg string, args ...string) {
	cmd := exec.Command(programName, args...)
	err := cmd.Run()
	if err != nil {
		logger.WriteError(fmt.Sprintf(`%s: %s`, errorMsg, err))
	}
}

func MustChangeDirectoryTo(path string) {
	err := os.Chdir(path)

	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to change directory to %q: %s`, path, err))
	}
}

func GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

func MustGetUserConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		logger.WriteError(fmt.Sprintf(`failed to get user config directory: %s`, err))
	}

	return configDir
}
