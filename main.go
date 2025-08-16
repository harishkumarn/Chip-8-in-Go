package main

import (
	"sync"

	"github.com/harishkumarn/Chip-8-in-Go/engine"
)

func main() {

	wg := sync.WaitGroup{}
	wg.Add(1)
	engine.StartEmulation(RomPath, &wg)
	wg.Wait()
}
