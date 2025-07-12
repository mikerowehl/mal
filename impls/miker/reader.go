package mal

import (
	"fmt"
	"regexp"
	"strings"
)

type MalType any
type MalSymbol string
type MalList []any
type MalVector []any

type Reader struct {
	Tokens  []string
	Current int
}

func (r *Reader) Peek() string {
	return r.Tokens[r.Current]
}

func (r *Reader) Next() string {
	s := r.Tokens[r.Current]
	r.Current += 1
	return s
}

func (r *Reader) Done() bool {
	return r.Current >= len(r.Tokens)
}

func NewReader(s string) *Reader {
	tokens := tokenize(s)
	return &Reader{
		Tokens:  tokens,
		Current: 0,
	}
}

func tokenize(s string) []string {
	var validToken = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)
	matches := validToken.FindAllStringSubmatch(s, -1)
	tokens := make([]string, len(matches))
	for i := range matches {
		tokens[i] = matches[i][1]
	}
	return tokens
}

func escapestr(s string) string {
	var builder strings.Builder
	for _, r := range s {
		switch r {
		case '\n':
			builder.WriteRune('\\')
			builder.WriteRune('n')
		case '\\':
			builder.WriteRune('\\')
			builder.WriteRune('\\')
		case '"':
			builder.WriteRune('\\')
			builder.WriteRune('"')
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func process_string(s string) (string, error) {
	var builder strings.Builder
	inBackslash := false
	for _, r := range s {
		if inBackslash {
			switch r {
			case 'n':
				builder.WriteRune('\n')
			case '\\', '"':
				builder.WriteRune(r)
			default:
				return "", fmt.Errorf("process_string: unbalanced backslash")
			}
			inBackslash = false
		} else {
			if r == '\\' {
				inBackslash = true
			} else {
				builder.WriteRune(r)
			}
		}
	}
	if inBackslash {
		return "", fmt.Errorf("process_string: unbalanced backslash")
	}
	return builder.String(), nil
}

func read_atom(reader *Reader) (MalType, error) {
	next := reader.Next()
	if len(next) == 0 {
		return nil, fmt.Errorf("read_atom: empty token")
	}
	if next[0] == '"' {
		if len(next) < 2 || !strings.HasSuffix(next, `"`) {
			return nil, fmt.Errorf("read_atom: unbalanced quotes")
		}
		// We know it had a quote at the start and end, strip them and process
		return process_string(next[1 : len(next)-1])
	}
	if next[0] == ':' {
		keyword := next[1:]
		runeSlice := []rune(keyword)
		newRuneSlice := make([]rune, len(runeSlice)+1)
		newRuneSlice[0] = rune(0x29e)
		copy(newRuneSlice[1:], runeSlice)
		return string(newRuneSlice), nil
	}
	return MalSymbol(next), nil
}

func read_list(reader *Reader) (MalList, error) {
	var list []any
	for {
		if reader.Done() {
			return nil, fmt.Errorf("read_list: EOF")
		}
		next := reader.Peek()
		if next == ")" {
			reader.Next()
			return list, nil
		}
		val, err := read_form(reader)
		if err != nil {
			return nil, fmt.Errorf("read_list: Error in subform: %w", err)
		}
		list = append(list, val)
	}
}

func read_vector(reader *Reader) (MalVector, error) {
	var vector []any
	for {
		if reader.Done() {
			return nil, fmt.Errorf("read_vector: EOF")
		}
		next := reader.Peek()
		if next == "]" {
			reader.Next()
			return vector, nil
		}
		val, err := read_form(reader)
		if err != nil {
			return nil, fmt.Errorf("read_vector: Error in subform: %w", err)
		}
		vector = append(vector, val)
	}
}

func read_form(reader *Reader) (MalType, error) {
	if reader.Done() {
		return nil, fmt.Errorf("read_form: Missing input")
	}
	next := reader.Peek()
	if next == "(" {
		reader.Next()
		return read_list(reader)
	}
	if next == "[" {
		reader.Next()
		return read_vector(reader)
	}
	if next == "'" {
		reader.Next()
		quoted, err := read_form(reader)
		if err != nil {
			return nil, fmt.Errorf("read_form: quoted error %w", err)
		}
		return MalList{MalSymbol("quote"), quoted}, nil
	}
	return read_atom(reader)
}

func Read_str(s string) (MalType, error) {
	reader := NewReader(s)
	return read_form(reader)
}

func Pr_str(o MalType, readably bool) {
	switch t := o.(type) {
	case string:
		if len(t) > 0 {
			runes := []rune(t)
			if runes[0] == rune(0x29e) {
				fmt.Print(":")
				fmt.Print(string(runes[1:]))
			} else {
				if readably {
					fmt.Print(`"` + escapestr(t) + `"`)
				} else {
					fmt.Print(`"` + t + `"`)
				}
			}
		} else {
			fmt.Print(`""`)
		}
	case MalSymbol:
		fmt.Print(t)
	case MalList:
		fmt.Print("(")
		sep := ""
		for _, sub := range t {
			fmt.Print(sep)
			Pr_str(sub, readably)
			sep = " "
		}
		fmt.Print(")")
	case MalVector:
		fmt.Print("[")
		sep := ""
		for _, sub := range t {
			fmt.Print(sep)
			Pr_str(sub, readably)
			sep = " "
		}
		fmt.Print("]")
	}
}

/*
func main() {
	// var validToken = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)

	// fmt.Printf("%q\n", validToken.FindAllStringSubmatch("(+ 2 (* 3 4))", -1))
	// r := NewReader("(+ 2 (* 3 4))")
	// for s, ok := r.Peek(); ok == true; s, ok = r.Peek() {
	// 	fmt.Println(s)
	// 	r.Next()
	// }
	testVal := read_str(" (   +   2 3 (/ 4   9 ) )")
	pr_str(testVal)
}
*/
