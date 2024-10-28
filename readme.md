## Smallcopy

Copy the text contents of a directory and its subdirectory safely, in bulk. Skipped files include those like: 

```go
var defaultIgnores = []string{
	// Virtual environments
	"venv", "virtualenv", "env", ".venv", ".env",
	// IDE and editor files
	".idea", ".vscode", "*.sublime-*", "*.swp", ".DS_Store",
	// Specific IDE files
	".idea/workspace.xml", ".idea/tasks.xml", ".idea/dictionaries",
	".idea/inspectionProfiles/profiles_settings.xml",
	// Build outputs
	"build", "dist", "*.egg-info", "*.pyc", "__pycache__",
	// Dependency directories
	"node_modules", "bower_components",
	// Log files
	"*.log",
	// OS generated files
	"Thumbs.db", ".DS_Store", "*.bak",
	// CI/CD configuration files
	".github/workflows/*", ".gitlab-ci.yml", ".travis.yml", "azure-pipelines.yml",
	// Specific CI/CD files
	".github/workflows/gcs_deploy.yaml",
	// Docker files
	"Dockerfile", "docker-compose.yml", "docker-compose.*.yml",
	// Configuration files
	"*.config.js", "*.config.ts", "*.conf", "*.cfg",
	// Package manager files
	"package-lock.json", "yarn.lock", "Pipfile.lock",
	// Database files
	"*.sqlite", "*.db",
	// Temporary files
	"*.tmp", "*.temp",
}
```

See implementation details in main.go to ensure the copy works as required. 

## Install 

```bash
go build -o smallcopy && sudo mv smallcopy /usr/local/bin/
```

## Usage


In the directory you want to copy, run:
```bash
smallcopy
```

