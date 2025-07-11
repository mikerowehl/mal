package main

import (
	"bufio"
	"fmt"
	"os"
)

// Returned values are:
//
//	string - token/line
//	bool - eof, true means end of input
//	error - set to nil unless there's an error
func READ(scanner *bufio.Scanner) (string, bool, error) {
	fmt.Print("user> ")

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", false, fmt.Errorf("READ unable to get line: %w", err)
		} else {
			return "", true, nil
		}
	}
	return scanner.Text(), false, nil
}

func EVAL(s string) string {
	return s
}

func PRINT(s string) {
	fmt.Println(s)
}

func rep() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		line, eof, err := READ(scanner)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			return
		}
		if eof {
			fmt.Println()
			return
		}
		PRINT(EVAL(line))
	}
}

func main() {
	rep()
}
