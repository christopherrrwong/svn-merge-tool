# SVN Merge Tool

A command-line tool written in Go that simplifies SVN branch-to-trunk and trunk-to-branch merging operations with interactive commit message generation.

## Features

- 🔄 **Automated Merging**: Seamlessly merge between SVN branches and trunk
- 📝 **Interactive Commit Messages**: Generate standardized commit messages with revision selection
- 📊 **Revision History**: View recent revisions from trunk or branch
- 🔧 **Configurable**: Environment-based configuration for different repositories
- ⚡ **Fast Workflow**: Streamlined commands for common SVN operations

## Prerequisites

- Go 1.16 or higher
- SVN (Subversion) command-line client
- Access to your SVN repository

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd svn-merge-tool
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the tool (optional):
```bash
go build -o svn-merge-tool svn-merge-tool.go
```

## Configuration

Create a `.env` file in the project root with your SVN repository configuration:

```env
# SVN Repository base URL (without /trunk or /branches)
REPO_BASE=https://your-svn-server.com/svn/your-project

# Branch path (e.g., /branches/feature-branch or /branches/release-1.0)
BRANCH_PATH=/branches/your-branch-name

# Local repository path where trunk and branch directories are located
LOCAL_REPO_PATH=/path/to/your/local/svn/workspace
```

### Example `.env` file:
```env
REPO_BASE=https://svn.company.com/projects/myproject
BRANCH_PATH=/branches/feature-authentication
LOCAL_REPO_PATH=/Users/developer/workspace/myproject
```

The tool expects your local workspace to have this structure:
```
/your/local/workspace/
├── trunk/          # SVN trunk checkout
└── branches/
    └── your-branch-name/  # SVN branch checkout
```

## Usage

### Available Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `branch-to-trunk` | `btt` | Merge branch changes to trunk |
| `trunk-to-branch` | `ttb` | Merge trunk changes to branch |
| `generate-commit` | `gc` | Generate commit message and commit changes |
| `revisions-trunk` | `rt` | Show last 10 revisions for trunk |
| `revisions-branch` | `rb` | Show last 10 revisions for branch |
| `help` | `-h`, `--help` | Show help message |

### Basic Workflow

#### 1. Merge Branch to Trunk
```bash
# Step 1: Perform the merge
go run svn-merge-tool.go branch-to-trunk

# Step 2: Generate commit message and commit
go run svn-merge-tool.go generate-commit
```

#### 2. Merge Trunk to Branch
```bash
# Step 1: Perform the merge
go run svn-merge-tool.go trunk-to-branch

# Step 2: Generate commit message and commit
go run svn-merge-tool.go generate-commit
```

### Interactive Commit Message Generation

The `generate-commit` command provides an interactive interface to:

1. **Select merge direction**: Branch-to-trunk or trunk-to-branch
2. **Choose revision type**: Single revision or revision range
3. **Select revisions**: Pick from the last 10 revisions with detailed information
4. **Add task number**: For branch-to-trunk merges (e.g., T123456)
5. **Review and confirm**: Preview the generated commit message before committing

#### Example Commit Messages

**Single Revision (Branch to Trunk):**
```
Merged r12345 from /branches/feature-auth to trunk (T123456)
```

**Revision Range (Trunk to Branch):**
```
Merged branches r12340 - r12345 from trunk to /branches/feature-auth
```

### Viewing Revision History

```bash
# View trunk revisions
go run svn-merge-tool.go revisions-trunk

# View branch revisions
go run svn-merge-tool.go revisions-branch
```

## Command Examples

### Using Full Command Names
```bash
go run svn-merge-tool.go branch-to-trunk
go run svn-merge-tool.go generate-commit
go run svn-merge-tool.go revisions-trunk
```

### Using Aliases
```bash
go run svn-merge-tool.go btt
go run svn-merge-tool.go gc
go run svn-merge-tool.go rt
```

### If Built as Binary
```bash
./svn-merge-tool btt
./svn-merge-tool gc
```



## Troubleshooting

### Common Issues

1. **"No .env file found"**
   - Create a `.env` file with your repository configuration
   - Ensure the file is in the same directory as the Go script

2. **"Failed to change directory"**
   - Verify your `LOCAL_REPO_PATH` is correct
   - Ensure trunk and branch directories exist and are SVN working copies

3. **"SVN command failed"**
   - Check that SVN is installed and accessible from command line
   - Verify repository URLs are accessible
   - Ensure you have proper SVN credentials



## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

