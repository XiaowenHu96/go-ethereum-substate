package vm

import (
	"fmt"
	"sort"
)

type Scheduler struct {
	searcher *ACTree
}

type AC_list []ACResult

func (ac_list AC_list) Less(i, j int) bool {
	if ac_list[i].end != ac_list[j].end {
		return ac_list[i].end < ac_list[j].end
	}
	return ac_list[i].start < ac_list[j].start
}

func (ac_list AC_list) Swap(i, j int) {
	ac_list[i], ac_list[j] = ac_list[j], ac_list[i]
}

func (ac_list AC_list) Len() int {
	return len(ac_list)
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func (s *Scheduler) maximal_interval_schedule(intervals []ACResult) ([]ACResult, int) {
	// sort by finish time
	sort.Sort(AC_list(intervals))
	// p[j] is the index of the last interval that finishes before j
	p := make([]int, len(intervals))
	p[0] = 0
	for i := 1; i < len(p); i++ {
		interval := intervals[i]
		begin := interval.start
		lo := 0
		hi := i - 1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if intervals[mid].end <= begin {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		if intervals[lo].end > begin {
			p[i] = 0
		} else {
			p[i] = lo
		}
	}
	m := make([]int, len(intervals))
	res := make([]ACResult, 0)
	m[0] = 0
	for i := 1; i < len(m); i++ {
		// end - start - 1 = length of super-instruction - 1 = dispatches saved
		take := intervals[i].end - intervals[i].start - 1 + m[p[i]]
		no_take := m[i-1]
		if take > no_take {
			m[i] = take
		} else {
			m[i] = no_take
		}
	}
	for i := len(intervals) - 1; i > 0; {
		if m[i] == m[i-1] {
			i--
			continue
		} else {
			res = append(res, intervals[i])
			i = p[i]
		}
	}
	// TODO: Sanity check, recovered schedule should have same gain
	gain := 0
	for i := 0; i < len(res); i++ {
		gain += int(res[i].end - res[i].start - 1)
	}
	if gain != m[len(m)-1] {
		panic("Gain is not equal!")
	}

	// TODO: Sanity check, makeing sure res are not overlapped
	sort.Sort(AC_list(res))
	for i := 0; i < len(res)-1; i++ {
		if res[i].end > res[i+1].start {
			panic(fmt.Sprintf("Overlapped super-instruction! [%d,%d) [%d, %d)",
				res[i].start, res[i].end, res[i+1].start, res[i+1].end))
		}
	}
	return res, m[len(m)-1]
}

func (s *Scheduler) Schedule(code []byte) []ACResult {
	candidates := s.searcher.search(code)
	// add dummy interval
	candidates = append(candidates, ACResult{0, 0, 0})
    schedule, _ := s.maximal_interval_schedule(candidates)
    // fmt.Printf("Speedup: %f\n", float64(len(code)) / float64(len(code) - gain))
    return schedule
}
