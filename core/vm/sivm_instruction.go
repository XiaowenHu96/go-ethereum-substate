package vm

type SIVM_OpCode uint16

// Code for the macro EVM is a slice of instructions
type SIVM_Code []SIVM_OpCode

func (code_stream SIVM_Code) GetOp(n uint64) SIVM_OpCode {
	if n < uint64(len(code_stream)) {
		return code_stream[n]
	}
	return 0
}
