package main

import (
	"bufio"
	"fmt"
	"os"

	"peihao/cs5400/minrkt"
)

func main() {
	colorRed := "\033[31m"
	colorReset := "\033[0m"
	fmt.Println("Welcome to minimalistic racket phase 1 !")
	scanner := bufio.NewScanner(os.Stdin)
	mapIdentifier := make(map[string]interface{})
	callStack := make([]map[string]interface{}, 5)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		fmt.Printf("Entered input (expression): %q.\n", line)
		tokens, err := minrkt.Tokenize(line)
		if err != nil {
			fmt.Println(colorRed, "error in toknizer phase: ", err, colorReset)
			continue
		}
		root, err := minrkt.Parse(tokens)
		if err != nil {
			fmt.Println(colorRed, "error in parser phase: ", err, colorReset)
			continue
		}
		if result, t, err := root.Eval(minrkt.Params{MapIdentifier: mapIdentifier, CallStack: callStack}); err != nil {
			fmt.Println(colorRed, "error in evaluation phase: ", err, colorReset)
		} else if t == minrkt.TYPE_FLOAT64 {
			fmt.Println("Result is: ", result.(float64))
		} else if t == minrkt.TYPE_DEFINE {
			// define statement: we don't need to do anything
			fmt.Println("define :", mapIdentifier)
		} else if t == minrkt.TYPE_NOTIFICATION {
			fmt.Println(result)
		} else {
			boolRes := result.(bool)
			var stringRes string
			if boolRes {
				stringRes = "#t"
			} else {
				stringRes = "#f"
			}
			fmt.Println("Result is: ", stringRes)
		}
	}

}
