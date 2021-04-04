package cpu

// Flags
const (
	Z = 1 << iota
	S
	P
	CY
	AC
)

// SetZSP flags based on result
func (cpu *CPU) SetZSP(result uint8, size uint8) {
	// Zero flag
	cpu.Flags.Set(Z, result == 0)

	// Sign flag
	cpu.Flags.Set(S, (result&(1<<(size-1)))>>(size-1) == 1)

	// Parity flag
	// parity := bits.ByteParity(result)
	cpu.Flags.Set(P, (result&1) == 0)
}
