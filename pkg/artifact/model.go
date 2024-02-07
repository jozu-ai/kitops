package artifact

type Model struct {
	Repository string
	Tag        string
	Layers     []*ModelLayer
	Config     *JozuFile
}
