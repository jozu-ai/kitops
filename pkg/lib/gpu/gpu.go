package gpu

import (
	"runtime"
)

type GPUInfo struct {
	Library string
	Variant string
}

func GetGPUInfo() *GPUInfo {

	if runtime.GOARCH == "amd64" {
		return &GPUInfo{
			Library: "cpu",
			Variant: getCPUVariant(),
		}
	}
	return &GPUInfo{
		Library: "metal",
		
	}
}
