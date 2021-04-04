package cpu

import (
	"fmt"
	"os"

	"github.com/protoshark/invaders8080/bits"
)

func unimplemented(cpu *CPU) {
	fmt.Printf("Unimplemented instruction: %02x\n", cpu.Memory[cpu.PC-1])
	os.Exit(1)
}

// ! data group

func lxi(cpu *CPU, rh *uint8, rl *uint8, _word uint16) {
	*rh = uint8(_word >> 8)
	*rl = uint8(_word & 0xff)

	cpu.PC += 2
} // OK

func lxiSP(cpu *CPU, _word uint16) {
	cpu.SP = _word

	cpu.PC += 2
} // OK

func mvi(cpu *CPU, r *uint8, _byte uint8) {
	*r = _byte
	cpu.PC++
} // OK

func mviM(cpu *CPU, _byte uint8) {
	cpu.MemWrite(cpu.hlValue(), _byte)
	cpu.PC++
} // OK

func lda(cpu *CPU, addr uint16) {
	cpu.A = cpu.MemRead(addr)
	cpu.PC += 2
} // OK

func ldax(cpu *CPU, rh uint8, rl uint8) {
	offset := (uint16(rh) << 8) | uint16(rl)
	cpu.A = cpu.MemRead(offset)
} // OK

func shld(cpu *CPU, addr uint16) {
	cpu.MemWrite(addr+1, cpu.H)
	cpu.MemWrite(addr, cpu.L)

	cpu.PC += 2
} // OK

func lhld(cpu *CPU, addr uint16) {
	cpu.H = cpu.MemRead(addr + 1)
	cpu.L = cpu.MemRead(addr)

	cpu.PC += 2
} // OK

func sta(cpu *CPU, addr uint16) {
	cpu.MemWrite(addr, cpu.A)
	cpu.PC += 2
} // OK

func stax(cpu *CPU, rh uint8, rl uint8) {
	addr := (uint16(rh) << 8) | uint16(rl)
	cpu.MemWrite(addr, cpu.A)
} // OK

func mov(cpu *CPU, r1 *uint8, r2 uint8) {
	*r1 = r2
} // OK

func movM(cpu *CPU, r uint8) {
	cpu.MemWrite(cpu.hlValue(), r)
} // OK

func xchg(cpu *CPU) {
	_h := cpu.H
	cpu.H = cpu.D
	cpu.D = _h

	_l := cpu.L
	cpu.L = cpu.E
	cpu.E = _l
} // OK

// ! arith group

func dcr(cpu *CPU, r *uint8) {
	result := *r - 1
	*r = result
	cpu.SetZSP(result, 8)
} // OK
func dcrM(cpu *CPU) {
	result := cpu.MemRead(cpu.hlValue()) - 1
	cpu.MemWrite(cpu.hlValue(), result)
	cpu.SetZSP(result, 16)
} // OK

func dcx(cpu *CPU, rh *uint8, rl *uint8) {
	result := *rl - 1
	*rl = result
	if result == 0xff {
		result = *rh - 1
		*rh = result
	}
} // OK

func inr(cpu *CPU, r *uint8) {
	result := *r + 1
	*r = result
	cpu.SetZSP(result, 8)
} // OK
func inrM(cpu *CPU) {
	result := cpu.MemRead(cpu.hlValue()) + 1
	cpu.MemWrite(cpu.hlValue(), result)
	cpu.SetZSP(result, 16)
} // OK

func inx(cpu *CPU, rh *uint8, rl *uint8) {
	result := *rl + 1
	*rl = result
	if *rl == 0 {
		result = *rh + 1
		*rh = result
	}
} // OK

func add(cpu *CPU, r uint8) {
	result := uint16(cpu.A) + uint16(r)

	cpu.A = uint8(result & 0xff)

	cpu.Flags.Set(CY, result > 0xff)
	cpu.SetZSP(uint8(result&0xff), 8)
} // OK

func adc(cpu *CPU, r uint8) {
	result := uint16(cpu.A) + uint16(r) + uint16(cpu.Flags.GetValue(CY))

	cpu.A = uint8(result & 0xff)

	cpu.Flags.Set(CY, result > 0xff)
	cpu.SetZSP(uint8(result&0xff), 8)

} // OK

func adi(cpu *CPU, _byte uint8) {
	result := uint16(cpu.A) + uint16(_byte)

	cpu.A = uint8(result & 0xff)

	cpu.Flags.Set(CY, result > 0xff)
	cpu.SetZSP(uint8(result&0xff), 8)

	cpu.PC++
} // OK

func sub(cpu *CPU, r uint8) {
	result := cpu.A - r

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, cpu.A < r)

	cpu.A = result
} // OK

func sui(cpu *CPU, _byte uint8) {
	result := cpu.A - _byte

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, cpu.A < _byte)

	cpu.A = result
	cpu.PC++
} // OK

func sbi(cpu *CPU, _byte uint8) {
	result := cpu.A - _byte - cpu.Flags.GetValue(CY)

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, cpu.A < (_byte+cpu.Flags.GetValue(CY)))

	cpu.A = result
	cpu.PC++
} // OK

// ! logic group

func dad(cpu *CPU, rh uint8, rl uint8) {
	hl := (uint16(cpu.H) << 8) | uint16(cpu.L)
	rp := (uint16(rh) << 8) | uint16(rl)

	result := uint32(hl) + uint32(rp)

	// set carry flag
	cpu.Flags.Set(CY, (result&0xFFFF_0000) != 0)

	cpu.H = uint8(result & 0xFF00 >> 8)
	cpu.L = uint8(result & 0xff)
} // OK

func dadSP(cpu *CPU) {
	hl := (uint16(cpu.H) << 8) | uint16(cpu.L)

	result := uint32(hl) + uint32(cpu.SP)

	// set carry flag
	cpu.Flags.Set(CY, (result&0xFFFF_0000) != 0)

	cpu.H = uint8(result & 0xff00 >> 8)
	cpu.L = uint8(result & 0xff)
} // OK

func rrc(cpu *CPU) {
	temp := cpu.A
	cpu.A = ((temp & 1) << 7) | (temp >> 1)
	cpu.Flags.Set(CY, (temp&1) == 1)
} // OK

func rlc(cpu *CPU) {
	temp := cpu.A
	cpu.A = (temp << 1) | (temp&0x80)>>7
	cpu.Flags.Set(CY, (temp&0x80)>>7 != 0)
} // OK

func rar(cpu *CPU) {
	temp := cpu.A
	cpu.A = (cpu.Flags.GetValue(CY) << 7) | (temp >> 1)
	cpu.Flags.Set(CY, (temp&1) == 1)
} // OK

func xra(cpu *CPU, r uint8) {
	result := cpu.A ^ r

	cpu.SetZSP(result, 8)

	// the CY and AC flags are cleared
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
} // OK

func ana(cpu *CPU, r uint8) {
	result := cpu.A & r
	cpu.A = result

	cpu.SetZSP(result, 8)
	// the CY flag is cleared
	cpu.Flags.Set(CY, false)
} // OK

func ani(cpu *CPU, _byte uint8) {
	result := cpu.A & _byte

	cpu.SetZSP(result, 8)

	// the CY and AC flags is cleared
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
	cpu.PC++
} // OK

func cmp(cpu *CPU, r uint8) {
	result := cpu.A - r

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, cpu.A < r)
} // OK

func cpi(cpu *CPU, _byte uint8) {
	result := cpu.A - _byte

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, cpu.A < _byte)

	cpu.PC++
} // OK

func ora(cpu *CPU, r uint8) {
	result := cpu.A | r

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
} // OK

func ori(cpu *CPU, _byte uint8) {
	result := cpu.A | _byte

	cpu.SetZSP(result, 8)
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
	cpu.PC++
} // OK

// ! branch group

func jmp(cpu *CPU, addr uint16) {
	cpu.PC = addr
} // OK

func jnz(cpu *CPU, addr uint16) {
	if !cpu.Flags.Get(Z) {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
} // OK

func jz(cpu *CPU, addr uint16) {
	if cpu.Flags.Get(Z) {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
} // OK

func jnc(cpu *CPU, addr uint16) {
	if !cpu.Flags.Get(CY) {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
} // OK

func jc(cpu *CPU, addr uint16) {
	if cpu.Flags.Get(CY) {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
} // OK

func jm(cpu *CPU, addr uint16) {
	if cpu.Flags.Get(S) {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
} // OK

func call(cpu *CPU, addr uint16) {
	cpu.pushStack(cpu.PC + 2)
	cpu.PC = addr
} // OK

func cnz(cpu *CPU, addr uint16) {
	if !cpu.Flags.Get(Z) {
		call(cpu, addr)
	} else {
		cpu.PC += 2
		cpu.Cycles += 6
	}
} // OK

func cz(cpu *CPU, addr uint16) {
	if cpu.Flags.Get(Z) {
		call(cpu, addr)
	} else {
		cpu.PC += 2
		cpu.Cycles += 6
	}
} // OK

func cnc(cpu *CPU, addr uint16) {
	if !cpu.Flags.Get(CY) {
		call(cpu, addr)
	} else {
		cpu.PC += 2
		cpu.Cycles += 6
	}
} // OK

func ret(cpu *CPU) {
	cpu.PC = cpu.popStack()
}

func rnz(cpu *CPU) {
	if !cpu.Flags.Get(Z) {
		cpu.PC = cpu.popStack()
	}
} // OK

func rz(cpu *CPU) {
	if cpu.Flags.Get(Z) {
		cpu.PC = cpu.popStack()
	}
}

func rnc(cpu *CPU) {
	if !cpu.Flags.Get(CY) {
		cpu.PC = cpu.popStack()
	}
} // OK

func rc(cpu *CPU) {
	if cpu.Flags.Get(CY) {
		cpu.PC = cpu.popStack()
	}
} // OK

func pchl(cpu *CPU) {
	cpu.PC = (uint16(cpu.H) << 8) | uint16(cpu.L)
} // OK

// ! stack group

func push(cpu *CPU, rh uint8, rl uint8) {
	pair := (uint16(rh) << 8) | uint16(rl)
	cpu.pushStack(pair)
} // OK

func pop(cpu *CPU, rh *uint8, rl *uint8) {
	value := cpu.popStack()
	*rh = uint8(value >> 8)
	*rl = uint8(value & 0xff)
} // OK

func pushPSW(cpu *CPU) {
	value := (uint16(cpu.A) << 8) | uint16(cpu.Flags)
	cpu.pushStack(value)
} // OK

func popPSW(cpu *CPU) {
	value := cpu.popStack()

	hvalue := uint8(value >> 8)
	lvalue := uint8(value & 0xff)

	cpu.A = hvalue
	cpu.Flags = bits.Bitfield(lvalue)
} // OK

func xthl(cpu *CPU) {

	temp := cpu.popStack()

	cpu.pushStack(cpu.hlValue())

	cpu.H = uint8(temp >> 8)
	cpu.L = uint8(temp & 0xff)
} // OK

func sphl(cpu *CPU) {
	cpu.SP = cpu.hlValue()
} // OK

// func rst(cpu *CPU, n uint8) {
// 	cpu.pushStack(cpu.PC + 2)
// 	cpu.PC = uint16(8 * n)
// }
