package filesystem

import (
	"fmt"
	"io/fs"
	"kitops/pkg/lib/constants"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// VerifySubpath checks that path.Join(context, subDir) is a subdirectory of context, following
// symlinks if present.
func VerifySubpath(context, subDir string) (absPath string, err error) {
	// Get absolute path for context and context + subDir
	absContext, err := filepath.Abs(context)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s: %w", context, err)
	}
	fullPath := filepath.Clean(filepath.Join(absContext, subDir))

	// Get actual paths, ignoring symlinks along the way
	resolvedContext, err := filepath.EvalSymlinks(absContext)
	if err != nil {
		return "", fmt.Errorf("error resolving %s: %w", absContext, err)
	}
	resolvedFullPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		return "", fmt.Errorf("error resolving %s: %w", absContext, err)
	}

	// Get relative path between context and the full path to check if the
	// actual full, absolute path is a subdirectory of context
	relPath, err := filepath.Rel(resolvedContext, resolvedFullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	if strings.Contains(relPath, "..") {
		return "", fmt.Errorf("paths must be within context directory")
	}

	return resolvedFullPath, nil
}

func PathExists(path string) (fs.FileInfo, bool) {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil, false
	}
	return fi, true
}

// Searches for a kit file in the given context directory.
// It checks for accepted kitfile names and returns the path of the first found kitfile.
// If no kitfile is found, it returns the path (contextDir + kitfile)
// of the default kitfile.
func FindKitfileInPath(contextDir string) string {
	var defaultKitFileNames = []string{"Kitfile", "kitfile", ".kitfile"}
	for _, fileName := range defaultKitFileNames {
		if _, exists := PathExists(filepath.Join(contextDir, fileName)); exists {
			return path.Join(contextDir, fileName)
		}
	}
	return path.Join(contextDir, constants.DefaultKitFileName)
}
