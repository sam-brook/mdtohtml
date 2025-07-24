package parser

import (
	"strings"
)

type HtmlElement struct {
	Open  string
	Close string
}

type MarkdownTag struct {
	StartIndex int
	EndIndex   int
	Element    HtmlElement
}

type Parser struct {
	Lines       []string
	LineIndex   int
	CharIndex   int
	SyntaxStack []MarkdownTag
	Output      strings.Builder
}

var HtmlElements = map[string]HtmlElement{
	"# ":     {"<h1>", "</h1>"},
	"## ":    {"<h2>", "</h2>"},
	"h3":     {"<h3>", "</h3>"},
	"h4":     {"<h4>", "</h4>"},
	"h5":     {"<h5>", "</h5>"},
	"h6":     {"<h6>", "</h6>"},
	"em":     {"<em>", "</em>"},
	"strong": {"<strong>", "</strong>"},
	"code":   {"<code>", "</code>"},
	"del":    {"<del>", "</del>"},
	"ul":     {"<ul><li>", "</li></ul>"},
}
