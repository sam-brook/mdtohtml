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
}

type Parser struct {
	Lines          []string
	LineIndex      int
	CharIndex      int
	MultiLineTags  SyntaxStack
	BlockLevelTags SyntaxStack
	InTextTags     SyntaxStack
	Output         strings.Builder
}

func (s *SyntaxStack) Push(val string) {
	s.topIndex++
	if len(s.syntax) > s.topIndex+1 {
		s.syntax = append(s.syntax, val)
	} else {
		s.syntax[s.topIndex] = val
	}
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
	return blockhtmlTags[s.syntax[s.topIndex]]
}

func (s *SyntaxStack) Clear() {
	s.topIndex = -1
}

var blockhtmlTags = map[string]HTMLElement{
	"# ":      {"<h1>", "</h1>"},
	"## ":     {"<h2>", "</h2>"},
	"### ":    {"<h3>", "</h3>"},
	"#### ":   {"<h4>", "</h4>"},
	"##### ":  {"<h5>", "</h5>"},
	"###### ": {"<h6>", "</h6>"},
	"> ":      {"<blockquote>", "</blockquote>"},
	"- ":      {"<ul><li>", "</li></ul>"},
}

var multiLineTags = map[string]HTMLElement{
	"```": {"<code>", "</code>"},
}

// TODO add links
var inTextTags = map[string]HTMLElement{
	"*":  {"<em>", "</em>"},
	"**": {"<strong>", "</strong>"},
	"`":  {"<code>", "</code>"},
	"~~": {"<del>", "</del>"},
}

func NewParser(input *string) *Parser {
	return &Parser{
		Lines:     strings.Split(*input, "\n"),
		LineIndex: 0,
		CharIndex: 0,
		BlockLevelTags: SyntaxStack{
			syntax:   make([]string, 1),
			topIndex: -1,
		},
		MultiLineTags: SyntaxStack{
			syntax:   make([]string, 1),
			topIndex: -1,
		},
		InTextTags: SyntaxStack{
			syntax:   make([]string, 10),
			topIndex: -1,
		},
	}
}

func GetPrefixToken(input string) string {
	first_space_index := strings.IndexAny(input, " ")
	return input[0 : first_space_index+1]
}

func WrapLineWithTag(p *Parser, line string) {
	element := blockhtmlTags[p.BlockLevelTags.Pop()]
	p.Output.WriteString(element.Open)
	p.Output.WriteString(line)
	p.Output.WriteString(element.Close)
	p.Output.WriteString("\n")
}

func stripBlockTag(current_line string) string {
	return current_line[strings.IndexAny(current_line, " ")+1:]
}

func Parse(input string) string {
	p := NewParser(&input)

	fmt.Println("test input value: ", input)

	for ; p.LineIndex < len(p.Lines); p.LineIndex++ {
		current_line := p.Lines[p.LineIndex]
		first_word := GetPrefixToken(current_line)
		_, ok := blockhtmlTags[first_word]
		if ok {
			p.BlockLevelTags.Push(first_word)
			fmt.Println("test stack peek", p.BlockLevelTags.Peek())
		}
		WrapLineWithTag(p, stripBlockTag(current_line))

	}

	return p.Output.String()
}
