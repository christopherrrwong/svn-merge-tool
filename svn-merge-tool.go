package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Helper function to get environment variable with default fallback
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


var (
	RepoBase      string
	BranchName    string
	LocalRepoPath string
	TrunkUrl      string
	BranchUrl     string
	localTrunkPath    string
	localBranchPath   string
	BranchNameDisplay string
)

func init() {
	// Load .env file FIRST
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using defaults")
	}
	
	// THEN initialize variables after .env is loaded, please add default values for the variables if needed 
	RepoBase = getEnvWithDefault("REPO_BASE", "")
	BranchName = getEnvWithDefault("BRANCH_PATH", "")
	LocalRepoPath = getEnvWithDefault("LOCAL_REPO_PATH", "")
	
	// Derived variables
	TrunkUrl = RepoBase + "/trunk"
	BranchUrl = RepoBase + BranchName
	localTrunkPath = LocalRepoPath + "/trunk"
	localBranchPath = LocalRepoPath + BranchName
	BranchNameDisplay = BranchName
	
	fmt.Println("=== Configuration ===")
	fmt.Println("REPO_BASE:", RepoBase)
	fmt.Println("BRANCH_PATH:", BranchName)
	fmt.Println("LOCAL_REPO_PATH:", LocalRepoPath)
	fmt.Println("TrunkUrl:", TrunkUrl)
	fmt.Println("BranchUrl:", BranchUrl)
	fmt.Println("=====================")
}

type Revision struct {
	Number  string
	Author  string
	Date    string
	Message string
}


func UpdateLocalRepo(path string) error {
	fmt.Printf("Updating local repo at: %s\n", path)

	cmd := exec.Command("svn", "update", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	fmt.Println("Update successful")
	return nil
}

func GetTop10Revisions(path string) ([]Revision, error) {
	// First update the local repository to get latest changes
	fmt.Printf("Updating local repository: %s\n", path)
	if err := UpdateLocalRepo(path); err != nil {
		fmt.Printf("Warning: Failed to update local repo: %v\n", err)
	}

	// Change to the SVN directory
	if err := os.Chdir(path); err != nil {
		return nil, fmt.Errorf("failed to change directory: %v", err)
	}

	// Get last 10 revisions
	cmd := exec.Command("svn", "log", "-l", "10")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("svn log failed: %v", err)
	}

	return parseSVNLog(string(output))
}

func parseSVNLog(logOutput string) ([]Revision, error) {
	var revisions []Revision

	
	entries := strings.Split(logOutput, "------------------------------------------------------------------------")


	revRegex := regexp.MustCompile(`^r(\d+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*\d+\s*lines?\s*$`)

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		lines := strings.Split(entry, "\n")
		if len(lines) < 2 {
			continue
		}

	
		headerLine := strings.TrimSpace(lines[0])
		matches := revRegex.FindStringSubmatch(headerLine)

		if len(matches) >= 4 {
	
			var messageLines []string
			for i := 1; i < len(lines); i++ {
				line := strings.TrimSpace(lines[i])
				if line != "" {
					messageLines = append(messageLines, line)
				}
			}

			message := strings.Join(messageLines, " ")
	
			if len(message) > 60 {
				message = message[:57] + "..."
			}

			revision := Revision{
				Number:  "r" + matches[1],
				Author:  strings.TrimSpace(matches[2]),
				Date:    strings.TrimSpace(matches[3]),
				Message: message,
			}

			revisions = append(revisions, revision)
		}
	}

	return revisions, nil
}

func GetUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func DisplayRevisions(revisions []Revision) {
	fmt.Println("\n=== Available Revisions ===")
	fmt.Printf("%-3s | %-8s | %-15s | %s\n", "No.", "Revision", "Author", "Message")
	fmt.Println("----+----------+-----------------+----------------------------------------")

	for i, rev := range revisions {
		author := rev.Author
		if len(author) > 15 {
			author = author[:12] + "..."
		}

		fmt.Printf("%-3d | %-8s | %-15s | %s\n", i+1, rev.Number, author, rev.Message)
	}
	fmt.Println()
}

func SelectRevisionMode() int {
	fmt.Println("Select commit message type:")
	fmt.Println("1. Single revision")
	fmt.Println("2. Revision range (2 revisions)")

	for {
		choice := GetUserInput("Enter choice (1 or 2): ")
		switch choice {
		case "1":
			return 1
		case "2":
			return 2
		default:
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

func SelectSingleRevision(revisions []Revision) string {
	for {
		input := GetUserInput(fmt.Sprintf("Select revision number (1-%d): ", len(revisions)))
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(revisions) {
			fmt.Printf("Invalid choice. Please enter a number between 1 and %d.\n", len(revisions))
			continue
		}

		selected := revisions[choice-1]
		fmt.Printf("Selected: %s - %s\n", selected.Number, selected.Message)
		return selected.Number
	}
}

func SelectTwoRevisions(revisions []Revision) (string, string) {
	var first, second string

	fmt.Println("Select FIRST revision:")
	for {
		input := GetUserInput(fmt.Sprintf("Enter number (1-%d): ", len(revisions)))
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(revisions) {
			fmt.Printf("Invalid choice. Please enter a number between 1 and %d.\n", len(revisions))
			continue
		}

		selected := revisions[choice-1]
		first = selected.Number
		fmt.Printf("First: %s - %s\n", selected.Number, selected.Message)
		break
	}

	fmt.Println("\nSelect SECOND revision:")
	for {
		input := GetUserInput(fmt.Sprintf("Enter number (1-%d): ", len(revisions)))
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(revisions) {
			fmt.Printf("Invalid choice. Please enter a number between 1 and %d.\n", len(revisions))
			continue
		}

		selected := revisions[choice-1]
		second = selected.Number
		fmt.Printf("Second: %s - %s\n", selected.Number, selected.Message)
		break
	}

	return first, second
}

func DetectMergeDirection() (string, error) {

	// // If no pending changes, ask user
	fmt.Println("Please specify merge direction:")
	fmt.Println("1. Branch to Trunk")
	fmt.Println("2. Trunk to Branch")

	for {
		choice := GetUserInput("Enter choice (1 or 2): ")
		switch choice {
		case "1":
			return "branch-to-trunk", nil
		case "2":
			return "trunk-to-branch", nil
		default:
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

func GenerateCommitMessage(mode int, direction string, startRev, endRev string, taskNumber string) string {
	switch mode {
	case 1: // Single revision
		if direction == "trunk-to-branch" {
			return fmt.Sprintf("Merged %s from trunk to %s", startRev, BranchNameDisplay)
		} else {
			return fmt.Sprintf("Merged %s from %s to trunk (%s)", startRev, BranchNameDisplay, taskNumber)
		}
	case 2: // Range
		if direction == "trunk-to-branch" {
			return fmt.Sprintf("Merged branches %s - %s from trunk to %s", startRev, endRev, BranchNameDisplay)
		} else {
			return fmt.Sprintf("Merged branches %s - %s from %s to trunk (%s)", startRev, endRev, BranchNameDisplay, taskNumber)
		}
	}
	return "Merge commit"
}

func GenerateTaskNumber() string {
	fmt.Println("Branch to trunk required. Please insert your task number e.g, T123456.")
	taskNumber := GetUserInput("Enter task number: ")
	return taskNumber
}

func InteractiveCommitMessage() error {
	fmt.Println("=== Generate Commit Message ===")

	// Detect merge direction
	direction, err := DetectMergeDirection()
	if err != nil {
		return err
	}

	fmt.Printf("Merge direction: %s\n", direction)

	var revisions []Revision
	var sourcePath string

	// Get revisions from source location
	if direction == "trunk-to-branch" {
		sourcePath = localTrunkPath
		fmt.Println("Getting trunk revisions...")
	} else {
		sourcePath = localBranchPath
		fmt.Println("Getting branch revisions...")
	}

	revisions, err = GetTop10Revisions(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to get revisions: %v", err)
	}

	if len(revisions) == 0 {
		return fmt.Errorf("no revisions found")
	}


	DisplayRevisions(revisions)


	mode := SelectRevisionMode()

	//commit message for branch to trunk requried task number
	var commitMessage string
	switch mode {
	case 1:
		// Single revision

		selected := SelectSingleRevision(revisions)

		if direction == "branch-to-trunk" {
			taskNumber := GenerateTaskNumber()
			commitMessage = GenerateCommitMessage(1, direction, selected, "", taskNumber)
		} else {
			commitMessage = GenerateCommitMessage(1, direction, selected, "", "")
		}

	case 2:
		// Two revisions
		first, second := SelectTwoRevisions(revisions)
		if direction == "branch-to-trunk" {
			taskNumber := GenerateTaskNumber()
			commitMessage = GenerateCommitMessage(2, direction, first, second, taskNumber)
		} else {
			commitMessage = GenerateCommitMessage(2, direction, first, second, "")
		}

	}

	
	fmt.Printf("\nGenerated commit message:\n\"%s\"\n\n", commitMessage)

	// Ask for confirmation before committing
	confirm := GetUserInput("Do you want to commit with this message? (y/N): ")
	if strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes" {
		return CommitWithMessage(commitMessage, direction)
	}

	fmt.Println("Commit cancelled. You can commit manually later with:")
	fmt.Printf("svn commit -m \"%s\"\n", commitMessage)
	return nil
}

func CommitWithMessage(message string, direction string) error {
	if direction != "branch-to-trunk" && direction != "trunk-to-branch" {
		return fmt.Errorf("invalid direction: %s", direction)
	}

	var targetDir string
	if direction == "branch-to-trunk" {
		targetDir = localTrunkPath
		fmt.Printf("Committing changes in trunk directory: %s\n", targetDir)
	} else {
		targetDir = localBranchPath
		fmt.Printf("Committing changes in branch directory: %s\n", targetDir)
	}

	// Create the SVN commit command and set the working directory
	cmd := exec.Command("svn", "commit", "-m", message)
	cmd.Dir = targetDir  // Set the working directory for the command
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}

	fmt.Println("Commit successful!")
	return nil
}

// update trunk and branch first to make sure we have latest changes and merge branch to trunk
func MergeBranchToTrunk() error {
	fmt.Println("=== Starting merge of branch to trunk ===")


	if err := UpdateLocalRepo(localTrunkPath); err != nil {
		return err
	}

	
	if err := UpdateLocalRepo(localBranchPath); err != nil {
		return err
	}

	fmt.Println("Merging branch to trunk...")


	if err := os.Chdir(localTrunkPath); err != nil {
		return fmt.Errorf("failed to change to trunk directory: %v", err)
	}


	cmd := exec.Command("svn", "merge", BranchUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("merge failed: %v", err)
	}

	fmt.Println("Merge successful!")


	fmt.Println("\n=== Checking merge status ===")
	statusCmd := exec.Command("svn", "status")
	statusCmd.Stdout = os.Stdout
	statusCmd.Run()

	fmt.Println("\nMerge completed! To commit, run:")
	fmt.Println("go run svn-merge-tool.go generate-commit")

	return nil
}

// update trunk and branch first to make sure we have latest changes and merge trunk to branch
func MergeTrunkToBranch() error {
	fmt.Println("=== Starting merge of trunk to branch ===")

	
	if err := UpdateLocalRepo(localBranchPath); err != nil {
		return err
	}

	if err := UpdateLocalRepo(localTrunkPath); err != nil {
		return err
	}

	fmt.Println("Merging trunk to branch...")


	if err := os.Chdir(localBranchPath); err != nil {
		return fmt.Errorf("failed to change to branch directory: %v", err)
	}

	// Merge trunk into current directory (branch)
	cmd := exec.Command("svn", "merge", TrunkUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("merge failed: %v", err)
	}

	fmt.Println("Merge successful!")


	fmt.Println("\n=== Checking merge status ===")
	statusCmd := exec.Command("svn", "status")
	statusCmd.Stdout = os.Stdout
	statusCmd.Run()

	fmt.Println("\nMerge completed! To commit, run:")
	fmt.Println("go run svn-merge-tool.go generate-commit")

	return nil
}

// get recent revisions for trunk and branch
func GetRecentRevisions(path string, url string) error {
	fmt.Printf("=== Getting recent revisions for: %s ===\n", url)

	// First update local repository
	fmt.Printf("Updating local repository: %s\n", path)
	if err := UpdateLocalRepo(path); err != nil {
		fmt.Printf("Warning: Failed to update local repo: %v\n", err)
	}

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change to directory %s: %v", path, err)
	}

	cmd := exec.Command("svn", "log", "-l", "10", "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get revisions: %v", err)
	}

	return nil
}

// Get revisions directly from remote repository URL (most up-to-date)
func GetRecentRevisionsFromRemote(url string) error {
	fmt.Printf("=== Getting latest revisions directly from remote: %s ===\n", url)

	cmd := exec.Command("svn", "log", "-l", "10", "--verbose", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get remote revisions: %v", err)
	}

	return nil
}

func showHelp() {
	fmt.Println("SVN Merge Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run svn-merge-tool.go <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  branch-to-trunk    Merge branch to trunk (merge only)")
	fmt.Println("  trunk-to-branch    Merge trunk to branch (merge only)")
	fmt.Println("  generate-commit    Generate commit message and commit")
	fmt.Println("  revisions-trunk    Show last 10 revisions for trunk")
	fmt.Println("  revisions-branch   Show last 10 revisions for branch")
	fmt.Println("  help               Show this help message")
	fmt.Println("")
	fmt.Println("Workflow:")
	fmt.Println("  1. go run svn-merge-tool.go branch-to-trunk    # Do the merge")
	fmt.Println("  2. go run svn-merge-tool.go generate-commit    # Generate commit message")
	fmt.Println("")
	fmt.Println("Short aliases:")
	fmt.Println("  btt = branch-to-trunk")
	fmt.Println("  ttb = trunk-to-branch")
	fmt.Println("  gc  = generate-commit")
}

func main() {


	if len(os.Args) < 2 {
		fmt.Println("Error: No command specified")
		fmt.Println("")
		showHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "branch-to-trunk", "btt":
		if err := MergeBranchToTrunk(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "trunk-to-branch", "ttb":
		if err := MergeTrunkToBranch(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "generate-commit", "gc":
		if err := InteractiveCommitMessage(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "revisions-trunk", "rt":
		if err := GetRecentRevisions(localTrunkPath, TrunkUrl); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "revisions-branch", "rb":
		if err := GetRecentRevisions(localBranchPath, BranchUrl); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		showHelp()
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("")
		showHelp()
		os.Exit(1)
	}
}
