// optimize_compiled.go
package main

import "fmt"

type MultiplierSegment struct {
	inverse, danger byte
	operations      []Operation
	next            Segment
}

func (ms *MultiplierSegment) AddNext(s Segment) Segment {
	if ms.next != nil {
		ms.next = ms.AddNext(s)
	} else {
		ms.next = s
	}
	return ms
}
func (ms *MultiplierSegment) Run(data *Data) error {
	c := data.memory[data.memoryPointer]
	if c != 0 {
		if c%ms.danger != 0 {
			return fmt.Errorf("Infinite Loop")
		}
		multiplier := -c * ms.inverse / ms.danger
		for i := range ms.operations {
			data.memory[data.memoryPointer+
				ms.operations[i].offset] += ms.operations[i].value * multiplier
		}
	}
	if ms.next != nil {
		return ms.next.Run(data)
	}
	return nil
}

func Optimize(s Segment) Segment {
	if s == nil {
		return nil
	}
	switch t := s.(type) {
	case *LoopSegment:
		t.inner = Optimize(t.inner)
		t.next = Optimize(t.next)
		multiplier := FindMultipliers(t)
		if multiplier != nil {
			return multiplier
		}
	case *SimpleSegment:
		t.next = Optimize(t.next)
	case *ReadSegment:
		t.next = Optimize(t.next)
	case *ScopeError, *StartSegment:
	default:
		panic("Switch statement doesn't have full coverage")
	}
	return s
}

func FindMultipliers(ls *LoopSegment) *MultiplierSegment {
	if ls.inner == nil {
		return nil
	}
	ss, ok := ls.inner.(*SimpleSegment)
	if ok {
		if ss.offset != 0 || ss.next != nil {
			return nil
		}
		step := byte(0)
		for i := range ss.operations {
			if ss.operations[i].offset == 0 {
				step = ss.operations[i].value
				break
			}
		}
		if step == 0 {
			ls.inner = nil
			return nil
		}

		div := byte(gcd(int(step), 256))
		coprime := step / div
		inverse := ModInv_2(coprime)
		if len(ss.prints) == 0 {
			return &MultiplierSegment{
				inverse, div,
				ss.operations,
				ls.next,
			}
		} else {
			ls.danger = div
			return nil
		}
	}
	return nil
}

func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
