// test_interpreter.go
package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestInterpret(t *testing.T) {
	tests, err := filepath.Glob("../Tests/*.in")
	if err != nil {
		panic(err)
	}
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
		inputs := strings.SplitN(string(input), "!", 2)
		if len(inputs) == 1 {
			Interpret(input, []byte{}, &buffer)
		} else {
			inputs[1] = strings.Replace(inputs[1], "!", string(0), -1)
			Interpret([]byte(inputs[0]), []byte(inputs[1]), &buffer)
		}
		bufferbytes := buffer.Bytes()
		if len(bufferbytes) != len(output) {
			t.Errorf(
				"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
				tests[i],
				string(bufferbytes),
				string(output),
			)
		} else {
			for i := range bufferbytes {
				if bufferbytes[i] != output[i] {
					t.Errorf(
						"Failed %s\n\t(got \"%s\", wanted \"%s\")\n",
						tests[i],
						string(bufferbytes),
						string(output),
					)
					break
				}
			}
			t.Logf("Succeeded %s\n", tests[i])
		}
	}
}
