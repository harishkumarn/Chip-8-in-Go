package engine

import (
	"sync"

	"github.com/harishkumarn/Chip-8-in-Go/engine/processor"
	"github.com/harishkumarn/Chip-8-in-Go/util"
)

func StartEmulation(romPath string, wg *sync.WaitGroup) {
	defer wg.Done() // Needed ?
	proc := processor.Processor{}
	util.InitMemory(romPath, &proc.Memory)
	proc.Init()
	proc.Run()
}
