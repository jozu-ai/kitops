package filesystem

import (
	"kitops/pkg/artifact"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnoreMatches(t *testing.T) {
	tests := []struct {
		name         string
		kitIgnore    []string
		layerPaths   []string
		curPath      string
		curLayerPath string
		shouldIgnore bool
	}{
		{
			name:         "Ignores files in directory",
			kitIgnore:    []string{"dir1"},
			layerPaths:   []string{},
			curPath:      "dir1/subdir1/subdir2/file.txt",
			curLayerPath: "",
			shouldIgnore: true,
		},
		{
			name:         "Ignores files with wildcard",
			kitIgnore:    []string{"dir1/*.txt"},
			layerPaths:   []string{},
			curPath:      "dir1/testfile.txt",
			curLayerPath: "",
			shouldIgnore: true,
		},
		{
			name:         "Ignores files with '**' wildcard",
			kitIgnore:    []string{"**/testfile.txt"},
			layerPaths:   []string{},
			curPath:      "dir1/subdir1/subdir2/testfile.txt",
			curLayerPath: "",
			shouldIgnore: true,
		},
		{
			name:         "Can explicitly include files",
			kitIgnore:    []string{"dir1", "!dir1/testfile.txt"},
			layerPaths:   []string{},
			curPath:      "dir1/testfile.txt",
			curLayerPath: "",
			shouldIgnore: false,
		},
		{
			name:         "Test intersecting layers exclusion",
			kitIgnore:    []string{},
			layerPaths:   []string{"main", "main/subdir"},
			curPath:      "main/subdir/testfile.txt",
			curLayerPath: "main",
			shouldIgnore: true,
		},
		{
			name:         "Test intersecting layers inclusion",
			kitIgnore:    []string{},
			layerPaths:   []string{"main", "main/subdir"},
			curPath:      "main/subdir/testfile.txt",
			curLayerPath: "main/subdir",
			shouldIgnore: false,
		},
		{
			name:         "Test intersecting layers inclusion with kitignore",
			kitIgnore:    []string{"**/testfile.txt"},
			layerPaths:   []string{"main", "main/subdir"},
			curPath:      "main/subdir/testfile.txt",
			curLayerPath: "main/subdir",
			shouldIgnore: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testKitfile := &artifact.KitFile{}
			for _, layerPath := range tt.layerPaths {
				testKitfile.Code = append(testKitfile.Code, artifact.Code{Path: layerPath})
			}
			ignore, err := NewIgnore(tt.kitIgnore, testKitfile)
			assert.NoError(t, err)

			ignored, err := ignore.Matches(tt.curPath, tt.curLayerPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.shouldIgnore, ignored)
		})
	}
}
