package parser

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	testing_data := "# Title potato\n## hello potato"
	result := ParseMarkdown(testing_data)
	expected := "<h1>Title potato</h1>\n<h2>hello potato</h2>\n"

	if result != expected {
		t.Errorf("expected: %q\nresult: %q", expected, result)
	} else {
		fmt.Println("Test passed, expected result was\nexpected:", expected, "\nresult:", result)
	}
}
