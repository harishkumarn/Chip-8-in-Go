package engine

import (
	"fmt"
	"sync"

	"github.com/harishkumarn/Chip-8-in-Go/engine/processor"
	"github.com/harishkumarn/Chip-8-in-Go/util"
)

func StartEmulation(romPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	proc := processor.Processor{}
	util.InitMemory(romPath, &proc.Memory)
	proc.Init()
	fmt.Println("Starting the processor")
	proc.Run()
}
