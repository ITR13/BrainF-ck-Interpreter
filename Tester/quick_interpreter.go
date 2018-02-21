// quick_interpreter.go
package main

import (
	"bytes"
	"fmt"
)

type Data struct {
	memory        map[int]byte
	memoryPointer int
	input         []byte
	inputPointer  int
	output        *bytes.Buffer
}
type Operation struct {
	offset int
	value  byte
}
type Segment interface {
	Run(data *Data) error
	AddNext(nextSegment Segment) Segment
}

func MakeData(input []byte) *Data {
	return &Data{
		make(map[int]byte),
		0,
		input,
		0,
		&bytes.Buffer{},
	}
}

type StartSegment struct{}

func (ss *StartSegment) AddNext(s Segment) Segment {
	return s
}
func (ss *StartSegment) Run(*Data) error {
	return nil
}

type SimpleSegment struct {
	offset     int
	operations []Operation
	prints     []Operation
	next       Segment
}

func (ss *SimpleSegment) Run(data *Data) error {
	for i := range ss.prints {
		err := data.output.WriteByte(
			data.memory[data.memoryPointer+
				ss.prints[i].offset] +
				ss.prints[i].value,
		)
		if err != nil {
			return err
		}
	}
	for i := range ss.operations {
		data.memory[data.memoryPointer+
			ss.operations[i].offset] += ss.operations[i].value
	}
	data.memoryPointer += ss.offset
	if ss.next != nil {
		return ss.next.Run(data)
	}
	return nil
}
func (ss *SimpleSegment) AddNext(s Segment) Segment {
	if ss.next != nil {
		ss.next = ss.next.AddNext(s)
	} else {
		ss.next = s
	}
	return ss
}

type ScopeError struct{}

func (se *ScopeError) Run(data *Data) error {
	fmt.Fprintln(data.output, "\nScope Error.")
	return nil
}
func (se *ScopeError) AddNext(Segment) Segment {
	return se
}

type LoopSegment struct {
	danger byte //If a number isn't modulo this, then it will infinite loop
	inner  Segment
	next   Segment
}

func (ls *LoopSegment) Run(data *Data) error {
	b := data.memory[data.memoryPointer]
	if b != 0 {
		if ls.inner == nil || b%ls.danger != 0 {
			return fmt.Errorf("Infinite Loop")
		}
		for data.memory[data.memoryPointer] != 0 {
			err := ls.inner.Run(data)
			if err != nil {
				return err
			}
		}
	}
	if ls.next != nil {
		return ls.next.Run(data)
	}
	return nil
}

func (ls *LoopSegment) AddNext(s Segment) Segment {
	if ls.next != nil {
		ls.next = ls.next.AddNext(s)
	} else {
		ls.next = s
	}
	return ls
}

type ReadSegment struct {
	next Segment
}

func (ss *ReadSegment) AddNext(s Segment) Segment {
	if ss.next != nil {
		ss.next = ss.next.AddNext(s)
	} else {
		ss.next = s
	}
	return ss
}
func (ss *ReadSegment) Run(data *Data) error {
	if data.inputPointer < len(data.input) {
		data.memory[data.memoryPointer] = data.input[data.inputPointer]
		data.inputPointer++
	} else {
		data.memory[data.memoryPointer] = 0
	}
	if ss.next != nil {
		return ss.next.Run(data)
	}
	return nil
}

func Compile(program []byte) Segment {
	var start Segment
	start = &StartSegment{}
	for i := 0; i < len(program); i++ {
		switch program[i] {
		case '<', '>', '+', '-', '.':
			simple, offset := CreateSimpleSegment(program, i)
			start = start.AddNext(simple)
			i += offset
		case '[':
			loop, offset := CreateLoopSegment(program, i)
			start = start.AddNext(loop)
			i += offset
		case ']':
			start = start.AddNext(&ScopeError{})
		case ',':
			start = start.AddNext(&ReadSegment{nil})
		}
	}
	return start
}

func CreateSimpleSegment(
	program []byte,
	programPointer int,
) (*SimpleSegment, int) {
	programOffset := 0
	dataOffset := 0
	changes := make(map[int]byte, 0)
	prints := make([]Operation, 0)
simpleLoop:
	for i := programPointer; i < len(program); i++ {
		switch program[i] {
		case '+':
			changes[dataOffset]++
		case '-':
			changes[dataOffset]--
		case '>':
			dataOffset++
		case '<':
			dataOffset--
		case '.':
			prints = append(prints, Operation{dataOffset, changes[dataOffset]})
		case '[', ']', ',':
			break simpleLoop
		}
		programOffset++
	}

	operations := make([]Operation, 0)

	for key, value := range changes {
		if value != 0 {
			operations = append(operations, Operation{key, value})
		}
	}

	return &SimpleSegment{
		dataOffset,
		operations,
		prints,
		nil,
	}, programOffset - 1
}

func CreateLoopSegment(
	program []byte,
	programPointer int,
) (*LoopSegment, int) {
	length := 0
	depth := 0
	for i := programPointer; i < len(program); i++ {
		switch program[i] {
		case '[':
			depth++
		case ']':
			depth--
		}
		if depth == 0 {
			break
		}
		length++
	}
	var next Segment
	if depth != 0 {
		next = &ScopeError{}
	}
	return &LoopSegment{
		1,
		Compile(program[programPointer+1 : programPointer+length]),
		next,
	}, length
}
