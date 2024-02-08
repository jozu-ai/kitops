package artifact

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type ModelLayer struct {
	ContextDir string
	Descriptor ocispec.Descriptor
}

func NewLayer(context string) *ModelLayer {
	return &ModelLayer{
		ContextDir: context,
	}
}

func (layer *ModelLayer) Apply(writers ...io.Writer) error {
	// Check if ContextDir exists and is a directory
	fileInfo, err := os.Stat(layer.ContextDir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("ContextDir is not a valid directory")
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk the context dir and tar everything
	err = filepath.Walk(layer.ContextDir, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		parentDir := filepath.Dir(layer.ContextDir)

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(
			strings.Replace(file, parentDir, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
