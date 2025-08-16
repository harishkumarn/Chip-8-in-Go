package io

import (
	"fmt"

	"github.com/harishkumarn/Chip-8-in-Go/util"
)

type Keypad struct{}

func (kp *Keypad) GetKeyPress() <-chan uint8 {
	keyPress := make(chan uint8)

	go func() {
		var key rune
		for {
			fmt.Scanf("%c", &key)
			keyPress <- util.GetMappedKey(key)
		}
	}()
	return keyPress
}
