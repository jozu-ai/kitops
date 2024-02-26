package constants

import "path/filepath"

type ConfigKey struct{}

const (
	DefaultKitFileName  = "Kitfile"
	DefaultConfigSubdir = ".kitops"

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

func StoragePath(configBase string) string {
	return filepath.Join(configBase, StorageSubpath)
}

func CredentialsPath(configBase string) string {
	return filepath.Join(configBase, CredentialsSubpath)
}
