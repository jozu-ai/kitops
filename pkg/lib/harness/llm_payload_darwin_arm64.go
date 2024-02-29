package harness

import "embed"

//go:embed llama.cpp/build/darwin/arm64/*/bin/*
var serverEmbed embed.FS