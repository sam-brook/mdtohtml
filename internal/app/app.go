package app

import (
	"fmt"

	"github.com/sam-brook/mdtohtml/internal/parser"
)

func Run() {
	testing_data := "```\n# **Title potato**\n```\n\n## *~~hello potato~~*"
	result := parser.Parse(testing_data)
	expected := "<h1>Title potato</h1>\n<h2>hello potato<h2>"
	fmt.Println("Test passed, expected result was\nexpected: ", expected, "\nresult:", result)
}
