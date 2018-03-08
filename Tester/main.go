package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

func main() {
	//test_interpreter_quick(false)
	//test_interpreter(false)
	test_all_others()
}

func readFile(path string) ([]byte, error) {
	out, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	for i := len(out) - 1; i >= 0; i-- {
		if out[i] == 0x0D {
			out = append(out[:i], out[i+1:]...)
		}
	}
	return out, nil
}

func test_interpreter(testmeta bool) {
	interpreter, err := readFile("../compiled.bf")
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
		input, err := readFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := readFile(tests[i][:len(tests[i])-2] + "out")
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
	if !testmeta {
		return
	}
	fmt.Println("Meta Tests")
	interpreter = append(interpreter, 0)
metaLoop:
	for i := range tests {
		buffer := bytes.Buffer{}
		input, err := readFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := readFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		for i := range input {
			if input[i] == '!' {
				input[i] = 0
			}
		}
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

func test_interpreter_quick(testmeta bool) {
	interpreter, err := readFile("../compiled.bf")
	if err != nil {
		panic(err)
	}
	compiled := Compile(interpreter)
	compiled = Optimize(compiled)
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
		for i := range input {
			if input[i] == '!' {
				input[i] = 0
			}
		}
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
	if !testmeta {
		return
	}
	fmt.Println("Meta Tests")
	interpreter = append(interpreter, 0)
metaLoop:
	for i := range tests {
		input, err := readFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := readFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		for i := range input {
			if input[i] == '!' {
				input[i] = 0
			}
		}
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

func test_all_others() {
	interpreters, err := filepath.Glob("../Other/*.bf")
	if err != nil {
		panic(err)
	}
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
	compileds := make([]Segment, len(interpreters))
	seperatorsymbol := make([]byte, len(interpreters))
	names := make([][]byte, len(interpreters))
	for i := range compileds {
		interpreter, err := readFile(interpreters[i])
		if err != nil {
			panic(err)
		}
		compileds[i] = Compile(interpreter)
		compileds[i] = Optimize(compileds[i])
		symbol, err := readFile(interpreters[i] + ".sym")
		if err != nil {
			seperatorsymbol[i] = 0
		} else {
			seperatorsymbol[i] = symbol[0]
		}
		names[i], err = readFile(interpreters[i] + ".name")
		if err != nil {
			panic(err)
		}
		names[i] = append(names[i], []byte(fmt.Sprintf(" (%d)", i))...)
	}
	c := make(chan bool, len(compileds))

	for i := range tests {
		input, err := readFile(tests[i])
		if err != nil {
			panic(err)
		}
		output, err := readFile(tests[i][:len(tests[i])-2] + "out")
		if err != nil {
			panic(err)
		}
		fmt.Printf("----------------\nTesting %s\n----------------\n", tests[i])
		quit := false
		for j := range compileds {
			sepInput := make([]byte, len(input))
			copy(sepInput, input)
			go test_other_quick(
				compileds[j],
				sepInput,
				output,
				seperatorsymbol[j],
				names[j],
				c,
				&quit,
			)
		}
		go func() {
			time.Sleep(time.Minute * 15)
			quit = true
		}()
		for _ = range compileds {
			<-c
		}
	}
}

func test_other_quick(
	compiled Segment,
	input []byte,
	output []byte,
	sep byte,
	name []byte,
	c chan bool,
	quit *bool,
) {
	for i := range input {
		if input[i] == '!' {
			input[i] = sep
		}
	}
	data := MakeData(append(input, 0))
	err := compiled.RunWithTimeout(data, quit)
	if err != nil {
		fmt.Printf("%s failed\n\t%s\n", name, err)
		c <- false
		return
	}
	bufferbytes := data.output.Bytes()
	if len(bufferbytes) != len(output) {
		fmt.Printf(
			"%s failed\n\t(got \"%s\", wanted \"%s\")\n",
			name,
			string(bufferbytes),
			string(output),
		)
		c <- false
	} else {
		for i := range bufferbytes {
			if bufferbytes[i] != output[i] {
				fmt.Printf(
					"%s failed\n\t(got \"%s\", wanted \"%s\")\n",
					name,
					string(bufferbytes),
					string(output),
				)
				c <- false
				return
			}
		}
		fmt.Printf("%s succeeded\n", name)
		c <- true
	}
}
