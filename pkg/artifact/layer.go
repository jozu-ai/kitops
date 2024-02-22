package artifact

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ModelLayer struct {
	BaseDir   string
	MediaType string
}

func (layer *ModelLayer) Apply(writers ...io.Writer) error {
	// Check if path exists
	_, err := os.Stat(layer.BaseDir)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk the context dir and tar everything
	err = filepath.Walk(layer.BaseDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip anything that's not a regular file or directory
		if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
			return nil
		}
		// Skip the baseDir itself
		if file == layer.BaseDir {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// We want the path in the tarball to be relative to the layer's base directory
		subPath := strings.TrimPrefix(strings.Replace(file, layer.BaseDir, "", -1), string(filepath.Separator))
		header.Name = subPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if fi.Mode().IsRegular() {
			err := writeFileToTar(file, tw)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func writeFileToTar(file string, tw *tar.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(tw, f); err != nil {
		return err
	}
	return nil
}
