package vm

type SIVM_OpCode uint16

// // The encoding of each instruction for the MACRO EVM
// type SIVM_Instruction struct {
// 	// The op-code of this instruction.
// 	opcode SIVM_OpCode
// }

// Code for the macro EVM is a slice of instructions
type SIVM_Code []SIVM_OpCode

func (this *SIVM_Code) GetOp(n uint64) SIVM_OpCode {
	if n < uint64(len(*this)) {
		return (*this)[n]
	}
	return 0
}
