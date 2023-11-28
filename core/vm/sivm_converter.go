package vm

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"sync"
)

var scheduler *Scheduler

type CodeCache struct {
	sivm_code        SIVM_Code
	interpreter_code InterpreterCode
	ratio            float64
}

type InterpreterCode []byte

var cache = map[common.Hash]CodeCache{}
var mu = sync.Mutex{}

func (c InterpreterCode) GetOp(n uint64) OpCode {
	if n < uint64(len(c)) {
		return OpCode((c)[n])
	}
	return 0
}

func SIVMInit() {
	fmt.Printf("Number of super instructions: %d\n", NUM_OF_SI)
	scheduler = &Scheduler{InitSearcherFromDB()}
}

func Convert(code []byte, codeHash common.Hash) CodeCache {
	// dummy convert
	if SIVMDummyConvert {
		ret := make([]SIVM_OpCode, 0)
		for _, c := range code {
			ret = append(ret, SIVM_OpCode(c))
		}
		return CodeCache{ret, code, 1.0}
	}
	// TODO: why caching of 0x0 is not working?
	interpreter_code := make([]byte, len(code))
	copy(interpreter_code, code)
	// fmt.Printf("Code hash: %x\n", codeHash)
	mu.Lock()
	res, ok := cache[codeHash]
	if ok && codeHash != common.BytesToHash([]byte{0x0}) {
		mu.Unlock()
		return res
	} else {
		mu.Unlock()
		sivm_code, ratio := convert(&interpreter_code)
		new_cache := CodeCache{sivm_code, interpreter_code, ratio}
		mu.Lock()
		cache[codeHash] = new_cache
		mu.Unlock()
		return new_cache
	}
}

// 1. Run aho-corasick to identify SI offsets
// 2. Collect offsets and obtian best set using DP
// 3. Translate
func convert(code_ptr *[]byte) (SIVM_Code, float64) {
	// 1. produce clean opcode stream without push args
	// 2. store push args and later recover
	input := *code_ptr
	opcodes := make([]byte, 0)
	push_args := make([]byte, 0)
	for i := 0; i < len(input); i++ {
		opcodes = append(opcodes, input[i])
		if input[i] >= byte(PUSH1) && input[i] <= byte(PUSH32) {
			num := int(input[i]) - int(PUSH1) + 1
			for l := 0; l < num; l++ {
				if i+1+l < len(input) {
					push_args = append(push_args, input[i+1+l])
				}
			}
			i += num
		}
	}
	// start convesion
	push_args_idx := 0
	schedules, esitmated_ratio := scheduler.Schedule(opcodes)
	res := make([]SIVM_OpCode, 0)
	j := 0
	next_start := -1
	for i := 0; i < len(opcodes); {
		// fetch next schedule location
		if j < len(schedules) {
			next_start = schedules[j].start
		} else {
			next_start = len(opcodes)
		}
		// parse until next schedule
		for i < next_start {
			res = append(res, SIVM_OpCode(opcodes[i]))
			if opcodes[i] >= byte(PUSH1) && opcodes[i] <= byte(PUSH32) {
				num := int(opcodes[i]) - int(PUSH1) + 1
				for l := 0; l < num && push_args_idx < len(push_args); l++ {
					res = append(res, SIVM_OpCode(push_args[push_args_idx]))
					push_args_idx++
				}
			}
			i++
		}
		// start parsing the schedule
		if j >= len(schedules) {
			break
		}
		res = append(res, SIVM_OpCode(schedules[j].sym))
        // Currently not used, please ignore:
		(*code_ptr)[len(res)-1] = 0
		// fix push args
		if i != schedules[j].start {
			panic("SIVM conversion failed: schedule start not match")
		}
		for i < schedules[j].end {
			// NOP
			if i != schedules[j].start {
				res = append(res, SIVM_NOP)
			}
			if opcodes[i] >= byte(PUSH1) && opcodes[i] <= byte(PUSH32) {
				num := int(opcodes[i]) - int(PUSH1) + 1
				for l := 0; l < num && push_args_idx < len(push_args); l++ {
					res = append(res, SIVM_OpCode(push_args[push_args_idx]))
					push_args_idx++
				}
			}
			i++
		}
		j++
	}
	// we should obtain a converted code whose length equals the original code
	if len(res) != len(input) {
		panic(fmt.Sprintf("SIVM conversion failed: Length does not match after: %d v.s. before: %d", len(res), len(input)))
	}
	return res, esitmated_ratio
}
