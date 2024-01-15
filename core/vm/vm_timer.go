package vm

import (
	"math/big"
	"time"
	// "sync/atomic"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// var timer_vm_duration time.Duration

// func resetTimerDuration() {
// 	atomic.StoreInt64((*int64)(&timer_vm_duration), 0)
// }

// func addTimerDuration(delta time.Duration) {
// 	atomic.AddInt64((*int64)(&timer_vm_duration), (int64)(delta))
// }

// func getTimerDuration() time.Duration {
// 	return time.Duration(atomic.LoadInt64((*int64)(&timer_vm_duration)))
// }

type VMTimer struct {
	started      bool // whether timer is started
	in_execution bool // whether timer is in the interpreter execution loop
	last_time    time.Time
	total_time   time.Duration
	dispatches   uint64
}

func (t *VMTimer) GetTotalTime() time.Duration {
	return t.total_time
}

func (t *VMTimer) GetDispatches() uint64 {
	return t.dispatches
}

func (t *VMTimer) StartTimer() {
	if t.started {
		panic("Timer already started")
	}
	if t.in_execution {
		panic("Already in execution Loop")
	}
	t.started = true
	t.last_time = time.Now()
	t.in_execution = true
}

func (t *VMTimer) DBStartTimer() {
	if !t.in_execution {
		return
	}
	if t.started {
		panic("Timer already started")
	}
	t.started = true
	t.last_time = time.Now()
}

func (t *VMTimer) StopTimer() {
	if t.started == false {
		panic("Timer already stopped")
	}
	if t.in_execution == false {
		panic("Not in execution Loop")
	}
	elapsed := time.Since(t.last_time)
	t.total_time += elapsed
	t.started = false
	t.in_execution = false
}

func (t *VMTimer) DBStopTimer() {
	if !t.in_execution {
		return
	}
	if t.started == false {
		panic("Timer already stopped")
	}
	elapsed := time.Since(t.last_time)
	t.total_time += elapsed
	t.started = false
}

func InstallTimer(vm *EVM, timer *VMTimer) {
	vm.VMTimer = timer
}

// create a proxy for stateDB that has a built-in timer installed
func NewLoggerProxy(db StateDB, timer *VMTimer) StateDB {
	return &timerProxiedDB{
		db:      db,
		vmTimer: timer,
	}
}

type timerProxiedDB struct {
	db      StateDB
	vmTimer *VMTimer
}

func (s *timerProxiedDB) CreateAccount(addr common.Address) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.CreateAccount(addr)
}

func (s *timerProxiedDB) Exist(addr common.Address) bool {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.Exist(addr)
	return res
}

func (s *timerProxiedDB) Empty(addr common.Address) bool {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.Empty(addr)
	return res
}

func (s *timerProxiedDB) Suicide(addr common.Address) bool {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.Suicide(addr)
	return res
}

func (s *timerProxiedDB) HasSuicided(addr common.Address) bool {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.HasSuicided(addr)
	return res
}

func (s *timerProxiedDB) GetBalance(addr common.Address) *big.Int {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetBalance(addr)
	return res
}

func (s *timerProxiedDB) AddBalance(addr common.Address, value *big.Int) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddBalance(addr, value)
}

func (s *timerProxiedDB) SubBalance(addr common.Address, value *big.Int) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.SubBalance(addr, value)
}

func (s *timerProxiedDB) GetNonce(addr common.Address) uint64 {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetNonce(addr)
	return res
}

func (s *timerProxiedDB) SetNonce(addr common.Address, value uint64) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.SetNonce(addr, value)
}

func (s *timerProxiedDB) GetCommittedState(addr common.Address, key common.Hash) common.Hash {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetCommittedState(addr, key)
	return res
}

func (s *timerProxiedDB) GetState(addr common.Address, key common.Hash) common.Hash {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetState(addr, key)
	return res
}

func (s *timerProxiedDB) SetState(addr common.Address, key common.Hash, value common.Hash) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.SetState(addr, key, value)
}

func (s *timerProxiedDB) GetCode(addr common.Address) []byte {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetCode(addr)
	return res
}

func (s *timerProxiedDB) GetCodeSize(addr common.Address) int {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetCodeSize(addr)
	return res
}

func (s *timerProxiedDB) GetCodeHash(addr common.Address) common.Hash {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetCodeHash(addr)
	return res
}

func (s *timerProxiedDB) SetCode(addr common.Address, code []byte) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.SetCode(addr, code)
}

func (s *timerProxiedDB) Snapshot() int {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.Snapshot()
	return res
}

func (s *timerProxiedDB) RevertToSnapshot(id int) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.RevertToSnapshot(id)
}

func (s *timerProxiedDB) AddRefund(amount uint64) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddRefund(amount)
}

func (s *timerProxiedDB) SubRefund(amount uint64) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.SubRefund(amount)
}

func (s *timerProxiedDB) GetRefund() uint64 {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.GetRefund()
	return res
}

func (s *timerProxiedDB) PrepareAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.PrepareAccessList(sender, dest, precompiles, txAccesses)
}

func (s *timerProxiedDB) AddressInAccessList(addr common.Address) bool {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	res := s.db.AddressInAccessList(addr)
	return res
}

func (s *timerProxiedDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	a, b := s.db.SlotInAccessList(addr, slot)
	return a, b
}

func (s *timerProxiedDB) AddAddressToAccessList(addr common.Address) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddAddressToAccessList(addr)
}

func (s *timerProxiedDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddSlotToAccessList(addr, slot)
}

func (s *timerProxiedDB) AddLog(entry *types.Log) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddLog(entry)
}

func (s *timerProxiedDB) AddPreimage(hash common.Hash, data []byte) {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	s.db.AddPreimage(hash, data)
}

func (s *timerProxiedDB) ForEachStorage(addr common.Address, op func(common.Hash, common.Hash) bool) error {
	defer func() {
		s.vmTimer.DBStartTimer()
	}()
	s.vmTimer.DBStopTimer()
	return s.db.ForEachStorage(addr, op)
}
