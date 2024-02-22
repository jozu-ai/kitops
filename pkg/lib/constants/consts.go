package constants

const (
	DefaultModelFileName = "Kitfile"

	// Media type for the model layer
	ModelLayerMediaType = "application/vnd.kitops.modelkit.model.v1.tar+gzip"
	// Media type for the dataset layer
	DataSetLayerMediaType = "application/vnd.kitops.modelkit.dataset.v1.tar+gzip"
	// Media type for the code layer
	CodeLayerMediaType = "application/vnd.kitops.modelkit.code.v1.tar+gzip"
	// Media type for the model config (Kitfile)
	ModelConfigMediaType = "application/vnd.kitops.modelkit.config.v1+json"
)
