package parser

import (
	"fmt"
	"strings"
)

type HTMLElement struct {
	Open  string
	Close string
}

type SyntaxStack struct {
	syntax   []string
	topIndex int
}

type Stack interface {
	Push(val string)
	Pop() string
	Peek() string
	GetElement() HTMLElement
	Clear()
}

type Parser struct {
	Lines       []string
	LineIndex   int
	CharIndex   int
	Syntaxstack SyntaxStack
	Output      strings.Builder
}

func (s *SyntaxStack) Push(val string) {
	s.syntax[s.topIndex] = val
}

func (s *SyntaxStack) Pop() string {
	popped_elem := s.syntax[s.topIndex]
	s.topIndex--
	return popped_elem
}

func (s *SyntaxStack) Peek() string {
	popped_elem := s.syntax[s.topIndex]
	return popped_elem
}

func (s *SyntaxStack) GetElement() HTMLElement {
	return htmlTags[s.syntax[s.topIndex]]
}

func (s SyntaxStack) Clear() {
	s.topIndex = 0
}

var htmlTags = map[string]HTMLElement{
	"# ":      {"<h1>", "</h1>"},
	"## ":     {"<h2>", "</h2>"},
	"### ":    {"<h3>", "</h3>"},
	"#### ":   {"<h4>", "</h4>"},
	"##### ":  {"<h5>", "</h5>"},
	"###### ": {"<h6>", "</h6>"},
	"*":       {"<em>", "</em>"},
	"**":      {"<strong>", "</strong>"},
	"`":       {"<code>", "</code>"},
	"~~":      {"<del>", "</del>"},
	"- ":      {"<ul><li>", "</li></ul>"},
}

func NewParser(input *string) *Parser {
	return &Parser{
		Lines:     strings.Split(*input, "\n"),
		LineIndex: 0,
		CharIndex: 0,
		Syntaxstack: SyntaxStack{
			syntax:   make([]string, 10),
			topIndex: 0,
		},
	}
}

func GetPrefixToken(input string) string {
	first_space_index := strings.IndexAny(input, " ")
	return input[0 : first_space_index+1]
}

func WrapWithTag(p *Parser, line string) {
	element := htmlTags[p.Syntaxstack.Pop()]
	p.Output.WriteString(element.Open)
	p.Output.WriteString(line)
	p.Output.WriteString(element.Close)
	p.Output.WriteString("\n")
}

func Parse(input string) string {
	p := NewParser(&input)

	fmt.Println("test input value: ", input)

	for ; p.LineIndex < len(p.Lines); p.LineIndex++ {
		current_line := p.Lines[p.LineIndex]
		fmt.Println("test current_line value:", current_line)
		first_word := GetPrefixToken(current_line)
		fmt.Println("test first_word value:", first_word)
		_, ok := htmlTags[first_word]
		if ok {
			p.Syntaxstack.Push(first_word)
			fmt.Println("test stack peek", p.Syntaxstack.Peek())
		}
		WrapWithTag(p, current_line[strings.IndexAny(current_line, " ")+1:])

	}

	return p.Output.String()
}
