package mal

import (
	"fmt"
	"regexp"
)

type MalType any
type MalList []any

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

func read_atom(reader *Reader) MalType {
	next := reader.Next()
	return next
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

func read_form(reader *Reader) (MalType, error) {
	if reader.Done() {
		return nil, fmt.Errorf("read_form: Missing input")
	}
	next := reader.Peek()
	if next == "(" {
		reader.Next()
		return read_list(reader)
	}
	return read_atom(reader), nil
}

func Read_str(s string) (MalType, error) {
	reader := NewReader(s)
	return read_form(reader)
}

func Pr_str(o MalType) {
	switch t := o.(type) {
	case string:
		fmt.Print(t)
	case MalList:
		fmt.Print("(")
		sep := ""
		for _, sub := range t {
			if sep != "" {
				fmt.Print(sep)
			}
			Pr_str(sub)
			sep = " "
		}
		fmt.Print(")")
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
