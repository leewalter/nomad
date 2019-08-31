package fingerprint

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/shirou/gopsutil/mem"
)

const bytesInMB int64 = 1024 * 1024

// MemoryFingerprint is used to fingerprint the available memory on the node
type MemoryFingerprint struct {
	StaticFingerprinter
	logger log.Logger
}

// NewMemoryFingerprint is used to create a Memory fingerprint
func NewMemoryFingerprint(logger log.Logger) Fingerprint {
	f := &MemoryFingerprint{
		logger: logger.Named("memory"),
	}
	return f
}

func (f *MemoryFingerprint) Fingerprint(req *FingerprintRequest, resp *FingerprintResponse) error {
	var totalMemory int64
	cfg := req.Config
	if cfg.MemoryMB != 0 {
		totalMemory = int64(cfg.MemoryMB) * bytesInMB
	} else {
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			f.logger.Warn("error reading memory information", "error", err)
			return err
		}
		if memInfo.Total > 0 {
			totalMemory = int64(memInfo.Total)
		}
	}

	if totalMemory > 0 {
		resp.AddAttribute("memory.totalbytes", fmt.Sprintf("%d", totalMemory))

		resp.NodeResources = &structs.NodeResources{
			Memory: structs.NodeMemoryResources{
				MemoryMB: totalMemory / bytesInMB,
			},
		}
	}

	return nil
}
