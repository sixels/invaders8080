package cpu

import (
	"fmt"

	"github.com/protoshark/invaders8080/bits"
)

// Instruction data -> the name and the number of cycles each instruction run
type Instruction struct {
	Name   string
	Cycles uint8
	Size   uint8
}

// InstructionTable for disassembling
var InstructionTable = []Instruction{
	/*0x00*/ {"NOP", 0x04, 1}, {"LXI B,$", 0x0a, 3}, {"STAX B", 0x07, 1}, {"INX B", 0x05, 1}, {"INR B", 0x05, 1}, {"DCR B", 0x05, 1}, {"MVI B,#$", 0x07, 2}, {"RLC", 0x04, 1},
	/*0x08*/ {"UNK", 0x04, 1}, {"DAD B", 0x0a, 1}, {"LDAX B", 0x07, 1}, {"DCX B", 0x05, 1}, {"INR C", 0x05, 1}, {"DCR C", 0x05, 1}, {"MVI C,#$", 0x07, 2}, {"RRC", 0x04, 1},
	/*0x10*/ {"UNK", 0x04, 1}, {"LXI D,$", 0x0a, 3}, {"STAX D", 0x07, 1}, {"INX D", 0x05, 1}, {"INR D", 0x05, 1}, {"DCR D", 0x05, 1}, {"MVI D,#$", 0x07, 2}, {"RAL", 0x04, 1},
	/*0x18*/ {"UNK", 0x04, 1}, {"DAD D", 0x0a, 1}, {"LDAX D", 0x07, 1}, {"DCX D", 0x05, 1}, {"INR E", 0x05, 1}, {"DCR E", 0x05, 1}, {"MVI E,#$", 0x07, 2}, {"RAR", 0x04, 1},
	/*0x20*/ {"UNK", 0x04, 1}, {"LXI H,$", 0x0a, 3}, {"SHLD", 0x10, 1}, {"INX H", 0x05, 1}, {"INR H", 0x05, 1}, {"DCR H", 0x05, 1}, {"MVI H,#$", 0x07, 2}, {"DAA", 0x04, 1},
	/*0x28*/ {"UNK", 0x04, 1}, {"DAD H", 0x0a, 1}, {"LHLD $", 0x10, 3}, {"DCX H", 0x05, 1}, {"INR L", 0x05, 1}, {"DCR L", 0x05, 1}, {"MVI L,#$", 0x07, 2}, {"CMA", 0x04, 1},
	/*0x30*/ {"UNK", 0x04, 1}, {"LXI SP,$", 0x0a, 3}, {"STA $", 0x0d, 3}, {"INX SP", 0x05, 1}, {"INR M", 0x0a, 1}, {"DCR M", 0x0a, 1}, {"MVI M,#$", 0x0a, 2}, {"STC", 0x04, 1},
	/*0x38*/ {"UNK", 0x04, 1}, {"DAD SP", 0x0a, 1}, {"LDA $", 0x0d, 3}, {"DCX SP", 0x05, 1}, {"INR A", 0x05, 1}, {"DCR A", 0x05, 1}, {"MVI A,#$", 0x07, 2}, {"CMC", 0x04, 1},

	/*0x40*/ {"MOV B,B", 0x05, 1}, {"MOV B,C", 0x05, 1}, {"MOV B,D", 0x05, 1}, {"MOV B,E", 0x05, 1}, {"MOV B,H", 0x05, 1}, {"MOV B,L", 0x05, 1}, {"MOV B,M", 0x07, 1}, {"MOV B,A", 0x05, 1},
	/*0x48*/ {"MOV C,B", 0x05, 1}, {"MOV C,C", 0x05, 1}, {"MOV C,D", 0x05, 1}, {"MOV C,E", 0x05, 1}, {"MOV C,H", 0x05, 1}, {"MOV C,L", 0x05, 1}, {"MOV C,M", 0x07, 1}, {"MOV C,A", 0x05, 1},
	/*0x50*/ {"MOV D,B", 0x05, 1}, {"MOV D,C", 0x05, 1}, {"MOV D,D", 0x05, 1}, {"MOV D,E", 0x05, 1}, {"MOV D,H", 0x05, 1}, {"MOV D,L", 0x05, 1}, {"MOV D,M", 0x07, 1}, {"MOV D,A", 0x05, 1},
	/*0x58*/ {"MOV E,B", 0x05, 1}, {"MOV E,C", 0x05, 1}, {"MOV E,D", 0x05, 1}, {"MOV E,E", 0x05, 1}, {"MOV E,H", 0x05, 1}, {"MOV E,L", 0x05, 1}, {"MOV E,M", 0x07, 1}, {"MOV E,A", 0x05, 1},
	/*0x60*/ {"MOV H,B", 0x05, 1}, {"MOV H,C", 0x05, 1}, {"MOV H,D", 0x05, 1}, {"MOV H,E", 0x05, 1}, {"MOV H,H", 0x05, 1}, {"MOV H,L", 0x05, 1}, {"MOV H,M", 0x07, 1}, {"MOV H,A", 0x05, 1},
	/*0x68*/ {"MOV L,B", 0x05, 1}, {"MOV L,C", 0x05, 1}, {"MOV L,D", 0x05, 1}, {"MOV L,E", 0x05, 1}, {"MOV L,H", 0x05, 1}, {"MOV L,L", 0x05, 1}, {"MOV L,M", 0x07, 1}, {"MOV L,A", 0x05, 1},
	/*0x70*/ {"MOV M,B", 0x07, 1}, {"MOV M,C", 0x07, 1}, {"MOV M,D", 0x07, 1}, {"MOV M,E", 0x07, 1}, {"MOV M,H", 0x07, 1}, {"MOV M,L", 0x07, 1}, {"HLT", 0x07, 1}, {"MOV M,A", 0x07, 1},
	/*0x78*/ {"MOV A,B", 0x05, 1}, {"MOV A,C", 0x05, 1}, {"MOV A,D", 0x05, 1}, {"MOV A,E", 0x05, 1}, {"MOV A,H", 0x05, 1}, {"MOV A,L", 0x05, 1}, {"MOV A,M", 0x07, 1}, {"MOV A,A", 0x05, 1},

	/*0x80*/ {"ADD B", 0x04, 1}, {"ADD C", 0x04, 1}, {"ADD D", 0x04, 1}, {"ADD E", 0x04, 1}, {"ADD H", 0x04, 1}, {"ADD L", 0x04, 1}, {"ADD M", 0x07, 1}, {"ADD A", 0x04, 1},
	/*0x88*/ {"ADC B", 0x04, 1}, {"ADC C", 0x04, 1}, {"ADC D", 0x04, 1}, {"ADC E", 0x04, 1}, {"ADC H", 0x04, 1}, {"ADC L", 0x04, 1}, {"ADC M", 0x07, 1}, {"ADC A", 0x04, 1},
	/*0x90*/ {"SUB B", 0x04, 1}, {"SUB C", 0x04, 1}, {"SUB D", 0x04, 1}, {"SUB E", 0x04, 1}, {"SUB H", 0x04, 1}, {"SUB L", 0x04, 1}, {"SUB M", 0x07, 1}, {"SUB A", 0x04, 1},
	/*0x98*/ {"SBB B", 0x04, 1}, {"SBB C", 0x04, 1}, {"SBB D", 0x04, 1}, {"SBB E", 0x04, 1}, {"SBB H", 0x04, 1}, {"SBB L", 0x04, 1}, {"SBB M", 0x07, 1}, {"SBB A", 0x04, 1},
	/*0xa0*/ {"ANA B", 0x04, 1}, {"ANA C", 0x04, 1}, {"ANA D", 0x04, 1}, {"ANA E", 0x04, 1}, {"ANA H", 0x04, 1}, {"ANA L", 0x04, 1}, {"ANA M", 0x07, 1}, {"ANA A", 0x04, 1},
	/*0xa8*/ {"XRA B", 0x04, 1}, {"XRA C", 0x04, 1}, {"XRA D", 0x04, 1}, {"XRA E", 0x04, 1}, {"XRA H", 0x04, 1}, {"XRA L", 0x04, 1}, {"XRA M", 0x07, 1}, {"XRA A", 0x04, 1},
	/*0xb0*/ {"ORA B", 0x04, 1}, {"ORA C", 0x04, 1}, {"ORA D", 0x04, 1}, {"ORA E", 0x04, 1}, {"ORA H", 0x04, 1}, {"ORA L", 0x04, 1}, {"ORA M", 0x07, 1}, {"ORA A", 0x04, 1},
	/*0xb8*/ {"CMP B", 0x04, 1}, {"CMP C", 0x04, 1}, {"CMP D", 0x04, 1}, {"CMP E", 0x04, 1}, {"CMP H", 0x04, 1}, {"CMP L", 0x04, 1}, {"CMP M", 0x07, 1}, {"CMP A", 0x04, 1},

	/*0xc0*/ {"RNZ", 0x05, 1}, {"POP B", 0x0a, 1}, {"JNZ $", 0x0a, 3}, {"JMP $", 0x0a, 3}, {"CNZ $", 0x0b, 3}, {"PUSH B", 0x0b, 1}, {"ADI #$", 0x07, 2}, {"RST 0", 0x04, 1},
	/*0xc8*/ {"RZ", 0x05, 1}, {"RET", 0x0a, 1}, {"JZ $", 0x0a, 3}, {"UNK", 0x04, 1}, {"CZ $", 0x0b, 3}, {"CALL $", 0x11, 3}, {"ACI #$", 0x07, 2}, {"RST 1", 0x04, 1},
	/*0xd0*/ {"RNC", 0x05, 1}, {"POP D", 0x0a, 1}, {"JNC $", 0x0a, 3}, {"OUT #$", 0x0a, 2}, {"CNC $", 0x0b, 3}, {"PUSH D", 0x0b, 1}, {"SUI #$", 0x07, 2}, {"RST 2", 0x04, 1},
	/*0xd8*/ {"RC", 0x05, 1}, {"UNK", 0x04, 1}, {"JC $", 0x0a, 3}, {"IN #$", 0x0a, 2}, {"CC $", 0x0b, 3}, {"UNK", 0x04, 1}, {"SBI #$", 0x07, 2}, {"RST 3", 0x04, 1},
	/*0xe0*/ {"RPO", 0x05, 1}, {"POP H", 0x0a, 1}, {"JPO $", 0x0a, 3}, {"XTHL", 0x12, 1}, {"CPO $", 0x0b, 3}, {"PUSH H", 0x0b, 1}, {"ANI #$", 0x07, 2}, {"RST 4", 0x04, 1},
	/*0xe8*/ {"RPE", 0x05, 1}, {"PCHL", 0x05, 1}, {"JPE $", 0x0a, 3}, {"XCHG", 0x04, 1}, {"CPE $", 0x0b, 3}, {"UNK", 0x04, 1}, {"XRI #$", 0x07, 2}, {"RST 5", 0x04, 1},
	/*0xf0*/ {"RP", 0x05, 1}, {"POP PSW", 0x0a, 1}, {"JP $", 0x0a, 3}, {"DI", 0x04, 1}, {"CP $", 0x0b, 3}, {"PUSH PSW", 0x0b, 1}, {"ORI #$", 0x07, 2}, {"RST 6", 0x04, 1},
	/*0xf8*/ {"RM", 0x05, 1}, {"SPHL", 0x05, 1}, {"JM $", 0x0a, 3}, {"EI", 0x04, 1}, {"CM $", 0x0b, 3}, {"UNK", 0x04, 1}, {"CPI #$", 0x07, 2}, {"RST 7", 0x04, 1},
}

// Memory information
const (
	// CPU ROM end region
	RomOffset  = 0x2000
	VRAMOffset = 0x2400
)

// The CPU structure
type CPU struct {
	// registers
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8
	// CPU flags
	Flags bits.Bitfield
	// Memory
	Memory []byte
	// Pointers
	PC uint16
	SP uint16

	// Interruption control
	IntEnable bool

	// keep control of cycles
	Cycles uint32
}

// Disassembly a buffer
func Disassembly(buffer []byte, offset uint16) uint8 {
	opcode := buffer[offset]
	instruction := InstructionTable[opcode]

	fmt.Printf("%04x %s", offset, instruction.Name)

	for i := instruction.Size - 1; i > 0; i-- {
		fmt.Printf("%02x", buffer[offset+uint16(i)])
	}

	fmt.Printf("\t")

	return instruction.Size
}

// New Cpu
func New() CPU {
	cpu := CPU{}
	cpu.Memory = make([]byte, 0x10000) // 16Kb memory
	cpu.PC = 0

	return cpu
}

// MemRead reads byte from memory
func (cpu *CPU) MemRead(offset uint16) uint8 {
	if offset >= 0x4000 {
		fmt.Printf("%04x\n", offset)
		panic("Attempt to Read over the RAM limit")
	}
	return cpu.Memory[offset]
}

// MemWrite writes byte to memory
func (cpu *CPU) MemWrite(offset uint16, value uint8) {
	if offset < RomOffset {
		fmt.Println(offset)
		panic("Attempt to write ROM Memory")
	}
	if offset >= 0x4000 {
		panic("Attempt to write over the RAM limit")
	}
	cpu.Memory[offset] = value
}

// NextWord from cpu memory at pc
func (cpu *CPU) NextWord() uint16 {
	pc := cpu.PC
	return (uint16(cpu.MemRead(pc+1)) << 8) | uint16(cpu.MemRead(pc))
}

// NextByte from cpu memory at pc
func (cpu *CPU) NextByte() uint8 {
	if cpu.PC >= 0x4000 {
		fmt.Printf("%04x\n", cpu.PC)
		panic("Attempt to Read over the RAM limit")
	}
	return cpu.Memory[cpu.PC]
}

// push a word to stack
func (cpu *CPU) pushStack(value uint16) {
	sp := cpu.SP
	cpu.MemWrite(sp-1, uint8(value>>8))
	cpu.MemWrite(sp-2, uint8(value))
	cpu.SP -= 2
}

// pop a word from stack
func (cpu *CPU) popStack() uint16 {
	value := (uint16(cpu.MemRead(cpu.SP+1)) << 8) | uint16(cpu.MemRead(cpu.SP))
	cpu.SP += 2
	return value
}

// HL register pair value
func (cpu *CPU) hlValue() uint16 {
	return (uint16(cpu.H) << 8) | uint16(cpu.L)
}

// Step the CPU execution and return the number of cycles
func (cpu *CPU) Step(debug bool) {
	instructionOffset := cpu.PC
	cpu.runInstruction(cpu.PC)

	if debug {
		cpu.debug(instructionOffset)
	}
}

func (cpu *CPU) debug(instructionOffset uint16) {
	Disassembly(cpu.Memory, instructionOffset)

	var (
		zs = "."
		ss = "."
		ps = "."
		cs = "."
		as = "."
	)

	if cpu.Flags.Get(Z) {
		zs = "Z"
	}
	if cpu.Flags.Get(S) {
		ss = "S"
	}
	if cpu.Flags.Get(P) {
		ps = "P"
	}
	if cpu.Flags.Get(CY) {
		cs = "CY"
	}
	if cpu.Flags.Get(AC) {
		as = "AC"
	}

	fmt.Printf("%s%s%s%s%s ", zs, ss, ps, cs, as)

	fmt.Printf("A %02x, B %02x, C %02x, D %02x, E %02x, H %02x, L %02x, SP %04x C %09d\n",
		cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L, cpu.SP, cpu.Cycles)
}

// Interrupt the cpu execution
func (cpu *CPU) Interrupt(address uint8) {
	cpu.pushStack(cpu.PC)
	cpu.PC = uint16(address)
	cpu.IntEnable = false
}

// Run a single instruction at the offset
func (cpu *CPU) runInstruction(offset uint16) {
	opcode := cpu.MemRead(offset)
	cpu.PC++

	cpu.Cycles += uint32(InstructionTable[opcode].Cycles)

	switch opcode {
	case 0x00: // NOP
		break
	case 0x01: // LXI B,word
		lxi(cpu, &cpu.B, &cpu.C, cpu.NextWord())
	case 0x02: // STAX B
		stax(cpu, cpu.B, cpu.C)
	case 0x03: // INX B
		inx(cpu, &cpu.B, &cpu.C)
	case 0x04: // INR B
		inr(cpu, &cpu.B)
	case 0x05: // DCR B
		dcr(cpu, &cpu.B)
	case 0x06: // MVI B,byte
		mvi(cpu, &cpu.B, cpu.NextByte())
	case 0x07: // RLC
		rlc(cpu)

	// case 0x08: // -
	// 	break
	case 0x09: // DAD B
		dad(cpu, cpu.B, cpu.C)
	case 0x0a: // LDAX B
		ldax(cpu, cpu.B, cpu.C)
	case 0x0b: // DCX B
		dcx(cpu, &cpu.B, &cpu.C)
	case 0x0c: // INR C
		inr(cpu, &cpu.C)
	case 0x0d: // DCR C
		dcr(cpu, &cpu.C)
	case 0x0e: // MVI C,byte
		mvi(cpu, &cpu.C, cpu.NextByte())
	case 0x0f: // RRC
		rrc(cpu)

	// case 0x10: // -
	// 	break
	case 0x11: // LXI D,word
		lxi(cpu, &cpu.D, &cpu.E, cpu.NextWord())
	case 0x12: // STAX D
		stax(cpu, cpu.D, cpu.E)
	case 0x13: // INX D
		inx(cpu, &cpu.D, &cpu.E)
	case 0x14: // INR D
		inr(cpu, &cpu.D)
	case 0x15: // DCR D
		dcr(cpu, &cpu.D)
	case 0x16: // MVI D,byte
		mvi(cpu, &cpu.D, cpu.NextByte())
	case 0x17: //RAL
		ral(cpu)

	// case 0x18: // -
	// 	break
	case 0x19: // DAD D
		dad(cpu, cpu.D, cpu.E)
	case 0x1a: // LDAX D
		ldax(cpu, cpu.D, cpu.E)
	case 0x1b: // DCX D
		dcx(cpu, &cpu.D, &cpu.E)
	case 0x1c: // INR E
		inr(cpu, &cpu.E)
	case 0x1d: // DCR E
		dcr(cpu, &cpu.E)
	case 0x1e: // MVI E,byte
		mvi(cpu, &cpu.E, cpu.NextByte())
	case 0x1f: // RAR
		rar(cpu)

	// case 0x20: // -
	// 	break
	case 0x21: // LXI H,word
		lxi(cpu, &cpu.H, &cpu.L, cpu.NextWord())
	case 0x22: //SHLD addr
		shld(cpu, cpu.NextWord())
	case 0x23: // INX H
		inx(cpu, &cpu.H, &cpu.L)
	case 0x24: // INR H
		inr(cpu, &cpu.H)
	case 0x25: // DCR H
		dcr(cpu, &cpu.H)
	case 0x26: // MVI H,byte
		mvi(cpu, &cpu.H, cpu.NextByte())

	// case 0x28: // -
	// 	break
	case 0x29: // DAD H
		dad(cpu, cpu.H, cpu.L)
	case 0x2a: // LHLD addr
		lhld(cpu, cpu.NextWord())
	case 0x2b: // DCX H
		dcx(cpu, &cpu.H, &cpu.L)
	case 0x2c: // INR L
		inr(cpu, &cpu.L)
	case 0x2d: // DCR L
		dcr(cpu, &cpu.L)
	case 0x2e: // MVI L,byte
		mvi(cpu, &cpu.L, cpu.NextByte())
	case 0x2f: // CMA
		cpu.A = ^cpu.A

	// case 0x30: // -
	// 	break
	case 0x31: // LXI SP,word
		lxiSP(cpu, cpu.NextWord())
	case 0x32: // STA addr
		sta(cpu, cpu.NextWord())
	case 0x33: // INX SP
		cpu.SP++
	case 0x34: // INR M
		inrM(cpu)
	case 0x35: // DCR M
		dcrM(cpu)
	case 0x36: // MVI M,byte
		mviM(cpu, cpu.NextByte())
	case 0x37: // STC
		cpu.Flags.Set(CY, true)

	// case 0x38: // -
	// 	break
	case 0x39: // DAD SP
		dadSP(cpu)
	case 0x3a: // LDA addr
		lda(cpu, cpu.NextWord())
	case 0x3b: // DCX SP
		cpu.SP--
	case 0x3c: // INR A
		inr(cpu, &cpu.A)
	case 0x3d: // DCR A
		dcr(cpu, &cpu.A)
	case 0x3e: // MVI A,byte
		mvi(cpu, &cpu.A, cpu.NextByte())
	case 0x3f: // CMC
		cpu.Flags.Set(CY, !cpu.Flags.Get(CY))

	case 0x40: // MOV B,B
		mov(cpu, &cpu.B, cpu.B)
	case 0x41: // MOV B,C
		mov(cpu, &cpu.B, cpu.C)
	case 0x42: // MOV B,D
		mov(cpu, &cpu.B, cpu.D)
	case 0x43: // MOV B,E
		mov(cpu, &cpu.B, cpu.E)
	case 0x44: // MOV B,H
		mov(cpu, &cpu.B, cpu.H)
	case 0x45: // MOV B,L
		mov(cpu, &cpu.B, cpu.L)
	case 0x46: // MOV B,M
		mov(cpu, &cpu.B, cpu.MemRead(cpu.hlValue()))
	case 0x47: // MOV B,A
		mov(cpu, &cpu.B, cpu.A)

	case 0x48: // MOV C,B
		mov(cpu, &cpu.C, cpu.B)
	case 0x49: // MOV C,C
		mov(cpu, &cpu.C, cpu.C)
	case 0x4a: // MOV C,D
		mov(cpu, &cpu.C, cpu.D)
	case 0x4b: // MOV C,E
		mov(cpu, &cpu.C, cpu.E)
	case 0x4c: // MOV C,H
		mov(cpu, &cpu.C, cpu.H)
	case 0x4d: // MOV C,L
		mov(cpu, &cpu.C, cpu.L)
	case 0x4e: // MOV C,M
		mov(cpu, &cpu.C, cpu.MemRead(cpu.hlValue()))
	case 0x4f: // MOV C,A
		mov(cpu, &cpu.C, cpu.A)

	case 0x50: // MOV D,B
		mov(cpu, &cpu.D, cpu.B)
	case 0x51: // MOV D,C
		mov(cpu, &cpu.D, cpu.C)
	case 0x52: // MOV D,D
		mov(cpu, &cpu.D, cpu.D)
	case 0x53: // MOV D,E
		mov(cpu, &cpu.D, cpu.E)
	case 0x54: // MOV D,H
		mov(cpu, &cpu.D, cpu.H)
	case 0x55: // MOV D,L
		mov(cpu, &cpu.D, cpu.L)
	case 0x56: // MOV D,M
		mov(cpu, &cpu.D, cpu.MemRead(cpu.hlValue()))
	case 0x57: // MOV D,A
		mov(cpu, &cpu.D, cpu.A)

	case 0x58: // MOV E,B
		mov(cpu, &cpu.E, cpu.B)
	case 0x59: // MOV E,C
		mov(cpu, &cpu.E, cpu.C)
	case 0x5a: // MOV E,D
		mov(cpu, &cpu.E, cpu.D)
	case 0x5b: // MOV E,E
		mov(cpu, &cpu.E, cpu.E)
	case 0x5c: // MOV E,H
		mov(cpu, &cpu.E, cpu.H)
	case 0x5d: // MOV E,L
		mov(cpu, &cpu.E, cpu.L)
	case 0x5e: // MOV E,M
		mov(cpu, &cpu.E, cpu.MemRead(cpu.hlValue()))
	case 0x5f: // MOV E,A
		mov(cpu, &cpu.E, cpu.A)

	case 0x60: // mov H,B
		mov(cpu, &cpu.H, cpu.B)
	case 0x61: // mov H,C
		mov(cpu, &cpu.H, cpu.C)
	case 0x62: // mov H,D
		mov(cpu, &cpu.H, cpu.D)
	case 0x63: // mov H,E
		mov(cpu, &cpu.H, cpu.E)
	case 0x64: // mov H,H
		mov(cpu, &cpu.H, cpu.H)
	case 0x65: // mov H,L
		mov(cpu, &cpu.H, cpu.L)
	case 0x66: // MOV H,M
		mov(cpu, &cpu.H, cpu.MemRead(cpu.hlValue()))
	case 0x67: // MOV H,A
		mov(cpu, &cpu.H, cpu.A)

	case 0x68: // MOV L,B
		mov(cpu, &cpu.L, cpu.B)
	case 0x69: // MOV L,C
		mov(cpu, &cpu.L, cpu.C)
	case 0x6a: // MOV L,D
		mov(cpu, &cpu.L, cpu.D)
	case 0x6b: // MOV L,E
		mov(cpu, &cpu.L, cpu.E)
	case 0x6c: // MOV L,H
		mov(cpu, &cpu.L, cpu.H)
	case 0x6d: // MOV L,L
		mov(cpu, &cpu.L, cpu.L)
	case 0x6e: // MOV L,M
		mov(cpu, &cpu.L, cpu.MemRead(cpu.hlValue()))
	case 0x6f: // MOV L,A
		mov(cpu, &cpu.L, cpu.A)

	case 0x70: // MOV M,B
		movM(cpu, cpu.B)
	case 0x71: // MOV M,C
		movM(cpu, cpu.C)
	case 0x72: // MOV M,D
		movM(cpu, cpu.D)
	case 0x73: // MOV M,E
		movM(cpu, cpu.E)
	case 0x74: // MOV M,H
		movM(cpu, cpu.H)
	case 0x75: // MOV M,L
		movM(cpu, cpu.L)
	case 0x76: // HLT
		unimplemented(cpu)
	case 0x77: // MOV M,A
		movM(cpu, cpu.A)

	case 0x78: // MOV A,B
		mov(cpu, &cpu.A, cpu.B)
	case 0x79: // MOV A,C
		mov(cpu, &cpu.A, cpu.C)
	case 0x7a: // MOV A,D
		mov(cpu, &cpu.A, cpu.D)
	case 0x7b: // MOV A,E
		mov(cpu, &cpu.A, cpu.E)
	case 0x7c: // MOV A,H
		mov(cpu, &cpu.A, cpu.H)
	case 0x7d: // MOV A,L
		mov(cpu, &cpu.A, cpu.L)
	case 0x7e: // MOV A,M
		mov(cpu, &cpu.A, cpu.MemRead(cpu.hlValue()))
	case 0x7f: // MOV A,A
		mov(cpu, &cpu.A, cpu.A)

	case 0x80: // ADD B
		add(cpu, cpu.B)
	case 0x81: // ADD C
		add(cpu, cpu.C)
	case 0x82: // ADD D
		add(cpu, cpu.D)
	case 0x83: // ADD E
		add(cpu, cpu.E)
	case 0x84: // ADD H
		add(cpu, cpu.H)
	case 0x85: // ADD L
		add(cpu, cpu.L)
	case 0x86: // ADD M
		add(cpu, cpu.MemRead(cpu.hlValue()))
	case 0x87: // ADD A
		add(cpu, cpu.A)

	case 0x88: // ADC B
		adc(cpu, cpu.B)
	case 0x89: // ADC C
		adc(cpu, cpu.C)
	case 0x8a: // ADC D
		adc(cpu, cpu.D)
	case 0x8b: // ADC E
		adc(cpu, cpu.E)
	case 0x8c: // ADC H
		adc(cpu, cpu.H)
	case 0x8d: // ADC L
		adc(cpu, cpu.L)
	case 0x8e: // ADC M
		adc(cpu, cpu.MemRead(cpu.hlValue()))
	case 0x8f: // ADC A
		adc(cpu, cpu.A)

	case 0x90: // SUB B
		sub(cpu, cpu.B)
	case 0x91: // SUB C
		sub(cpu, cpu.C)
	case 0x92: // SUB D
		sub(cpu, cpu.D)
	case 0x93: // SUB E
		sub(cpu, cpu.E)
	case 0x94: // SUB H
		sub(cpu, cpu.H)
	case 0x95: // SUB L
		sub(cpu, cpu.L)
	case 0x96: // SUB M
		sub(cpu, cpu.MemRead(cpu.hlValue()))
	case 0x97: // SUB A
		sub(cpu, cpu.A)

	case 0x98: // SBB B
		sbb(cpu, cpu.B)
	case 0x99: // SBB C
		sbb(cpu, cpu.C)
	case 0x9a: // SBB D
		sbb(cpu, cpu.D)
	case 0x9b: // SBB E
		sbb(cpu, cpu.E)
	case 0x9c: // SBB H
		sbb(cpu, cpu.H)
	case 0x9d: // SBB L
		sbb(cpu, cpu.L)
	case 0x9e: // SBB M
		sbb(cpu, cpu.MemRead(cpu.hlValue()))
	case 0x9f: // SBB A
		sbb(cpu, cpu.A)

	case 0xa0: // ANA B
		ana(cpu, cpu.B)
	case 0xa1: // ANA C
		ana(cpu, cpu.C)
	case 0xa2: // ANA D
		ana(cpu, cpu.D)
	case 0xa3: // ANA E
		ana(cpu, cpu.E)
	case 0xa4: // ANA H
		ana(cpu, cpu.H)
	case 0xa5: // ANA L
		ana(cpu, cpu.L)
	case 0xa6: // ANA M
		ana(cpu, cpu.MemRead(cpu.hlValue()))
	case 0xa7: // ANA A
		ana(cpu, cpu.A)

	case 0xa8: // XRA B
		xra(cpu, cpu.B)
	case 0xa9: // XRA C
		xra(cpu, cpu.C)
	case 0xaa: // XRA D
		xra(cpu, cpu.D)
	case 0xab: // XRA E
		xra(cpu, cpu.E)
	case 0xac: // XRA H
		xra(cpu, cpu.H)
	case 0xad: // XRA L
		xra(cpu, cpu.L)
	case 0xae: // XRA M
		xra(cpu, cpu.MemRead(cpu.hlValue()))
	case 0xaf: // XRA A
		xra(cpu, cpu.A)

	case 0xb0: // ORA B
		ora(cpu, cpu.B)
	case 0xb1: // ORA C
		ora(cpu, cpu.C)
	case 0xb2: // ORA D
		ora(cpu, cpu.D)
	case 0xb3: // ORA E
		ora(cpu, cpu.E)
	case 0xb4: // ORA H
		ora(cpu, cpu.H)
	case 0xb5: // ORA L
		ora(cpu, cpu.L)
	case 0xb6: // ORA M
		ora(cpu, cpu.MemRead(cpu.hlValue()))
	case 0xb7: // ORA A
		ora(cpu, cpu.A)

	case 0xb8: // CMP B
		cmp(cpu, cpu.B)
	case 0xb9: // CMP C
		cmp(cpu, cpu.C)
	case 0xba: // CMP D
		cmp(cpu, cpu.D)
	case 0xbb: // CMP E
		cmp(cpu, cpu.E)
	case 0xbc: // CMP H
		cmp(cpu, cpu.H)
	case 0xbd: // CMP L
		cmp(cpu, cpu.L)
	case 0xbe: // CMP M
		cmp(cpu, cpu.MemRead(cpu.hlValue()))
	case 0xbf: // CMP A
		cmp(cpu, cpu.A)

	case 0xc0: // RNZ
		rnz(cpu)
	case 0xc1: // POP B
		pop(cpu, &cpu.B, &cpu.C)
	case 0xc2: // JNZ addr
		jnz(cpu, cpu.NextWord())
	case 0xc3: // JMP addr
		jmp(cpu, cpu.NextWord())
	case 0xc4: // CNZ
		cnz(cpu, cpu.NextWord())
	case 0xc5: // PUSH B
		push(cpu, cpu.B, cpu.C)
	case 0xc6: // ADI byte
		adi(cpu, cpu.NextByte())
	case 0xc7: // RST 0
		rst(cpu, 0)

	case 0xc8: // RZ
		rz(cpu)
	case 0xc9: // RET
		ret(cpu)
	case 0xca: // JZ
		jz(cpu, cpu.NextWord())
	// case 0xcb: // -
	// 	break
	case 0xcc: // CZ
		cz(cpu, cpu.NextWord())
	case 0xcd: // CALL addr
		call(cpu, cpu.NextWord())
	case 0xce: // ACI byte
		aci(cpu, cpu.NextByte())
	case 0xcf: // RST 1
		rst(cpu, 1)

	case 0xd0: // RNC
		rnc(cpu)
	case 0xd1: // POP D
		pop(cpu, &cpu.D, &cpu.E)
	case 0xd2: // JNC addr
		jnc(cpu, cpu.NextWord())
	case 0xd3: // OUT byte (hardware specific)
		break
	case 0xd4: // CNC addr
		cnc(cpu, cpu.NextWord())
	case 0xd5: // PUSH D
		push(cpu, cpu.D, cpu.E)
	case 0xd6: // SUI byte
		sui(cpu, cpu.NextByte())
	case 0xd7: // RST 2
		rst(cpu, 2)

	case 0xd8: // RC
		rc(cpu)
	// case 0xd9: // -
	// 	break
	case 0xda: // JC addr
		jc(cpu, cpu.NextWord())
	case 0xdb: // IN byte (hardware specific)
		break
	case 0xdc: // CC addr
		cc(cpu, cpu.NextWord())
	// case 0xdd: // -
	// 	break
	case 0xde: // SBI byte
		sbi(cpu, cpu.NextByte())
	case 0xdf: // RST 3
		rst(cpu, 3)

	case 0xe0: // RPO
		rpo(cpu)
	case 0xe1: // POP H
		pop(cpu, &cpu.H, &cpu.L)
	case 0xe2: // JPO
		jpo(cpu, cpu.NextWord())
	case 0xe3: // XTHL
		xthl(cpu)
	case 0xe4: // CPO addr
		cpo(cpu, cpu.NextWord())
	case 0xe5: // PUSH H
		push(cpu, cpu.H, cpu.L)
	case 0xe6: // ANI byte
		ani(cpu, cpu.NextByte())
	case 0xe7: // RST 4
		rst(cpu, 4)

	case 0xe8: // RPE
		rpe(cpu)
	case 0xe9: // PCHL
		pchl(cpu)
	case 0xea: // JPE addr
		jpe(cpu, cpu.NextWord())
	case 0xeb: // XCHG
		xchg(cpu)
	case 0xec: // CPE addr
		cpe(cpu, cpu.NextWord())
	// case 0xed: // -
	// 	break
	case 0xee: // XRI byte
		xri(cpu, cpu.NextByte())
	case 0xef: // RST 5
		rst(cpu, 5)

	case 0xf0: // RP
		rp(cpu)
	case 0xf1: // POP PSW
		popPSW(cpu)
	case 0xf2: // JP addr
		jp(cpu, cpu.NextWord())
	case 0xf3: // DI
		cpu.IntEnable = false
	case 0xf4: // CP addr
		cp(cpu, cpu.NextWord())
	case 0xf5: // PUSH PSW
		pushPSW(cpu)
	case 0xf6: // ORI byte
		ori(cpu, cpu.NextByte())
	case 0xf7: // RST 6
		rst(cpu, 6)

	case 0xf8: // RM
		rm(cpu)
	case 0xf9: // SPHL
		sphl(cpu)
	case 0xfa: // JM addr
		jm(cpu, cpu.NextWord())
	case 0xfb: // EI
		cpu.IntEnable = true
	case 0xfc: // CM addr
		cm(cpu, cpu.NextWord())
	// case 0xfd: // -
	// 	break
	case 0xfe: // CPI byte
		cpi(cpu, cpu.NextByte())
	case 0xff: // RST 7
		rst(cpu, 7)
	default:
		unimplemented(cpu)
	}
}
