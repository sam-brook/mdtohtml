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
	syntax     []string
	topIndex   int
	stack_type map[string]HTMLElement
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
	if len(s.syntax) > s.topIndex {
		s.syntax[s.topIndex] = val
	} else {
		s.syntax = append(s.syntax, val)
	}
}

func (s *SyntaxStack) Pop() string {
	if s.topIndex == -1 {
		return ""
	}
	popped_elem := s.syntax[s.topIndex]
	s.topIndex--
	return popped_elem
}

func (s *SyntaxStack) Peek() string {
	if s.topIndex == -1 {
		return ""
	}
	popped_elem := s.syntax[s.topIndex]
	return popped_elem
}

func (s *SyntaxStack) GetElement() HTMLElement {
	key := s.Peek()
	html_tag, ok := s.stack_type[key]
	if !ok {
		return HTMLElement{}
	}

	return html_tag
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
			syntax:     make([]string, 0),
			topIndex:   -1,
			stack_type: blockhtmlTags,
		},
		MultiLineTags: SyntaxStack{
			syntax:     make([]string, 0),
			topIndex:   -1,
			stack_type: multiLineTags,
		},
		InTextTags: SyntaxStack{
			syntax:     make([]string, 0),
			topIndex:   -1,
			stack_type: inTextTags,
		},
	}
}

func GetPrefixToken(input string) string {
	first_space_index := strings.IndexAny(input, " ")
	return input[0 : first_space_index+1]
}

func concatenate_bytes(a, b byte) string {
	return string([]byte{a, b})
}

func isTagChar(char byte) bool {
	switch char {
	case '*', '~', '`':
		return true
	default:
		return false
	}
}

func WriteTagPrefix(stack *SyntaxStack, md_tag string, output *strings.Builder) {
	stack.Push(md_tag)
	html_tag := stack.GetElement()
	output.WriteString(html_tag.Open)
}

func WriteTagSuffix(stack *SyntaxStack, output *strings.Builder) {
	html_tag := stack.GetElement()
	stack.Pop()
	output.WriteString(html_tag.Close)
}

func Parse(input string) string {
	p := NewParser(&input)

	fmt.Println("test input value: ", input)

	// Line level loop
	for ; p.LineIndex < len(p.Lines); p.LineIndex++ {
		current_line := p.Lines[p.LineIndex]

		first_word := GetPrefixToken(current_line)
		_, block_tag_ok := blockhtmlTags[first_word]
		if block_tag_ok {
			WriteTagPrefix(&p.BlockLevelTags, first_word, &p.Output)
			p.CharIndex = p.CharIndex + len(first_word)
		}

		// Char level loop
		for ; p.CharIndex < len(current_line); p.CharIndex++ {
			c := current_line[p.CharIndex]
			if isTagChar(c) {
				if ((p.CharIndex + 1) < len(current_line)) && (c == current_line[p.CharIndex+1]) {
					tag := concatenate_bytes(c, current_line[p.CharIndex+1])
					fmt.Println("current char is: " + string(c))
					fmt.Println("next char is: " + string(current_line[p.CharIndex+1]))
					fmt.Println("current tag is: " + tag)

					if p.InTextTags.Peek() == tag {
						WriteTagSuffix(&p.InTextTags, &p.Output)
					} else {
						WriteTagPrefix(&p.InTextTags, tag, &p.Output)
					}
					p.CharIndex = p.CharIndex + 1
				} else {
					fmt.Println("current tag is: " + string(c))
					if p.InTextTags.Peek() == string(c) {
						WriteTagSuffix(&p.InTextTags, &p.Output)
					} else {
						WriteTagPrefix(&p.InTextTags, string(c), &p.Output)
					}
				}
			} else {
				p.Output.WriteByte(c)
			}
		}
		for p.InTextTags.Peek() != "" {
			WriteTagSuffix(&p.InTextTags, &p.Output)
		}
		if block_tag_ok {
			WriteTagSuffix(&p.BlockLevelTags, &p.Output)
		}

		p.Output.WriteByte('\n')
		p.CharIndex = 0

	}

	return p.Output.String()
}
