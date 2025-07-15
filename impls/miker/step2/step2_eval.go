package main

import (
	"bufio"
	"fmt"
	"mal"
	"os"
)

// Returned values are:
//
//	string - token/line
//	bool - eof, true means end of input
//	error - set to nil unless there's an error
func READ(scanner *bufio.Scanner) (mal.MalType, bool, error) {
	fmt.Print("user> ")

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", false, fmt.Errorf("READ unable to get line: %w", err)
		} else {
			return "", true, nil
		}
	}
	val, err := mal.Read_str(scanner.Text())
	return val, false, err
}

func APPLY(v mal.MalType) (mal.MalType, error) {
	if l, ok := v.(mal.MalList); !ok {
		fmt.Println("Applying something that isn't a list")
	} else {
		a, ok := l[0].(mal.MalFunc)
		if !ok {
			return nil, fmt.Errorf("Error converting apply function")
		}
		return a(l[1:])
	}
	return nil, fmt.Errorf("Error applying")
}

func listCast(raw mal.MalType, min int) (mal.MalList, error) {
	l, ok := raw.(mal.MalList)
	if !ok {
		return nil, fmt.Errorf("listCast: expected list: %v", raw)
	}
	if len(l) < min {
		return nil, fmt.Errorf("listCast: expected at least %d elements, got %d", min, len(l))
	}
	return l, nil
}

var add = func(a mal.MalType) (mal.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int first arg to add: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int second arg to add: %v", l[1])
	}
	return mal.MalType(i1 + i2), nil
}

var sub = func(a mal.MalType) (mal.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int second arg: %v", l[1])
	}
	return mal.MalType(i1 - i2), nil
}

var mul = func(a mal.MalType) (mal.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int second arg: %v", l[1])
	}
	return mal.MalType(i1 * i2), nil
}

var div = func(a mal.MalType) (mal.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int second arg: %v", l[1])
	}
	return mal.MalType(i1 / i2), nil
}
var env = map[string]mal.MalType{
	"+": mal.MalFunc(add),
	"-": mal.MalFunc(sub),
	"*": mal.MalFunc(mul),
	"/": mal.MalFunc(div),
}

func EVAL(v mal.MalType, env map[string]mal.MalType) mal.MalType {
	switch t := v.(type) {
	case mal.MalSymbol:
		for name, entry := range env {
			if name == string(t) {
				return entry
			}
		}
	case mal.MalList:
		if len(t) == 0 {
			return v
		}
		evaled := mal.MalList{}
		for _, entry := range t {
			n := EVAL(entry, env)
			evaled = append(evaled, n)
		}
		app, err := APPLY(evaled)
		if err != nil {
			fmt.Println("Error returned from apply")
			return nil
		}
		return app
	case mal.MalVector:
		evaled := mal.MalVector{}
		for _, entry := range t {
			n := EVAL(entry, env)
			evaled = append(evaled, n)
		}
		return evaled
	case mal.MalHashmap:
		evaled := mal.MalHashmap{}
		for i, entry := range t {
			if (i % 2) == 1 {
				n := EVAL(entry, env)
				evaled = append(evaled, n)
			} else {
				evaled = append(evaled, entry)
			}
		}
		return evaled
	}
	return v
}

func PRINT(v mal.MalType) {
	mal.Pr_str(v, true)
	fmt.Println()
}

func rep() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		line, eof, err := READ(scanner)
		if eof {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		} else {
			PRINT(EVAL(line, env))
		}
	}
}

func main() {
	rep()
}
