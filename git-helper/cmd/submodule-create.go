package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MakeNowJust/heredoc"
	commandhandler "github.com/pjkaufman/go-go-gadgets/pkg/command-handler"
	filehandler "github.com/pjkaufman/go-go-gadgets/pkg/file-handler"
	"github.com/pjkaufman/go-go-gadgets/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	submoduleName           string
	repoFolderPath          string
	ticketAbbreviation      string
	branchName              string
	branchPrefix            string
	pathToSubmodule         = []string{"src", "modules"}
	getCurrentBranchArgs    = []string{"branch", "--show-current"}
	getMasterBranchNameArgs = []string{"symbolic-ref", "--short", "refs/remotes/origin/HEAD"}
)

var prLinkRegex = regexp.MustCompile(`https:[^\n]*`)

const (
	TicketArgEmpty         = "ticket-abbreviation must have a non-whitespace value"
	BranchNameArgEmpty     = "branch-name must have a non-whitespace value"
	RepoParentPathArgEmpty = "repo-parent-path must have a non-whitespace value"
	SubmoduleNameArgEmpty  = "submodule must have a non-whitespace value"
	BranchPrefixArgEmpty   = "branch-prefix must have a non-whitespace value"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates the branch in the specified submodule if it does not already exist",
	Example: heredoc.Doc(`git-tools submodule create -s Submodule -d ./repos/ -a abbrev -b fix-bug -p fix
	will go ahead and look at all git repos in the folder repos with the submodule called Submodule and check if that repo currently has "abbrev" in the current branch name.
	If it does not, it will create a branch with the submodule branch set to "fix-bug" and push those changes up on the regular repo with a branch name of fix/abbrev-update-Submodule.`),
	Long: `Creates the specified branch in the provided submodule for all instances of the submodule in the provided folder so long as it is not already present.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := ValidateSubmoduleCreate(ticketAbbreviation, branchName, repoFolderPath, submoduleName, branchPrefix)
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
		var prLinks []string
		for _, folder := range folders {
			commandhandler.MustChangeDirectoryTo(folder)

			masterBranch = getGitMasterBranch()
			currentBranch = commandhandler.MustGetCommandOutput(gitProgramName, fmt.Sprintf(`failed to get current branch for %q`, folder), getCurrentBranchArgs...)
			if strings.Contains(currentBranch, ticketAbbreviation) {
				continue
			}

			currentBranch = strings.TrimSpace(currentBranch)
			logger.WriteInfof("%q does not contain %q\n", currentBranch, ticketAbbreviation)

			prLinks = append(prLinks, createSubmoduleUpdateBranch(folder, submoduleName, branchPrefix, masterBranch))
		}

		if len(prLinks) != 0 {
			logger.WriteInfo("\nPR Links:")

			for _, link := range prLinks {
				logger.WriteInfo("- " + link)
			}
		}
	},
}

func init() {
	submoduleCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&submoduleName, "submodule", "s", "", "the name of the submodule to operate on")
	err := createCmd.MarkFlagRequired("submodule")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"submodule\" as required on create command: %v\n", err)
	}

	createCmd.Flags().StringVarP(&repoFolderPath, "repo-parent-path", "d", "", "the path to the parent folder of the repos to operate on")
	err = createCmd.MarkFlagRequired("repo-parent-path")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"submodule\" as required on create command: %v\n", err)
	}

	err = createCmd.MarkFlagDirname("repo-parent-path")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"repo-parent-path\" as directory on create command: %v\n", err)
	}

	createCmd.Flags().StringVarP(&ticketAbbreviation, "ticket-abbreviation", "a", "", "the ticket abbreviation to use to determine whether we should update a repo and to help determine the name for submodule branch")
	err = createCmd.MarkFlagRequired("ticket-abbreviation")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"ticket-abbreviation\" as required on create command: %v\n", err)
	}

	createCmd.Flags().StringVarP(&branchName, "branch-name", "b", "", "the submodule branch name to checkout and use")
	err = createCmd.MarkFlagRequired("branch-name")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"branch-name\" as required on create command: %v\n", err)
	}

	createCmd.Flags().StringVarP(&branchPrefix, "branch-prefix", "p", "", "the branch prefix to use for the created branch names")
	err = createCmd.MarkFlagRequired("branch-prefix")
	if err != nil {
		logger.WriteErrorf("failed to mark flag \"branch-prefix\" as required on create command: %v\n", err)
	}
}

func getListOfFoldersWithSubmodule(path, submoduleName string) []string {
	var folders []string
	allFolders, err := filehandler.GetFoldersInCurrentFolder(path)
	if err != nil {
		logger.WriteError(err.Error())
	}

	for _, dir := range allFolders {
		var pathParts = []string{path, dir}
		var folderPath = filepath.Join(pathParts...)
		pathParts = append(pathParts, pathToSubmodule...)
		pathParts = append(pathParts, submoduleName)
		var submoduleFolderPath = filepath.Join(pathParts...)

		exists, err := filehandler.FolderExists(submoduleFolderPath)
		if err != nil {
			logger.WriteError(err.Error())
		}

		if !exists {
			continue
		}

		folders = append(folders, folderPath)
	}

	return folders
}

func createSubmoduleUpdateBranch(folder, submodule, branchPrefix, masterBranch string) string {
	logger.WriteInfo("Creating the DE branch for " + folder)
	checkoutLatestFromMaster(folder, masterBranch)

	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to pull latest changes for %q`, folder), "checkout", "-B", branchPrefix+"/"+ticketAbbreviation+"-update-"+submodule)

	var submoduleDir = filepath.Join(append(pathToSubmodule, submodule)...)
	commandhandler.MustChangeDirectoryTo(filepath.Join(append(pathToSubmodule, submodule)...))

	checkoutLatestFromMaster(submoduleDir, masterBranch)

	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to checkout %q for %q`, branchName, folder), "checkout", branchName)

	commandhandler.MustChangeDirectoryTo(upADirectory)

	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to stage changes to %q for %q`, submodule, folder), "add", submodule)
	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to commit changes for %q`, folder), "commit", "-m", fmt.Sprintf(`"Updated %s"`, submodule))
	pushOutput := commandhandler.MustGetCommandOutput(gitProgramName, fmt.Sprintf(`failed to push changes for %q`, folder), "push")

	return GetPullRequestLink(pushOutput)
}

func checkoutLatestFromMaster(folder, masterBranch string) {
	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to checkout master for %q`, folder), "checkout", masterBranch)
	commandhandler.MustRunCommand(gitProgramName, fmt.Sprintf(`failed to pull latest changes for %q`, folder), "pull")
}

func ValidateSubmoduleCreate(ticketAbbreviation, branchName, repoFolderPath, submoduleName, branchPrefix string) error {
	if strings.TrimSpace(ticketAbbreviation) == "" {
		return errors.New(TicketArgEmpty)
	}

	if strings.TrimSpace(branchName) == "" {
		return errors.New(BranchNameArgEmpty)
	}

	if strings.TrimSpace(repoFolderPath) == "" {
		return errors.New(RepoParentPathArgEmpty)
	}

	if strings.TrimSpace(submoduleName) == "" {
		return errors.New(SubmoduleNameArgEmpty)
	}

	if strings.TrimSpace(branchPrefix) == "" {
		return errors.New(BranchPrefixArgEmpty)
	}

	return nil
}

func GetPullRequestLink(pushOutput string) string {
	var matches = prLinkRegex.FindAllString(pushOutput, 1)
	if len(matches) == 0 {
		return ""
	}

	return matches[0]
}

func getGitMasterBranch() string {
	shortBranch := commandhandler.MustGetCommandOutput(gitProgramName, "failed to get master branch name", getMasterBranchNameArgs...)

	actualBranchIndex := strings.Index(shortBranch, "/")
	var actualBranch = shortBranch
	if actualBranchIndex != -1 {
		actualBranch = shortBranch[actualBranchIndex+1:]
	}

	actualBranch = strings.TrimSpace(actualBranch)
	if actualBranch == "" {
		logger.WriteError("failed to get master branch name as it is empty")
	}

	return actualBranch
}
