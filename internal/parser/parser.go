package parser

import (
	"fmt"
	"strings"
)

type HtmlElement struct {
	Open  string
	Close string
}

type SyntaxStack struct {
	MarkdownSyntax []string
	topIndex       int
}

type Stack interface {
	Push()
	Pop()
	Peek()
	GetElement()
}

type Parser struct {
	Lines       []string
	LineIndex   int
	CharIndex   int
	Syntaxstack SyntaxStack
	Output      strings.Builder
}

func (s SyntaxStack) Push(new string) {
	s.MarkdownSyntax[s.topIndex] = new
	s.topIndex++
}

func (s SyntaxStack) Pop() string {
	popped_elem := s.MarkdownSyntax[s.topIndex]
	s.topIndex--
	return popped_elem
}

func (s SyntaxStack) Peek() string {
	popped_elem := s.MarkdownSyntax[s.topIndex]
	return popped_elem
}

func (s SyntaxStack) GetElement() HtmlElement {
	return HtmlElements[s.MarkdownSyntax[s.topIndex]]
}

func (s SyntaxStack) Clear() {
	s.topIndex = 0
}

var HtmlElements = map[string]HtmlElement{
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
			MarkdownSyntax: make([]string, 10),
			topIndex:       0,
		},
	}
}

func GetFirstWord(input string) string {
	first_space_index := strings.IndexAny(input, " ")
	return input[0 : first_space_index+1]
}

func WrapLine(p *Parser, line string) {
	element := HtmlElements[p.Syntaxstack.Pop()]
	p.Output.WriteString(element.Open)
	p.Output.WriteString(line)
	p.Output.WriteString(element.Close)
	p.Output.WriteString("\n")
}

func ParseMarkdown(input string) string {
	p := NewParser(&input)

	fmt.Println("test input value: ", input)

	for ; p.LineIndex < len(p.Lines); p.LineIndex++ {
		current_line := p.Lines[p.LineIndex]
		fmt.Println("test current_line value:", current_line)
		first_word := GetFirstWord(current_line)
		fmt.Println("test first_word value:", first_word)
		_, ok := HtmlElements[first_word]
		if ok {
			p.Syntaxstack.Push(first_word)
			fmt.Println("test stack peek", p.Syntaxstack.Peek())
		}
		WrapLine(p, current_line[strings.IndexAny(current_line, " ")+1:])

	}

	return p.Output.String()
}
