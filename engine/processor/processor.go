package processor

import (
	"math/rand/v2"

	"github.com/harishkumarn/Chip-8-in-Go/io"
	"github.com/harishkumarn/Chip-8-in-Go/util"
)

type Processor struct {
	Memory                                 [1 << 12]uint8
	display                                *io.Display
	Keypad                                 *io.Keypad
	registers                              [16]uint8
	stack                                  [16]uint16
	pc, index                              uint16
	stackPointer, delay, sound, pressedKey uint8
	blockingKeyPress                       chan uint8
	blockingIO                             bool
}

func (cpu *Processor) Init() {
	cpu.pc = 0x200
	cpu.pressedKey = 16
	go func() {
		kpChan := cpu.Keypad.GetKeyPress()
		for {
			key := <-kpChan
			if cpu.blockingIO {
				cpu.blockingKeyPress <- key
				cpu.blockingIO = false
			} else {
				cpu.pressedKey = key
			}
		}
	}()
}

func (cpu *Processor) getInstruction() (uint8, uint8) {
	high := cpu.Memory[cpu.pc]
	low := cpu.Memory[cpu.pc+1]
	cpu.pc += 2
	return high, low
}

func (cpu *Processor) stackPush(addr uint16) {
	cpu.stack[cpu.stackPointer] = addr
	cpu.stackPointer += 1
}

func (cpu *Processor) stackPop() uint16 {
	cpu.stackPointer -= 1
	return cpu.stack[cpu.stackPointer+1]
}

func parseAddress(high, low uint8) uint16 {
	var addr uint16 = uint16(high)
	addr <<= 8
	addr += uint16(low)
	return addr
}

func (cpu *Processor) Run() {
	for {
		if cpu.delay > 0 {
			cpu.delay -= 1
		}
		if cpu.sound > 0 {
			cpu.sound -= 1
		}
		high, low := cpu.getInstruction()
		if low&0xF0 == 0xE0 {
			cpu.display.CLS()
		} else if low == 0xEE {
			cpu.pc = cpu.stackPop()
		} else if high&0xF0 == 0x10 || high&0xF == 0x20 {
			addr := parseAddress(high, low)
			if high&0xF0 == 0x20 { // CALL
				cpu.stackPush(cpu.pc)
			}
			cpu.pc = addr
		} else if high&0xF0 == 0x30 {
			if cpu.registers[high&0x0F] == low {
				cpu.pc += 2
			}
		} else if high&0xF0 == 0x40 {
			if cpu.registers[high&0x0F] != low {
				cpu.pc += 2
			}
		} else if high&0xF0 == 0x50 {
			if cpu.registers[high&0x0F] == cpu.registers[(low>>4)] {
				cpu.pc += 2
			}
		} else if high&0xF0 == 0x60 {
			cpu.registers[high&0x0F] = low
		} else if high&0xF0 == 0x70 {
			cpu.registers[high&0x0F] += low
		} else if high&0xF0 == 0x80 {
			switch low & 0x0F {
			case 0:
				cpu.registers[high&0x0F] = cpu.registers[low>>4]
			case 0x1:
				cpu.registers[high&0x0F] |= cpu.registers[low>>4]
			case 0x2:
				cpu.registers[high&0x0F] &= cpu.registers[low>>4]
			case 0x3:
				cpu.registers[high&0x0F] ^= cpu.registers[low>>4]
			case 0x4:
				var sum uint16 = uint16(cpu.registers[high&0x0F]) + uint16(cpu.registers[low>>4])
				cpu.registers[high&0x0F] = uint8(sum & 0xFF)
				if sum > 0xFF {
					cpu.registers[0xF] = 1
				} else {
					cpu.registers[0xF] = 0
				}
			case 0x5:
				var diff int16 = int16(cpu.registers[high&0x0F]) - int16(cpu.registers[low>>4])
				cpu.registers[high&0x0F] = uint8(diff & 0xFF)
				if diff < 0 {
					cpu.registers[0xF] = 1
				} else {
					cpu.registers[0xF] = 0
				}
			case 0x6:
				var res uint8 = cpu.registers[low>>4] >> 1
				cpu.registers[high&0x0F] = res
				if cpu.registers[low>>4]&0x1 == 1 {
					cpu.registers[0xF] = 1
				} else {
					cpu.registers[0x0F] = 0
				}
			case 0x7:
				var diff int16 = int16(cpu.registers[low>>4]) - int16(cpu.registers[high&0x0F])
				cpu.registers[low>>4] = uint8(diff & 0xFF)
				if diff < 0 {
					cpu.registers[0xF] = 1
				} else {
					cpu.registers[0xF] = 0
				}
			case 0xE:
				var res uint8 = cpu.registers[low>>4] << 1
				cpu.registers[high&0x0F] = res
				if cpu.registers[low>>4]&0x80 == 0x80 {
					cpu.registers[0xF] = 1
				} else {
					cpu.registers[0x0F] = 0
				}
			}
		} else if high&0xF0 == 0x90 {
			if cpu.registers[high&0x0F] != cpu.registers[(low>>4)] {
				cpu.pc += 2
			}
		} else if high&0xF0 == 0xA0 {
			cpu.index = parseAddress(high, low)
		} else if high&0xF0 == 0xB0 {
			cpu.pc = parseAddress(high, low) + uint16(cpu.registers[0])
		} else if high&0xF0 == 0xC0 {
			cpu.registers[high&0xF0] = uint8(rand.IntN(0xFFFF)) & low
		} else if high&0xF0 == 0xD0 {
			cpu.display.Draw(high&0xF0, low>>4, low&0x0F)
		} else if high&0xF0 == 0xE0 {
			if low == 0x9E && cpu.registers[high&0x0F] == cpu.pressedKey {
				cpu.pc += 2
			} else if low == 0xA1 && cpu.registers[high&0x0F] != cpu.pressedKey {
				cpu.pc += 2
			}
			cpu.pressedKey = 16
		} else if high&0xF0 == 0xF0 {
			switch low {
			case 0x07:
				cpu.registers[high&0x0F] = cpu.delay
			case 0x0A:
				cpu.blockingIO = true
				cpu.registers[high&0x0F] = <-cpu.blockingKeyPress
			case 0x15:
				cpu.delay = cpu.registers[high&0x0F]
			case 0x18:
				cpu.sound = cpu.registers[high&0x0F]
			case 0x1E:
				cpu.index += uint16(cpu.registers[high&0x0F])
			case 0x29:
				cpu.index = util.SpriteAddress(high & 0x0F)
			case 0x33:
				val := cpu.registers[high&0x0F]
				cpu.Memory[cpu.index+2] = val % 10
				val /= 10
				cpu.Memory[cpu.index+1] = val % 10
				val /= 10
				cpu.Memory[cpu.index] = val
			case 0x55:
				var i uint16
				var x uint16 = uint16(high & 0x0F)
				for i = 0; i < x; i += 1 {
					cpu.Memory[cpu.index+i] = cpu.registers[i]
				}
			case 0x65:
				var i uint16
				var x uint16 = uint16(high & 0x0F)
				for i = 0; i < x; i += 1 {
					cpu.registers[i] = cpu.Memory[cpu.index+i]
				}
			}
		}
	}
}
