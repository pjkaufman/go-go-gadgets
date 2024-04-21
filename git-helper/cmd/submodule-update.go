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

		filehandler.FolderMustExist(repoFolderPath, "repo-parent-path")

		folders := getListOfFoldersWithSubmodule(repoFolderPath, submoduleName)
		var currentBranch string
		var masterBranch string
		for _, folder := range folders {
			var pathParts = append([]string{folder}, append(pathToSubmodule, submoduleName)...)
			var submoduleDir = filepath.Join(pathParts...)
			commandhandler.MustChangeDirectoryTo(submoduleDir)

			masterBranch = getGitMasterBranch()
			currentBranch = commandhandler.MustGetCommandOutput(gitProgramName, fmt.Sprintf(`failed to get current branch for "%s"`, folder), getCurrentBranchArgs...)
			if strings.Contains(currentBranch, branchName) {
				logger.WriteInfo(fmt.Sprintf(`Skipping "%s" since it already has "%s" as its branch`, submoduleDir, branchName))
				continue
			}

			logger.WriteInfo(fmt.Sprintf(`Updating "%s"'s branch to "%s"`, submoduleDir, branchName))

			checkoutLatestFromMaster(submoduleDir, masterBranch)

			if branchName == masterBranch {
				commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to pull latest changes for "%s"`, folder), "checkout", branchName)
			}

			commandhandler.MustChangeDirectoryTo(upADirectory)

			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to stage changes to "%s" for "%s"`, submoduleName, folder), "add", submoduleName)
			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to commit changes for "%s"`, folder), "commit", "-m", fmt.Sprintf(`"Updated %s"`, submoduleName))
			commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to push changes for "%s"`, folder), "push")
		}
	},
}

func init() {
	submoduleCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&submoduleName, "submodule", "s", "", "the name of the submodule to operate on")
	updateCmd.Flags().StringVarP(&repoFolderPath, "repo-parent-path", "d", "", "the path to the parent folder of the repos to operate on")
	updateCmd.Flags().StringVarP(&branchName, "branch-name", "b", "", "the submodule branch name to checkout and use")
	updateCmd.MarkFlagRequired("submodule")
	updateCmd.MarkFlagRequired("repo-parent-path")
	updateCmd.MarkFlagRequired("branch-name")
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
