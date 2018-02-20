package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func main() {
	test_interpreter_quick()
	//test_interpreter()
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

func test_interpreter_quick() {
	interpreter, err := ioutil.ReadFile("../compiled.bf")
	if err != nil {
		panic(err)
	}
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
testLoop:
	for i := range tests {
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
		compiled := Compile(interpreter)
		data := MakeData(append(input, 0))
		err = compiled.Run(data)
		if err != nil {
			panic(err)
		}
		bufferbytes := data.output.Bytes()
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
		compiled := Compile(interpreter)
		input = append(interpreter, input...)
		data := MakeData(append(input, 0))
		err = compiled.Run(data)
		if err != nil {
			panic(err)
		}
		bufferbytes := data.output.Bytes()
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
