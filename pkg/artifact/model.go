package artifact

import "kitops/pkg/lib/constants"

type Model struct {
	Repository string
	Layers     []ModelLayer
	Config     *KitFile
}

type ModelLayer struct {
	Path      string
	MediaType string
}

func (l *ModelLayer) Type() string {
	switch l.MediaType {
	case constants.CodeLayerMediaType:
		return "code"
	case constants.DataSetLayerMediaType:
		return "dataset"
	case constants.ModelConfigMediaType:
		return "config"
	case constants.ModelLayerMediaType:
		return "model"
	}
	return "<unknown>"
}
