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

type fileType int

const (
	fileTypeModel fileType = iota
	fileTypeDataset
	fileTypeCode
	fileTypeDocs
	fileTypeMetadata
	fileTypeUnknown
)

var modelWeightsSuffixes = []string{
	".safetensors", ".pkl", ".joblib",
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
	".json", ".yaml", ".xml", ".txt",
}

var datasetSuffixes = []string{
	".tar", ".zip", ".parquet", ".csv",
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
	var modelFiles, metadataPaths []string
	var detectedLicenseType string
	for _, d := range ds {
		filename := d.Name()
		if d.IsDir() {
			dirModelFiles, err := addDirToKitfile(kitfile, filename, d)
			if err != nil {
				unprocessedDirPaths = append(unprocessedDirPaths, filename)
			}
			modelFiles = append(modelFiles, dirModelFiles...)
			continue
		}

		// Check for "special" files (e.g. readme, license)
		if strings.HasPrefix(strings.ToLower(filename), "readme") {
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{
				Path:        filename,
				Description: "Readme file",
			})
			continue
		} else if strings.HasPrefix(strings.ToLower(filename), "license") {
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{
				Path:        filename,
				Description: "License file",
			})
			licenseType, err := detectLicense(filepath.Join(baseDir, filename))
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
		switch determineFileType(filename) {
		case fileTypeModel:
			modelFiles = append(modelFiles, filename)
		case fileTypeMetadata:
			// Metadata should be included in either Model or Datasets, depending on
			// other contents
			metadataPaths = append(metadataPaths, filename)
		case fileTypeDocs:
			kitfile.Docs = append(kitfile.Docs, artifact.Docs{Path: filename})
		case fileTypeDataset:
			kitfile.DataSets = append(kitfile.DataSets, artifact.DataSet{Path: filename})
		default:
			// File is either code or unknown; we'll have to include it in a catch-all section
			includeCatchallSection = true
		}
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

func addDirToKitfile(kitfile *artifact.KitFile, path string, d fs.DirEntry) (modelFiles []string, err error) {
	switch d.Name() {
	case "docs":
		kitfile.Docs = append(kitfile.Docs, artifact.Docs{
			Path: path,
		})
		return nil, nil
	case "src", "pkg", "lib", "build":
		kitfile.Code = append(kitfile.Code, artifact.Code{
			Path: path,
		})
		return nil, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", path, err)
	}

	// Sort entries in the directory to try and figure out what it contains. We'll reuse the
	// fact that the fileTypes are enumerated using iota (and so are ints) to index correctly.
	// Avoid using maps here since they iterate in a random order.
	directoryContents := [int(fileTypeUnknown) + 1][]string{}
	for _, entry := range entries {
		relPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			// TODO: we can potentially recurse further here if we find we need to
			directoryContents[int(fileTypeUnknown)] = append(directoryContents[int(fileTypeUnknown)], relPath)
			continue
		}
		fileType := determineFileType(entry.Name())
		if fileType == fileTypeModel {
			modelFiles = append(modelFiles, relPath)
		}
		directoryContents[int(fileType)] = append(directoryContents[int(fileType)], relPath)
	}

	// Try to detect directories that contain e.g. only datasets so we can add them
	overallFiletype := fileTypeUnknown
	directoryHasMixedContents := false
	for fType, files := range directoryContents {
		if len(files) > 0 && fileType(fType) != fileTypeMetadata {
			if overallFiletype != fileTypeUnknown {
				directoryHasMixedContents = true
			}
			overallFiletype = fileType(fType)
		}
	}
	if directoryHasMixedContents {
		return modelFiles, fmt.Errorf("mixed content in directory; unable to determine type")
	}
	switch overallFiletype {
	case fileTypeModel:
		// Include any metadata files as modelParts later
		modelFiles = append(modelFiles, directoryContents[int(fileTypeMetadata)]...)
	case fileTypeDataset:
		kitfile.DataSets = append(kitfile.DataSets, artifact.DataSet{Path: path})
	case fileTypeDocs:
		kitfile.Docs = append(kitfile.Docs, artifact.Docs{Path: path})
	default:
		// If it's overall code, metadata, or unknown, just return it as unprocessed and let it be added as a Code section
		// later
		return modelFiles, fmt.Errorf("directory should be handled as Code")
	}

	return modelFiles, nil
}

func determineFileType(filename string) fileType {
	if anySuffix(filename, modelWeightsSuffixes) {
		return fileTypeModel
	}
	// Metadata should be included in either Model or Datasets, depending on
	// other contents
	if anySuffix(filename, metadataSuffixes) {
		return fileTypeMetadata
	}
	if anySuffix(filename, docsSuffixes) {
		return fileTypeDocs
	}
	if anySuffix(filename, datasetSuffixes) {
		return fileTypeDataset
	}
	return fileTypeUnknown

}

func addModelToKitfile(kitfile *artifact.KitFile, baseDir string, modelPaths []string) error {
	if len(modelPaths) == 0 {
		return nil
	}

	if len(modelPaths) == 1 {
		filename := filepath.Base(modelPaths[0])
		kitfile.Model = &artifact.Model{
			Path: modelPaths[0],
			Name: strings.TrimSuffix(filename, filepath.Ext(filename)),
		}
		return nil
	}

	// We'll handle two cases here: 1) the Model is split into multiple files (e.g. safetensors) or 2) there is a
	// main model plus smaller adaptor(s)
	largestFile := ""
	largestSize := int64(0)
	averageSize := int64(0)
	for _, modelFile := range modelPaths {
		info, err := os.Stat(modelFile)
		if err != nil {
			return fmt.Errorf("failed to process file %s: %w", filepath.Join(baseDir, modelFile), err)
		}
		size := info.Size()
		if size > largestSize {
			largestSize = size
			largestFile = modelFile
		}
		averageSize = averageSize + size
	}
	// Integer division is probably fine here; at most we're off by a byte.
	averageSize = averageSize / int64(len(modelPaths))

	// If the biggest file is 1.5x the average, make it the model and the rest parts; otherwise, add
	// all parts in lexical order
	if largestSize > averageSize+(averageSize/2) {
		kitfile.Model = &artifact.Model{
			Path: largestFile,
		}
		kitfile.Model.Name = strings.TrimSuffix(largestFile, filepath.Ext(largestFile))
		for _, modelFile := range modelPaths {
			if modelFile == largestFile {
				continue
			}
			kitfile.Model.Parts = append(kitfile.Model.Parts, artifact.ModelPart{
				Path: modelFile,
			})
		}
	} else {
		kitfile.Model = &artifact.Model{
			Path: modelPaths[0],
		}
		for _, modelFile := range modelPaths[1:] {
			kitfile.Model.Parts = append(kitfile.Model.Parts, artifact.ModelPart{
				Path: modelFile,
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
