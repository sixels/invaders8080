package bits

// Bitfield helper
type Bitfield uint8

// Get the field
func (b Bitfield) Get(field Bitfield) bool {
	return b&field != 0
}

// GetValue of field, since go doesn't convert bool to int
func (b Bitfield) GetValue(field Bitfield) uint8 {
	if b&field != 0 {
		return 1
	}
	return 0
}

// Set the given field
func (b *Bitfield) Set(field Bitfield, value bool) {
	if value {
		*b |= field
	} else {
		*b &= ^field
	}
}

// ByteParity return true if parity even
func ByteParity(x uint8) bool {
	// y := x
	// for i := 0; i < 3; i++ {
	// 	y ^= y >> (1 << i)
	// }
	// return (y & 0b01) == 0b01

	var p int = 0
	x = (x & ((1 << 8) - 1))
	for i := 0; i < 8; i++ {
		if (x & 0x1) == 1 {
			p++
		}
		x = x >> 1
	}
	return (0 == (p & 0x1))
}
