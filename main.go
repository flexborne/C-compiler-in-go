package main

import (
	"C_compiler/lexer"
	"C_compiler/parser"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadContentFromFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}
	return string(content), nil
}

func main() {
	fileName := "..//to_compile.c"
	codeInC, err := ReadContentFromFile(fileName)
	if err != nil {
		return
	}
	lex := lexer.NewLexer(codeInC)

	function, err := parser.ParseFunction(lex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(function)

	instructionMapping := map[string]string{
		"printf": "CALL puts",
	}

	var machineCode []string
	machineCode = append(machineCode, "main:", "push %rbp", "mov %rbp, %rsp")
	var LCLabels []string
	labelIdx := 0

	for _, statement := range function.Body {
		switch statement.(type) {
		case parser.FuncCallStatement:
			funcCall := statement.(parser.FuncCallStatement)
			if instruction, found := instructionMapping[funcCall.Name]; found {
				for _, arg := range funcCall.Args {
					LCLabels = append(LCLabels, ".LC"+strconv.Itoa(labelIdx)+"\n  .string "+"\""+arg+"\"")
					machineCode = append(machineCode, "mov %edi, OFFSET FLAT:.LC"+strconv.Itoa(labelIdx))
					labelIdx++
				}
				machineCode = append(machineCode, instruction)
			} else {
				machineCode = append(machineCode, "UNSUPPORTED_FUNCTION_CALL")
			}
		case parser.ReturnStatement:
			returnStatement := statement.(parser.ReturnStatement)
			machineCode = append(machineCode, "mov %eax, "+returnStatement.Expr, "pop %rbp", "ret")
		}
	}
	machineCode = append(LCLabels, machineCode...)
	err = os.WriteFile("assembly.s", []byte(strings.Join(machineCode, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing assembly code to file:", err)
		return
	}
}
