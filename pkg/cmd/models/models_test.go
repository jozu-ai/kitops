package models

import (
	"fmt"
	"jmm/pkg/artifact"
	"jmm/pkg/lib/constants"
	internal "jmm/pkg/lib/testing"
	"testing"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestListModels(t *testing.T) {
	tests := []struct {
		name                  string
		manifests             map[digest.Digest]ocispec.Manifest
		configs               map[digest.Digest]artifact.JozuFile
		index                 *ocispec.Index
		expectedOutputRegexps []string
		expectErrRegexp       string
	}{
		{
			name:            "Cannot read index.json",
			index:           nil,
			expectErrRegexp: "artifact not found",
		},
		{
			name: "Cannot find manifest from index.json",
			index: &ocispec.Index{
				Manifests: []ocispec.Descriptor{
					ManifestDesc("manifestA", true),
					ManifestDesc("manifestNotFound", true),
				},
			},
			manifests: map[digest.Digest]ocispec.Manifest{
				"manifestA": Manifest("configA", "layerA"),
				"manifestB": Manifest("configB", "layerB"),
			},
			configs: map[digest.Digest]artifact.JozuFile{
				"configA": Config("maintainerA", "formatA"),
				"configB": Config("maintainerB", "formatB"),
			},
			expectErrRegexp: "failed to read manifest manifestNotFound.*",
		},
		{
			name: "Cannot find config in manifest",
			index: &ocispec.Index{
				Manifests: []ocispec.Descriptor{
					ManifestDesc("manifestA", true),
					ManifestDesc("manifestB", true),
				},
			},
			manifests: map[digest.Digest]ocispec.Manifest{
				"manifestA": Manifest("configA", "layerA"),
				"manifestB": Manifest("configNotFound", "layerB"),
			},
			configs: map[digest.Digest]artifact.JozuFile{
				"configA": Config("maintainerA", "formatA"),
				"configB": Config("maintainerB", "formatB"),
			},
			expectErrRegexp: "failed to read config.*",
		},
		{
			name: "Prints summary of for each manifest line (layers are 1024 bytes)",
			index: &ocispec.Index{
				Manifests: []ocispec.Descriptor{
					ManifestDesc("manifestA", true),
					ManifestDesc("manifestB", true),
				},
			},
			manifests: map[digest.Digest]ocispec.Manifest{
				"manifestA": Manifest("configA", "layerA"),
				"manifestB": Manifest("configB", "layerB1", "layerB2", "layerB3"),
			},
			configs: map[digest.Digest]artifact.JozuFile{
				"configA": Config("maintainerA", "formatA"),
				"configB": Config("maintainerB", "formatB"),
			},
			expectedOutputRegexps: []string{
				"manifestA.*maintainerA.*formatA.*1.0 KiB",
				"manifestB.*maintainerB.*formatB.*3.0 KiB",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore := &internal.TestStore{
				Manifests: tt.manifests,
				Configs:   tt.configs,
				Index:     tt.index,
			}
			summary, err := listModels(testStore)
			if tt.expectErrRegexp != "" {
				// Should be error
				assert.Empty(t, summary, "Should not output summary on error")
				if assert.Error(t, err, "Should return an error") {
					return
				}
				assert.Regexp(t, tt.expectErrRegexp, err.Error())
			} else {
				if !assert.NoError(t, err, "Should not return an error") {
					return
				}
				for _, line := range tt.expectedOutputRegexps {
					// Assert all lines in expected output are somewhere in the summary
					assert.Regexp(t, line, summary)
				}
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input  int64
		output string
	}{
		{input: 0, output: "0 B"},
		{input: 500, output: "500 B"},
		{input: 1<<10 - 1, output: "1023 B"},
		{input: 1 << 10, output: "1.0 KiB"},
		{input: 4.5 * (1 << 10), output: "4.5 KiB"},
		{input: 1<<20 - 1, output: "1023.9 KiB"},
		{input: 1 << 20, output: "1.0 MiB"},
		{input: 6.5 * (1 << 20), output: "6.5 MiB"},
		{input: 1<<30 - 1, output: "1023.9 MiB"},
		{input: 1 << 30, output: "1.0 GiB"},
		{input: 1 << 40, output: "1.0 TiB"},
		{input: 1 << 50, output: "1.0 PiB"},
		{input: 500 * (1 << 50), output: "500.0 PiB"},
		{input: 1 << 60, output: "1024.0 PiB"},
	}
	for idx, tt := range tests {
		t.Run(fmt.Sprintf("test %d", idx), func(t *testing.T) {
			output := formatBytes(tt.input)
			assert.Equalf(t, tt.output, output, "Should convert %d to %s", tt.input, tt.output)
		})
	}
}

func Manifest(configDigest string, layerDigests ...string) ocispec.Manifest {
	manifest := ocispec.Manifest{
		Config: ocispec.Descriptor{
			MediaType: constants.ModelConfigMediaType,
			Digest:    digest.Digest(configDigest),
		},
	}
	for _, layerDigest := range layerDigests {
		manifest.Layers = append(manifest.Layers, ocispec.Descriptor{
			MediaType: constants.ModelLayerMediaType,
			Digest:    digest.Digest(layerDigest),
			Size:      1024,
		})
	}

	return manifest
}

func Config(maintainer, format string) artifact.JozuFile {
	config := artifact.JozuFile{
		Maintainer:  maintainer,
		ModelFormat: format,
	}

	return config
}

func ManifestDesc(digestStr string, valid bool) ocispec.Descriptor {
	size := internal.ValidSize
	if !valid {
		size = internal.InvalidSize
	}

	return ocispec.Descriptor{
		Digest:    digest.Digest(digestStr),
		MediaType: ocispec.MediaTypeImageManifest,
		Size:      size,
	}
}
