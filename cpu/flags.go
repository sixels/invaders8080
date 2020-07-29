package cpu

import "github.com/protoshark/invaders8080/bits"

// Flags
const (
	Z = 1 << iota
	S
	P
	CY
	AC
)

// SetZSP flags based on result
func (cpu *CPU) SetZSP(result uint8) {
	// Zero flag
	cpu.Flags.Set(Z, result == 0)

	// Sign flag
	cpu.Flags.Set(S, (result&0x80) == 0x80)

	// Parity flag
	parity := bits.ByteParity(result)
	cpu.Flags.Set(P, parity)
}
