package constants

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ConfigKey struct{}

const (
	DefaultKitFileName  = "Kitfile"
	DefaultConfigSubdir = "kitops"

	StorageSubpath     = "storage"
	CredentialsSubpath = "credentials.json"

	// Media type for the model layer
	ModelLayerMediaType = "application/vnd.kitops.modelkit.model.v1.tar+gzip"
	// Media type for the dataset layer
	DataSetLayerMediaType = "application/vnd.kitops.modelkit.dataset.v1.tar+gzip"
	// Media type for the code layer
	CodeLayerMediaType = "application/vnd.kitops.modelkit.code.v1.tar+gzip"
	// Media type for the model config (Kitfile)
	ModelConfigMediaType = "application/vnd.kitops.modelkit.config.v1+json"
)

func DefaultConfigPath() (string, error) {
	switch runtime.GOOS {
	case "linux":
		datahome := os.Getenv("XDG_DATA_HOME")
		if datahome == "" {
			// Use default ~/.local/share/
			userhome := os.Getenv("HOME")
			if userhome == "" {
				return "", fmt.Errorf("could not get $HOME directory")
			}
			datahome = filepath.Join(userhome, ".local", "share")
		}
		return filepath.Join(datahome, DefaultConfigSubdir), nil

	case "darwin":
		// Use ~/Library/Caches/kitops
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(cacheDir, DefaultConfigSubdir), nil

	case "windows":
		// Use %LOCALAPPDATA%\kitops
		appdata, err := os.UserCacheDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(appdata, DefaultConfigSubdir), nil

	default:
		return "", fmt.Errorf("Unrecognized operating system")
	}
}

func StoragePath(configBase string) string {
	return filepath.Join(configBase, StorageSubpath)
}

func CredentialsPath(configBase string) string {
	return filepath.Join(configBase, CredentialsSubpath)
}

func IndexJsonPath(configBase string) string {
	return filepath.Join(configBase, "index.json")
}
