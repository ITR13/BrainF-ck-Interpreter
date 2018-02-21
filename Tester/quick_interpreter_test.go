// test_interpreter.go
package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	simpleA := &SimpleSegment{
		0,
		[]Operation{{0, 'A'}},
		[]Operation{},
		&SimpleSegment{
			0,
			[]Operation{},
			[]Operation{{0, 0}},
			nil,
		},
	}
	data := MakeData([]byte{})
	simpleA.Run(data)
	outputbytes := data.output.Bytes()
	if len(outputbytes) != 1 {
		t.Errorf(
			"Failed simple A: wrong number of elements: got %d wanted %d (%s)",
			len(outputbytes), 1, string(outputbytes),
		)
	} else if outputbytes[0] != 'A' {
		t.Errorf(
			"Failed simple A: got %d wanted %d",
			outputbytes[0], 'A',
		)
	}
	simpleloop := &SimpleSegment{
		0,
		[]Operation{{0, 1}},
		[]Operation{},
		&LoopSegment{
			1,
			&SimpleSegment{
				0,
				[]Operation{{0, 1}},
				[]Operation{{0, 0}},
				nil,
			},
			nil,
		},
	}
	data = MakeData([]byte{})
	simpleloop.Run(data)
	outputbytes = data.output.Bytes()
	if len(outputbytes) != 255 {
		t.Errorf(
			"Failed simple loop: got %d elements wanted %d (%s)",
			len(outputbytes), 255, string(outputbytes),
		)
	} else {
		for i := byte(1); i != 0; i++ {
			if outputbytes[i-1] != i {
				t.Errorf(
					"Failed simple loop: got %d, wanted %d",
					outputbytes[i-1], i,
				)

			}
		}
	}
}
func TestCompileAndRun(t *testing.T) {
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
testLoop:
	for i := range tests {
		input, err := readFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := readFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		inputs := strings.SplitN(string(input), "!", 2)
		var data *Data
		var compiled Segment
		if len(inputs) == 1 {
			compiled = Compile(input)
			data = MakeData([]byte{})
		} else {
			inputs[1] = strings.Replace(inputs[1], "!", string(0), -1)
			compiled = Compile([]byte(inputs[0]))
			data = MakeData([]byte(inputs[1]))
		}
		compiled = Optimize(compiled)
		compiled.Run(data)
		bufferbytes := data.output.Bytes()
		//t.Log(compiled)
		if len(bufferbytes) != len(output) {
			t.Errorf(
				"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
				tests[i],
				string(bufferbytes),
				string(output),
			)
		} else {
			for j := range bufferbytes {
				if bufferbytes[j] != output[j] {
					t.Errorf(
						"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
						tests[i],
						string(bufferbytes),
						string(output),
					)
					continue testLoop
				}
			}
			t.Logf("Succeeded %s\n", tests[i])
		}
	}
}

func BenchmarkCompileAndRun(b *testing.B) {
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
	for i := range tests {
		b.Run(tests[i], func(b *testing.B) {
			input, err := readFile(tests[i])
			if err != nil {
				panic(err)
			}
			inputs := strings.SplitN(string(input), "!", 2)

			var compiled Segment
			if len(inputs) == 1 {
				compiled = Compile(input)
			} else {
				inputs[1] = strings.Replace(inputs[1], "!", string(0), -1)
				compiled = Compile([]byte(inputs[0]))
			}
			compiled = Optimize(compiled)

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				var data *Data
				if len(inputs) == 1 {
					data = MakeData([]byte{})
				} else {
					data = MakeData([]byte(inputs[1]))
				}
				b.StartTimer()
				compiled.Run(data)
			}
		})
	}
}

func BenchmarkMetaCompileAndRun(b *testing.B) {
	interpreter, err := readFile("../compiled.bf")
	compiled := Compile(interpreter)
	compiled = Optimize(compiled)
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
	for i := range tests {
		b.Run(tests[i], func(b *testing.B) {
			input, err := readFile(tests[i])
			if err != nil {
				panic(err)
			}
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				data := MakeData(input)
				b.StartTimer()
				compiled.Run(data)
			}
		})
	}
}
