package cpu

import (
	"fmt"
	"os"

	"github.com/vnteles/i8080/bits"
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
}

func lxiSP(cpu *CPU, _word uint16) {
	cpu.SP = _word

	cpu.PC += 2
}

func mvi(cpu *CPU, r *uint8, _byte uint8) {
	*r = _byte
	cpu.PC++
}

func mviM(cpu *CPU, _byte uint8) {
	cpu.MemWrite(cpu.hlValue(), _byte)
	cpu.PC++
}

func lda(cpu *CPU, addr uint16) {
	cpu.A = cpu.MemRead(addr)
	cpu.PC += 2
}

func ldax(cpu *CPU, rh uint8, rl uint8) {
	offset := (uint16(rh) << 8) | uint16(rl)
	cpu.A = cpu.MemRead(offset)
}

func shld(cpu *CPU, addr uint16) {
	cpu.MemWrite(addr, cpu.L)
	cpu.MemWrite(addr+1, cpu.H)

	cpu.PC += 2
}

func lhld(cpu *CPU, addr uint16) {
	cpu.L = cpu.MemRead(addr)
	cpu.H = cpu.MemRead(addr + 1)

	cpu.PC += 2
}

func sta(cpu *CPU, addr uint16) {
	cpu.MemWrite(addr, cpu.A)
	cpu.PC += 2
}

func stax(cpu *CPU, rh uint8, rl uint8) {
	addr := (uint16(rh) << 8) | uint16(rl)
	cpu.MemWrite(addr, cpu.A)
}

func mov(cpu *CPU, r1 *uint8, r2 uint8) {
	*r1 = r2
}

func movM(cpu *CPU, r uint8) {
	cpu.MemWrite(cpu.hlValue(), r)
}

func xchg(cpu *CPU) {
	h := cpu.H
	cpu.H = cpu.D
	cpu.D = h

	l := cpu.L
	cpu.L = cpu.E
	cpu.E = l
}

// ! arith group

func dcr(cpu *CPU, r *uint8) {
	*r = *r - 1
	cpu.SetZSP(*r)
}
func dcrM(cpu *CPU) {
	result := cpu.MemRead(cpu.hlValue()) - 1
	cpu.SetZSP(result)
	cpu.MemWrite(cpu.hlValue(), result)
}

func dcx(cpu *CPU, rh *uint8, rl *uint8) {
	*rl = *rl - 1
	if *rl == 0xff {
		*rh = *rh - 1
	}
}

func inr(cpu *CPU, r *uint8) {
	*r = *r + 1
	cpu.SetZSP(*r)
}
func inrM(cpu *CPU) {
	result := cpu.MemRead(cpu.hlValue()) + 1
	cpu.SetZSP(result)
	cpu.MemWrite(cpu.hlValue(), result)
}

func inx(cpu *CPU, rh *uint8, rl *uint8) {
	*rl = *rl + 1
	if *rl == 0 {
		*rh = *rh + 1
	}
}

func add(cpu *CPU, r uint8) {
	result := uint16(cpu.A) + uint16(r)

	cpu.SetZSP(uint8(result & 0xff))
	cpu.Flags.Set(CY, result > 0xff)

	cpu.A = uint8(result & 0xff)
}

func adc(cpu *CPU, r uint8) {
	result := uint16(cpu.A) + uint16(r) + uint16(cpu.Flags.GetValue(CY))

	cpu.SetZSP(uint8(result & 0xff))
	cpu.Flags.Set(CY, result > 0xff)

	cpu.A = uint8(result & 0xff)
}

func adi(cpu *CPU, _byte uint8) {
	result := uint16(cpu.A) + uint16(_byte)

	cpu.SetZSP(uint8(result & 0xff))
	cpu.Flags.Set(CY, result > 0xff)

	cpu.A = uint8(result & 0xff)
	cpu.PC++
}

func aci(cpu *CPU, _byte uint8) {
	result := uint16(cpu.A) + uint16(_byte) + uint16(cpu.Flags.GetValue(CY))

	cpu.SetZSP(uint8(result & 0xff))
	cpu.Flags.Set(CY, result > 0xff)

	cpu.A = uint8(result & 0xff)
	cpu.PC++
}

func sub(cpu *CPU, r uint8) {
	result := cpu.A - r

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < r)

	cpu.A = result
}

func sbb(cpu *CPU, r uint8) {
	c := r - cpu.Flags.GetValue(CY)
	result := cpu.A - c

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < c)

	cpu.A = result
}

func sui(cpu *CPU, _byte uint8) {
	result := cpu.A - _byte

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < _byte)

	cpu.A = result
	cpu.PC++
}

func sbi(cpu *CPU, _byte uint8) {
	c := _byte - cpu.Flags.GetValue(CY)
	result := cpu.A - c

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < c)

	cpu.A = result
	cpu.PC++
}

// ! logic group

func dad(cpu *CPU, rh uint8, rl uint8) {
	hl := (uint16(cpu.H) << 8) | uint16(cpu.L)
	rp := (uint16(rh) << 8) | uint16(rl)

	result := uint32(hl) + uint32(rp)

	// set carry flag
	cpu.Flags.Set(CY, (result&0xffff0000) != 0)

	cpu.H = uint8(uint16(result&0xff00) >> 8)
	cpu.L = uint8(result & 0xff)
}

func dadSP(cpu *CPU) {
	hl := (uint16(cpu.H) << 8) | uint16(cpu.L)

	result := uint32(hl) + uint32(cpu.SP)

	// set carry flag
	cpu.Flags.Set(CY, (result&0xffff0000) != 0)

	cpu.H = uint8(uint16(result&0xff00) >> 8)
	cpu.L = uint8(result & 0xff)
}

func rrc(cpu *CPU) {
	cpu.A = (cpu.A >> 1) | ((cpu.A & 1) << 7)
	cpu.Flags.Set(CY, (cpu.A&1) == 1)
}

func rlc(cpu *CPU) {
	cpu.A = (cpu.A << 1) | ((cpu.A & 0x80) >> 7)
	cpu.Flags.Set(CY, (cpu.A&0x80) == 0x80)
}

func ral(cpu *CPU) {
	cpu.A = (cpu.A << 1) | (cpu.Flags.GetValue(CY))
	cpu.Flags.Set(CY, (cpu.A&0x80) == 0x80)
}

func rar(cpu *CPU) {
	cpu.A = (cpu.A >> 1) | (cpu.Flags.GetValue(CY) << 7)
	cpu.Flags.Set(CY, (cpu.A&1) == 1)
}

func xra(cpu *CPU, r uint8) {
	result := cpu.A ^ r

	cpu.SetZSP(result)

	// the CY and AC flags are cleared
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
}

func xri(cpu *CPU, _byte uint8) {
	xra(cpu, _byte)
	cpu.PC++
}

func ana(cpu *CPU, r uint8) {
	result := cpu.A & r

	cpu.SetZSP(result)

	// the CY flag is cleared
	cpu.Flags.Set(CY, false)

	cpu.A = result
}

func ani(cpu *CPU, _byte uint8) {
	result := cpu.A & _byte

	cpu.SetZSP(result)

	// the CY and AC flags is cleared
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
	cpu.PC++
}

func cmp(cpu *CPU, r uint8) {
	result := cpu.A - r

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < r)
}

func cpi(cpu *CPU, _byte uint8) {
	result := cpu.A - _byte

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, cpu.A < _byte)

	cpu.PC++
}

func ora(cpu *CPU, r uint8) {
	result := cpu.A | r

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
}

func ori(cpu *CPU, _byte uint8) {
	result := cpu.A | _byte

	cpu.SetZSP(result)
	cpu.Flags.Set(CY, false)
	cpu.Flags.Set(AC, false)

	cpu.A = result
	cpu.PC++
}

// ! branch group

func jmp(cpu *CPU, addr uint16) {
	cpu.PC = addr
}

func jmpCondition(cpu *CPU, addr uint16, condition bool) {
	if condition {
		cpu.PC = addr
	} else {
		cpu.PC += 2
	}
}

func jnz(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, !cpu.Flags.Get(Z))
}

func jz(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, cpu.Flags.Get(Z))
}

func jnc(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, !cpu.Flags.Get(CY))
}

func jc(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, cpu.Flags.Get(CY))
}

func jpo(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, !cpu.Flags.Get(P))
}

func jpe(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, cpu.Flags.Get(P))
}

func jm(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, cpu.Flags.Get(S))
}

func jp(cpu *CPU, addr uint16) {
	jmpCondition(cpu, addr, !cpu.Flags.Get(S))
}

func call(cpu *CPU, addr uint16) {
	cpu.pushStack(cpu.PC + 2)
	cpu.PC = addr
}

func callCondition(cpu *CPU, addr uint16, condition bool) {
	if condition {
		call(cpu, addr)
	} else {
		cpu.PC += 2
		cpu.Cycles += 6
	}
}

func cnz(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, !cpu.Flags.Get(Z))
}

func cz(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, cpu.Flags.Get(Z))
}

func cnc(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, !cpu.Flags.Get(CY))
}

func cc(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, cpu.Flags.Get(CY))
}

func cpo(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, !cpu.Flags.Get(P))
}

func cpe(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, cpu.Flags.Get(P))
}

func cm(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, cpu.Flags.Get(S))
}

func cp(cpu *CPU, addr uint16) {
	callCondition(cpu, addr, !cpu.Flags.Get(S))
}

func ret(cpu *CPU) {
	cpu.PC = cpu.popStack()
}

func retCondition(cpu *CPU, condition bool) {
	if condition {
		cpu.PC = cpu.popStack()
	}
}

func rnz(cpu *CPU) {
	retCondition(cpu, !cpu.Flags.Get(Z))
}

func rz(cpu *CPU) {
	retCondition(cpu, cpu.Flags.Get(Z))
}

func rnc(cpu *CPU) {
	retCondition(cpu, !cpu.Flags.Get(CY))
}

func rc(cpu *CPU) {
	retCondition(cpu, cpu.Flags.Get(CY))
}

func rpo(cpu *CPU) {
	retCondition(cpu, !cpu.Flags.Get(P))
}

func rpe(cpu *CPU) {
	retCondition(cpu, cpu.Flags.Get(P))
}

func rm(cpu *CPU) {
	retCondition(cpu, cpu.Flags.Get(S))
}

func rp(cpu *CPU) {
	retCondition(cpu, !cpu.Flags.Get(S))
}

func pchl(cpu *CPU) {
	cpu.PC = (uint16(cpu.H) << 8) | uint16(cpu.L)
}

// ! stack group

func push(cpu *CPU, rh uint8, rl uint8) {
	rp := (uint16(rh) << 8) | uint16(rl)
	cpu.pushStack(rp)
}

func pop(cpu *CPU, rh *uint8, rl *uint8) {
	value := cpu.popStack()
	*rh = uint8(value >> 8)
	*rl = uint8(value & 0xff)
}

func pushPSW(cpu *CPU) {
	value := (uint16(cpu.A) << 8) | uint16(cpu.Flags)
	cpu.pushStack(value)
}

func popPSW(cpu *CPU) {
	value := cpu.popStack()

	hvalue := uint8(value >> 8)
	lvalue := uint8(value & 0xff)

	cpu.A = hvalue
	cpu.Flags = bits.Bitfield(lvalue)
}

func xthl(cpu *CPU) {
	l := cpu.L
	cpu.L = cpu.MemRead(cpu.SP)
	cpu.MemWrite(cpu.SP, l)

	h := cpu.H
	cpu.H = cpu.MemRead(cpu.SP + 1)
	cpu.MemWrite(cpu.SP+1, h)

}

func sphl(cpu *CPU) {
	cpu.SP = (uint16(cpu.H) << 8) | uint16(cpu.L)
}

func rst(cpu *CPU, n uint8) {
	cpu.pushStack(cpu.PC + 2)
	cpu.PC = uint16(8 * n)
}
