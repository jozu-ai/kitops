package artifact

type Model struct {
	Repository string
	Layers     []ModelLayer
	Config     *JozuFile
}
