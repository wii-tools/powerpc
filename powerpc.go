package powerpc

import "encoding/binary"

// fourByte returns 4 bytes, suitable for the given length.
func fourByte(value uint32) []byte {
	holder := make([]byte, 4)
	binary.BigEndian.PutUint32(holder, value)
	return holder
}

// uint24 returns 3 bytes for an uint24 represented as 4 bytes/an uint32.
func uint24(num uint32) [3]byte {
	if num > 0x00FFFFFF {
		panic("invalid uint24 passed")
	}

	result := fourByte(num)
	return [3]byte{result[1], result[2], result[3]}
}

// Instruction represents a 4-byte PowerPC instruction.
type Instruction [4]byte

// Instructions represents a group of PowerPC instructions.
type Instructions []Instruction

// Bytes returns the represented bytes of these instructions.
//
// As PowerPC has a fixed-width instruction set, the returned
// size will always be in a multiple of four.
func (i Instructions) Bytes() []byte {
	var contents []byte

	for _, instruction := range i {
		contents = append(contents, instruction[:]...)
	}

	return contents
}

// Padding is not an actual instruction - it represents 4 zeros.
var Padding Instruction = [4]byte{0x00, 0x00, 0x00, 0x00}

// BLR represents the blr mnemonic on PowerPC.
func BLR() Instruction {
	return [4]byte{0x4E, 0x80, 0x00, 0x20}
}

// CRXOR represents a common use of CRXOR on PowerPC.
// TODO: actually implement
func CRXOR() Instruction {
	return [4]byte{0x4c, 0xc6, 0x31, 0x82}
}

// ADDI represents the addi PowerPC instruction.
func ADDI(rT Register, rA Register, value uint16) Instruction {
	return EncodeInstrDForm(14, rT, rA, value)
}

// LI represents the li mnemonic on PowerPC.
func LI(rT Register, value uint16) Instruction {
	return ADDI(rT, 0, value)
}

// SUBI represents the subi mnemonic on PowerPC.
// TODO: handle negative values properly?
func SUBI(rT Register, rA Register, value uint16) Instruction {
	return ADDI(rT, 0, -value)
}

// ADDIS represents the addis PowerPC instruction.
func ADDIS(rT Register, rA Register, value uint16) Instruction {
	return EncodeInstrDForm(15, rT, rA, value)
}

// LIS represents the lis mnemonic on PowerPC.
func LIS(rT Register, value uint16) Instruction {
	return ADDIS(rT, 0, value)
}

// OR represents the or PowerPC instruction.
func OR(rS Register, rA Register, rB Register, rC bool) Instruction {
	return EncodeInstrXForm(31, rS, rA, rB, 444, rC)
}

// ORI represents the ori PowerPC instruction.
func ORI(rS Register, rA Register, value uint16) Instruction {
	return EncodeInstrDForm(24, rS, rA, value)
}

// STH represents the sth PowerPC instruction.
func STH(rS Register, offset uint16, rA Register) Instruction {
	return EncodeInstrDForm(44, rS, rA, offset)
}

// EIEIO represents the eieio PowerPC instruction.
func EIEIO() Instruction {
	return [4]byte{0x7C, 0x00, 0x06, 0xAC}
}

// STW represents the stw PowerPC instruction.
func STW(rS Register, offset uint16, rA Register) Instruction {
	return EncodeInstrDForm(36, rS, rA, offset)
}

// LWZ represents the lwz PowerPC instruction.
func LWZ(rT Register, offset uint16, rA Register) Instruction {
	return EncodeInstrDForm(32, rT, rA, offset)
}

// NOP represents the nop mnemonic for PowerPC.
func NOP() Instruction {
	return ORI(R0, R0, 0)
}

// CMPWI represents the cmpwi mnemonic for PowerPC.
// It does not support any other CR fields asides from 0.
func CMPWI(rA Register, value uint16) Instruction {
	return EncodeInstrDForm(11, 0, rA, value)
}

// SYNC is a hack, hardcoding sync 0.
// TODO(spotlightishere): actually encode this
func SYNC() Instruction {
	return [4]byte{0x7c, 0x00, 0x04, 0xac}
}

// MTSPR represents the mtspr PowerPC instruction.
func MTSPR(spr SpecialRegister, rS Register) Instruction {
	return EncodeInstrXFXForm(31, rS, spr, 467, false)
}

// MFSPR represents the mfspr PowerPC instruction.
func MFSPR(rS Register, spr SpecialRegister) Instruction {
	return EncodeInstrXFXForm(31, rS, spr, 339, false)
}

// STWU represents the stwu PowerPC instruction.
func STWU(rS Register, rA Register, offset uint16) Instruction {
	return EncodeInstrDForm(37, rS, rA, offset)
}

// calcDestination determines the proper offset from a given
// calling address and target address.
func calcDestination(from uint, target uint) [3]byte {
	// TODO(spotlightishere): Handle negative offsets properly
	offset := target - from

	// Sign-extend by two bytes
	calc := uint32(offset >> 2)
	return uint24(calc)
}

// BL represents the bl PowerPC instruction.
// It calculates the offset from the given current address and the given
// target address, saving the current address in the link register. It then branches.
func BL(current uint, target uint) Instruction {
	return EncodeInstrIForm(18, calcDestination(current, target), false, true)
}

// B represents the b PowerPC instruction.
// It calculates the offset from the given current address
// and the given target address, and then branches.
func B(current uint, target uint) Instruction {
	return EncodeInstrIForm(18, calcDestination(current, target), false, false)
}

// BNE represents the bne PowerPC instruction.
// It calculates the offset from the given current address
// and the given target address, and then branches.
func BNE(current uint, target uint) Instruction {
	return EncodeInstrBForm(16, 4, 2, calcDestination(current, target), false, false)
}
