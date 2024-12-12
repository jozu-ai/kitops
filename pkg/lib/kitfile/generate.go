// Copyright 2024 The KitOps Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package kitfile

import (
	"fmt"
	"io/fs"
	"kitops/pkg/artifact"
	"kitops/pkg/output"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/licensecheck"
)

var modelWeightsSuffixes = []string{
	".safetensors", ".pkl",
	// Pytorch suffixes
	".bin", ".pth", ".pt", ".mar", ".pt2", ".ptl",
	// Tensorflow
	".pb", ".ckpt", ".tflite", ".tfrecords",
	// NumPy
	".npy", ".npz",
	// Keras and others
	".keras", ".h5", ".caffemodel", ".pmml", ".coreml",
	// Other suffixes
	".gguf", ".ggml", ".ggmf", ".llamafile", ".onnx",
}

var docsSuffixes = []string{
	".md", ".adoc", ".html", ".pdf",
}

var metadataSuffixes = []string{
	".json", ".yaml", ".xml", ".csv", ".txt",
}

var datasetSuffixes = []string{
	".tar", ".zip",
}

// Generate a basic Kitfile by looking at the contents of a directory. Parameter
// packageOpt can be used to define metadata for the Kitfile (i.e. the package
// section), which is left empty if the parameter is nil.
func GenerateKitfile(baseDir string, packageOpt *artifact.Package) (*artifact.KitFile, error) {
	kitfile := &artifact.KitFile{
		ManifestVersion: "1.0.0",
	}
	if packageOpt != nil {
		kitfile.Package = *packageOpt
	}

	ds, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	// We can make sure all files are included by including a layer with path '.'
	// However, we only want to do this if it is necessary
	includeCatchallSection := false
	// Dirs we don't know how to handle automatically.
	var unprocessedDirPaths []string
	// Metadata files; we want these to be either model parts (if there is a model)
	// or datasets
	var metadataPaths []string
	var modelFiles []fs.DirEntry
	var detectedLicenseType string
	for _, d := range ds {
		name := d.Name()
		if d.IsDir() {
			err := addDirToKitfile(kitfile, name, d)
			if err != nil {
				unprocessedDirPaths = append(unprocessedDirPaths, name)
			}
			continue
		}

		// Check for "special" files (e.g. readme, license)
		if strings.HasPrefix(strings.ToLower(name), "readme") {
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{
				Path:        name,
				Description: "Readme file",
			})
			continue
		} else if strings.ToLower(name) == "license" {
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{
				Path:        name,
				Description: "License file",
			})
			licenseType, err := detectLicense(filepath.Join(baseDir, name))
			if err != nil {
				output.Debugf("Error determining license type: %s", err)
				output.Logf(output.LogLevelWarn, "Unable to determine license type")
			}
			detectedLicenseType = licenseType
			continue
		}

		// Try to determine type based on file extension
		// To support multi-part models, we need to collect all paths and decide
		// which one is the model and which one(s) are parts
		if anySuffix(name, modelWeightsSuffixes) {
			modelFiles = append(modelFiles, d)
			continue
		}
		// Metadata should be included in either Model or Datasets, depending on
		// other contents
		if anySuffix(name, metadataSuffixes) {
			metadataPaths = append(metadataPaths, name)
			continue
		}
		if anySuffix(name, docsSuffixes) {
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{Path: name})
			continue
		}
		if anySuffix(name, datasetSuffixes) {
			kitfile.DataSets = append(kitfile.DataSets, artifact.DataSet{Path: name})
			continue
		}

		// We don't know what this file is; we'll include it in a catch-all section
		includeCatchallSection = true
	}

	if len(modelFiles) > 0 {
		addModelToKitfile(kitfile, baseDir, modelFiles)
		for _, metadataPath := range metadataPaths {
			kitfile.Model.Parts = append(kitfile.Model.Parts, artifact.ModelPart{Path: metadataPath})
		}
	} else {
		for _, metadataPath := range metadataPaths {
			kitfile.DataSets = append(kitfile.DataSets, artifact.DataSet{Path: metadataPath})
		}
	}

	// Decide how to handle remaining paths. Either package them in one large code layer with basePath
	// or as separate layers for each directory.
	if includeCatchallSection || len(unprocessedDirPaths) > 5 {
		// Overwrite any code layers we added before; this is cleaner than e.g. having a layer for '.' and a layer for 'src'
		kitfile.Code = []artifact.Code{{Path: "."}}
	} else {
		for _, path := range unprocessedDirPaths {
			kitfile.Code = append(kitfile.Code, artifact.Code{Path: path})
		}
	}

	// If we detected a license, try to attach it to the Kitfile section that makes sense
	if kitfile.Model != nil && detectedLicenseType != "" {
		kitfile.Model.License = detectedLicenseType
	} else if len(kitfile.DataSets) == 1 {
		kitfile.DataSets[0].License = detectedLicenseType
	} else if len(kitfile.Code) == 1 {
		kitfile.Code[0].License = detectedLicenseType
	} else {
		kitfile.Package.License = detectedLicenseType
	}

	return kitfile, nil
}

func addDirToKitfile(kitfile *artifact.KitFile, path string, d fs.DirEntry) error {
	// TODO: consider looking into directories to see if we can figure out what they store? Might work for datasets
	switch d.Name() {
	case "docs":
		kitfile.Docs = append(kitfile.Docs, artifact.Docs{
			Path: path,
		})
	case "src", "pkg", "lib", "build":
		kitfile.Code = append(kitfile.Code, artifact.Code{
			Path: path,
		})
	default:
		return fmt.Errorf("could not determine data type for directory")
	}
	return nil
}

func addModelToKitfile(kitfile *artifact.KitFile, baseDir string, modelFiles []fs.DirEntry) error {
	if len(modelFiles) == 0 {
		return nil
	}

	if len(modelFiles) == 1 {
		filename := modelFiles[0].Name()
		kitfile.Model = &artifact.Model{
			Path: filename,
			Name: strings.TrimSuffix(filename, filepath.Ext(filename)),
		}
		return nil
	}

	// We'll handle two cases here: 1) the Model is split into multiple files (e.g. safetensors) or 2) there is a
	// main model plus smaller adaptor(s)
	largestFile := ""
	largestSize := int64(0)
	averageSize := int64(0)
	for _, modelFile := range modelFiles {
		info, err := modelFile.Info()
		if err != nil {
			return fmt.Errorf("failed to process file %s: %w", filepath.Join(baseDir, modelFile.Name()), err)
		}
		size := info.Size()
		if size > largestSize {
			largestSize = size
			largestFile = modelFile.Name()
		}
		averageSize = averageSize + size
	}
	// Integer division is probably fine here; at most we're off by a byte.
	averageSize = averageSize / int64(len(modelFiles))

	// If the biggest file is 1.5x the average, make it the model and the rest parts; otherwise, add
	// all parts in lexical order
	if largestSize > averageSize+(averageSize/2) {
		kitfile.Model = &artifact.Model{
			Path: largestFile,
		}
		kitfile.Model.Name = strings.TrimSuffix(largestFile, filepath.Ext(largestFile))
		for _, modelFile := range modelFiles {
			if modelFile.Name() == largestFile {
				continue
			}
			kitfile.Model.Parts = append(kitfile.Model.Parts, artifact.ModelPart{
				Path: modelFile.Name(),
			})
		}
	} else {
		kitfile.Model = &artifact.Model{
			Path: modelFiles[0].Name(),
		}
		for _, modelFile := range modelFiles[1:] {
			kitfile.Model.Parts = append(kitfile.Model.Parts, artifact.ModelPart{
				Path: modelFile.Name(),
			})
		}
	}
	return nil
}

func detectLicense(licensePath string) (string, error) {
	license, err := os.ReadFile(licensePath)
	if err != nil {
		return "", fmt.Errorf("failed to read license file: %w", err)
	}
	cov := licensecheck.Scan(license)
	if len(cov.Match) == 1 {
		return cov.Match[0].ID, nil
	} else {
		return "", fmt.Errorf("multiple licenses matched license file")
	}
}

func anySuffix(query string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(query, suffix) {
			return true
		}
	}
	return false
}
