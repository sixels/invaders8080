package main

import (
	"fmt"
	"os"

	"github.com/vnteles/i8080/cpu"
	"github.com/vnteles/i8080/invaders"
)

func loadFileIntoCPUMemory(cpu *cpu.CPU, path string, offset uint16) {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read the file %s\n", path))
	}

	file.Read(cpu.Memory[offset:])
}

// func main() {
// 	cpu := cpu.New()
// 	loadFileIntoCPUMemory(&cpu, "test/invaders.h", 0x0000)
// 	loadFileIntoCPUMemory(&cpu, "test/invaders.g", 0x0800)
// 	loadFileIntoCPUMemory(&cpu, "test/invaders.f", 0x1000)
// 	loadFileIntoCPUMemory(&cpu, "test/invaders.e", 0x1800)

// 	lastInterrupt := time.Now()

// 	for {
// 		cpu.Step(true)

// 		if time.Now().Sub(lastInterrupt).Seconds() > 1.0/60.0 {
// 			if cpu.IntEnable != 0 {
// 				// cpu.Interrupt(2)
// 				lastInterrupt = time.Now()
// 			}
// 		}
// 	}
// }

func main() {
	game := invaders.New()
	game.Run()
}
