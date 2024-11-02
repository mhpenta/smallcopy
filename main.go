package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gobwas/glob"
)

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

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	ignoredPatterns := getIgnoredPatterns(dir)
	files := getTextFiles(dir, ignoredPatterns)

	var output strings.Builder
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		relPath, err := filepath.Rel(dir, file)
		if err != nil {
			fmt.Printf("Error getting relative path for %s: %v\n", file, err)
			continue
		}
		commentStyle := getCommentStyle(filepath.Ext(file))
		output.WriteString(fmt.Sprintf("%s %s\n", commentStyle, relPath))

		if filepath.Ext(file) == ".md" {
			content = censorBashCommands(content)
		}

		output.Write(content)
		output.WriteString("\n\n")
	}

	err = clipboard.WriteAll(output.String())
	if err != nil {
		fmt.Println("Error copying to clipboard:", err)
		return
	}

	fmt.Println("Content copied to clipboard successfully! Total bytes:", output.Len())
}

func getIgnoredPatterns(dir string) []glob.Glob {
	patterns := readGitignore(dir)
	for _, ignore := range defaultIgnores {
		g, _ := glob.Compile(ignore)
		patterns = append(patterns, g)
	}
	return patterns
}

func readGitignore(dir string) []glob.Glob {
	gitignorePath := filepath.Join(dir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var patterns []glob.Glob
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			g, err := glob.Compile(line)
			if err == nil {
				patterns = append(patterns, g)
			}
		}
	}
	return patterns
}

func getTextFiles(dir string, ignoredPatterns []glob.Glob) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			relPath, _ := filepath.Rel(dir, path)
			for _, pattern := range ignoredPatterns {
				if pattern.Match(relPath) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		relPath, _ := filepath.Rel(dir, path)
		for _, pattern := range ignoredPatterns {
			if pattern.Match(relPath) {
				return nil
			}
		}
		if isTextFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking directory:", err)
	}
	return files
}

func isTextFile(filename string) bool {
    textExtensions := []string{
        ".txt", ".md", ".py", ".go", ".js", ".ts", ".tsx", 
        ".html", ".css", ".xml", ".json", ".yaml", ".yml", 
        ".sql", ".jsx"
    }
    ext := filepath.Ext(filename)
    for _, textExt := range textExtensions {
        if ext == textExt {
            return true
        }
    }
    return false
}

func getCommentStyle(ext string) string {
    switch ext {
    case ".py":
        return "#"
    case ".go", ".js", ".ts", ".tsx", ".jsx":
        return "//"
    case ".sql":
        return "--"
    default:
        return "#"
    }
}

func censorBashCommands(content []byte) []byte {
	re := regexp.MustCompile("(?s)```bash\\n(.*?)```")
	return re.ReplaceAllFunc(content, func(match []byte) []byte {
		lines := strings.Split(string(match), "\n")
		censored := make([]string, len(lines))
		for i, line := range lines {
			if i == 0 || i == len(lines)-1 {
				censored[i] = line
			} else {
				censored[i] = "[CENSORED BASH COMMAND]"
			}
		}
		return []byte(strings.Join(censored, "\n"))
	})
}
