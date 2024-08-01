package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

// updateCmd represents the create command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the submodule branch to the specified branch name",
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateSubmoduleUpdate(branchName, repoFolderPath, submoduleName)
		if err != nil {
			logger.WriteError(err.Error())
		}

		err = filehandler.FolderArgExists(repoFolderPath, "repo-parent-path")
		if err != nil {
			logger.WriteError(err.Error())
		}

		folders := getListOfFoldersWithSubmodule(repoFolderPath, submoduleName)
		var currentBranch string
		var masterBranch string
		for _, folder := range folders {
			var pathParts = append([]string{folder}, append(pathToSubmodule, submoduleName)...)
			var submoduleDir = filepath.Join(pathParts...)
			commandhandler.MustChangeDirectoryTo(submoduleDir)

			masterBranch = getGitMasterBranch()
			currentBranch = commandhandler.MustGetCommandOutput(gitProgramName, fmt.Sprintf(`failed to get current branch for %q`, folder), getCurrentBranchArgs...)
			if strings.Contains(currentBranch, branchName) {
				logger.WriteInfof("Skipping %q since it already has %q as its branch\n", submoduleDir, branchName)
				continue
			}

			logger.WriteInfof("Updating %q's branch to %q\n", submoduleDir, branchName)

			checkoutLatestFromMaster(submoduleDir, masterBranch)

			if branchName == masterBranch {
				commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to pull latest changes for %q`, folder), "checkout", branchName)
			}

			commandhandler.MustChangeDirectoryTo(upADirectory)

			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to stage changes to %q for %q`, submoduleName, folder), "add", submoduleName)
			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to commit changes for %q`, folder), "commit", "-m", fmt.Sprintf(`"Updated %s"`, submoduleName))
			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to push changes for %q`, folder), "push")
		}
	},
}

func init() {
	submoduleCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&submoduleName, "submodule", "s", "", "the name of the submodule to operate on")
	err := updateCmd.MarkFlagRequired("submodule")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"submodule\" as required on update command: %v\n", err)
	}

	updateCmd.Flags().StringVarP(&repoFolderPath, "repo-parent-path", "d", "", "the path to the parent folder of the repos to operate on")
	err = updateCmd.MarkFlagRequired("repo-parent-path")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"repo-parent-path\" as required on update command: %v\n", err)
	}

	err = updateCmd.MarkFlagDirname("repo-parent-path")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"repo-parent-path\" as directory on update command: %v\n", err)
	}

	updateCmd.Flags().StringVarP(&branchName, "branch-name", "b", "", "the submodule branch name to checkout and use")
	err = updateCmd.MarkFlagRequired("branch-name")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"branch-name\" as required on update command: %v\n", err)
	}
}

func ValidateSubmoduleUpdate(branchName, repoFolderPath, submoduleName string) error {
	if strings.TrimSpace(branchName) == "" {
		return errors.New(BranchNameArgEmpty)
	}

	if strings.TrimSpace(repoFolderPath) == "" {
		return errors.New(RepoParentPathArgEmpty)
	}

	if strings.TrimSpace(submoduleName) == "" {
		return errors.New(SubmoduleNameArgEmpty)
	}

	return nil
}
