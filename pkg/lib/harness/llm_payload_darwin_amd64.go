package harness

import "embed"

//go:embed llama.cpp/build/darwin/x86_64/*/bin/*
var serverEmbed embed.FS
