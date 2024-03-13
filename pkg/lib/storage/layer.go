package storage

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"kitops/pkg/artifact"
	"kitops/pkg/output"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// compressLayer compresses an *artifact.ModelLayer to a gzipped tar file. In order to return
// a descriptor (including hash) for the compressed file, the layer is saved to a temporary file
// on disk and must be moved to an appropriate location. It is the responsibility of the caller
// to clean up the temporary file when it is no longer needed.
func compressLayer(layer *artifact.ModelLayer) (tempFilePath string, desc ocispec.Descriptor, err error) {
	pathInfo, err := os.Stat(layer.Path)
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, err
	}
	tempFile, err := os.CreateTemp("", "kitops_layer_*")
	if err != nil {
		return "", ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempFileName := tempFile.Name()
	output.Debugf("Compressing layer to temporary file %s", tempFileName)

	digester := digest.Canonical.Digester()
	mw := io.MultiWriter(tempFile, digester.Hash())

	// Note: we have to close gzip writer before reading digest from digester as closing is what writes the GZIP footer
	gzw := gzip.NewWriter(mw)
	tw := tar.NewWriter(gzw)

	// Wrapper function for closing writers before returning an error
	handleErr := func(err error) (string, ocispec.Descriptor, error) {
		// Don't care about these errors since we'll be deleting the file anyways
		_ = tw.Close()
		_ = gzw.Close()
		_ = tempFile.Close()
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, err
	}

	if pathInfo.Mode().IsRegular() {
		if err := writeHeaderToTar(pathInfo.Name(), pathInfo, tw); err != nil {
			return handleErr(err)
		}
		if err := writeFileToTar(layer.Path, pathInfo, tw); err != nil {
			return handleErr(err)
		}
	} else if pathInfo.IsDir() {
		if err := writeDirToTar(layer.Path, tw); err != nil {
			return handleErr(err)
		}
	} else {
		return handleErr(fmt.Errorf("path %s is neither a file nor a directory", layer.Path))
	}

	callAndPrintError(tw.Close, "Failed to close tar writer: %s")
	callAndPrintError(gzw.Close, "Failed to close gzip writer: %s")

	tempFileInfo, err := tempFile.Stat()
	if err != nil {
		removeTempFile(tempFileName)
		return "", ocispec.DescriptorEmptyJSON, fmt.Errorf("failed to stat temporary file: %w", err)
	}
	callAndPrintError(tempFile.Close, "Failed to close temporary file: %s")

	desc = ocispec.Descriptor{
		MediaType: layer.MediaType,
		Digest:    digester.Digest(),
		Size:      tempFileInfo.Size(),
	}
	return tempFileName, desc, nil
}

// writeDirToTar walks the filesystem at basePath, compressing contents via the *tar.Writer.
// Any non-regular files and directories (e.g. symlinks) are skipped.
func writeDirToTar(basePath string, tw *tar.Writer) error {
	// We'll want paths in the tarball to be relative to the *parent* of basePath since we want
	// to compress the directory pointed at by basePath
	trimPath := filepath.Dir(basePath)
	return filepath.Walk(basePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip anything that's not a regular file or directory
		if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(strings.Replace(file, trimPath, "", -1), string(filepath.Separator))
		if relPath == "" {
			relPath = filepath.Base(basePath)
		}
		if err := writeHeaderToTar(relPath, fi, tw); err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		return writeFileToTar(file, fi, tw)
	})
}

func writeHeaderToTar(name string, fi os.FileInfo, tw *tar.Writer) error {
	header, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return fmt.Errorf("failed to generate header for %s: %w", name, err)
	}
	header.Name = name
	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	output.Debugf("Wrote header %s to tar file", header.Name)
	return nil
}

func writeFileToTar(file string, fi os.FileInfo, tw *tar.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file for archiving: %w", err)
	}
	defer f.Close()

	if written, err := io.Copy(tw, f); err != nil {
		return fmt.Errorf("failed to add file to archive: %w", err)
	} else if written != fi.Size() {
		return fmt.Errorf("error writing file: %w", err)
	}
	output.Debugf("Wrote file %s to tar file", file)
	return nil
}

// callAndPrintError is a wrapper to print an error message for a function that
// may return an error. The error is printed and then discarded.
func callAndPrintError(f func() error, msg string) {
	if err := f(); err != nil {
		output.Errorf(msg, err)
	}
}

func removeTempFile(filepath string) {
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		output.Errorf("Failed to clean up temporary file %s: %s", filepath, err)
	}
}
