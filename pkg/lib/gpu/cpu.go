package gpu

import (
	"kitops/pkg/output"

	"golang.org/x/sys/cpu"
)

func getCPUVariant() string {
	if cpu.X86.HasAVX2 {
		output.Debugln("CPU has AVX2")
		return "avx2"
	}
	if cpu.X86.HasAVX {
		output.Debugln("CPU has AVX")
		return "avx"
	}
	output.Debugln("CPU does not have vector extensions")
	return ""
}
