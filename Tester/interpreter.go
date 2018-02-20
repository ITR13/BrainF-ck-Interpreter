package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	test_interpreter()
}

func test_interpreter() {
	interpreter, err := ioutil.ReadFile("../compiled.bf")
	//interpreter, err := ioutil.ReadFile("../commented.bf")
	if err != nil {
		panic(err)
	}
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
testLoop:
	for i := range tests {
		buffer := bytes.Buffer{}
		input, err := ioutil.ReadFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := ioutil.ReadFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		for i := range input {
			if input[i] == '!' {
				input[i] = 0
			}
		}
		Interpret(interpreter, append(input, 0), &buffer)
		bufferbytes := buffer.Bytes()
		if len(bufferbytes) != len(output) {
			fmt.Printf(
				"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
				tests[i],
				string(bufferbytes),
				string(output),
			)
		} else {
			for i := range bufferbytes {
				if bufferbytes[i] != output[i] {
					fmt.Printf(
						"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
						tests[i],
						string(bufferbytes),
						string(output),
					)
					continue testLoop
				}
			}
			fmt.Printf("Succeeded %s\n", tests[i])
		}
	}
	fmt.Println("Meta Tests")
metaLoop:
	for i := range tests {
		buffer := bytes.Buffer{}
		input, err := ioutil.ReadFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := ioutil.ReadFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		for i := range input {
			if input[i] == '!' {
				input[i] = 0
			}
		}
		interpreter = append(interpreter, 0)
		input = append(interpreter, input...)
		Interpret(interpreter, append(input, 0), &buffer)
		bufferbytes := buffer.Bytes()
		if len(bufferbytes) != len(output) {
			fmt.Printf(
				"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
				tests[i],
				string(bufferbytes),
				string(output),
			)
		} else {
			for i := range bufferbytes {
				if bufferbytes[i] != output[i] {
					fmt.Printf(
						"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
						tests[i],
						string(bufferbytes),
						string(output),
					)
					continue metaLoop
				}
			}
			fmt.Printf("Succeeded %s\n", tests[i])
		}
	}
}

func Interpret(program []byte, input []byte, output io.Writer) {
	memory := make(map[int]byte)
	programPointer := 0
	memoryPointer := 0
	scope := 0
	for programPointer < len(program) {
		switch string(program[programPointer]) {
		case "+":
			memory[memoryPointer] += 1
		case "-":
			memory[memoryPointer] -= 1
		case "<":
			memoryPointer--
		case ">":
			memoryPointer++
		case "[":
			if memory[memoryPointer] == 0 {
				for depth := 1; depth > 0; {
					programPointer++
					if programPointer >= len(program) {
						fmt.Println("\nScope Error")
						return
					}
					c := program[programPointer]
					if c == '[' {
						depth++
					} else if c == ']' {
						depth--
					}
				}
			} else {
				scope++
			}
		case "]":
			if memory[memoryPointer] != 0 {
				for depth := 1; depth > 0; {
					programPointer--
					if programPointer < 0 {
						fmt.Println("\nScope Error")
						return
					}
					c := program[programPointer]
					if c == '[' {
						depth--
					} else if c == ']' {
						depth++
					}
				}
			} else {
				scope--
				if scope < 0 {
					fmt.Println("\nScope Error")
					return
				}
			}
		case ".":
			fmt.Fprint(output, string(memory[memoryPointer]))
		case ",":
			if len(input) == 0 {
				memory[memoryPointer] = 0
			} else {
				memory[memoryPointer] = input[0]
				input = input[1:]
			}
		case "(": //Used for debugging
			s := "("
			for depth := 1; depth > 0; {
				programPointer++
				c := program[programPointer]
				s += string(c)
				if c == '(' {
					depth++
				} else if c == ')' {
					depth--
				}
			}
			inner := strings.Split(s, "(")[2]
			inner = strings.Split(inner, ")")[0]

			split := strings.Split(inner, " ")

			band := strings.Split(split[0], "b")[1]
			posType := strings.Split(split[1], ":")[0]
			posIndex := strings.Split(split[1], ":")[1]

			var actualPos int
			if memoryPointer < 0 {
				actualPos = (memoryPointer - 4) / 5
			} else {
				actualPos = memoryPointer / 5
			}
			actualBand := memoryPointer % 5
			if actualBand < 0 {
				actualBand += 5
			}
			fmt.Fprintf(
				output,
				"%s - (b%d:%d)\n",
				s, actualBand, actualPos,
			)
			bandN, err := strconv.Atoi(band)
			if err != nil {
				panic(err)
			}
			posIndexN, err := strconv.Atoi(posIndex)
			if err != nil {
				panic(err)
			}

			if bandN != actualBand {
				panic(fmt.Errorf(
					"expected band %d but had %d",
					bandN,
					actualBand,
				))
			} else if posType == "Zero" {
				if posIndexN != actualPos {
					panic(fmt.Errorf(
						"expected posIndex %d but had %d",
						posIndexN,
						actualPos,
					))
				}
			} else if posType == "Right" {
				if posIndexN > actualPos {
					panic(fmt.Errorf(
						"expected posIndex >= %d but had %d",
						-posIndexN,
						actualPos,
					))
				}
			} else if posType == "Left" {
				if posIndexN < actualPos {
					panic(fmt.Errorf(
						"expected posIndex <= %d but had %d",
						-posIndexN,
						actualPos,
					))
				}
			}
		}
		programPointer++
	}
	if scope > 0 {
		fmt.Fprintln(output, "\nScope Error")
	}
}
