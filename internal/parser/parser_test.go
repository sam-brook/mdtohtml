package parser

import (
	"fmt"
	"testing"
)

func Test_IsTagChar(t *testing.T) {
	test_data_expect_true := []byte{'*', '~', '`'}
	test_data_expect_false := []byte{'c', '#', '&'}

	for _, tag_chars := range test_data_expect_true {
		if !isTagChar(tag_chars) {
			t.Errorf("expected true and is returning false for char %c", tag_chars)
		} else {
			fmt.Printf("test passed: %c\n", tag_chars)
		}
	}

	for _, non_tag_chars := range test_data_expect_false {
		if isTagChar(non_tag_chars) {
			t.Errorf("expected false and is returning true for char %c", non_tag_chars)
		} else {
			fmt.Printf("test passed: %c\n", non_tag_chars)
		}

	}
}

func Test(t *testing.T) {
	testing_data := "# Title potato\n## hello potato"
	result := Parse(testing_data)
	expected := "<h1>Title potato</h1>\n<h2>hello potato</h2>\n"

	if result != expected {
		t.Errorf("expected: %q\nresult: %q", expected, result)
	} else {
		fmt.Println("Test passed, expected result was\nexpected:", expected, "\nresult:", result)
	}

	testing_data = "# **Title potato**\n## hello potato"
	result = Parse(testing_data)
	expected = "<h1><strong>Title potato</strong></h1>\n<h2>hello potato</h2>\n"

	if result != expected {
		t.Errorf("expected: %q\nresult: %q", expected, result)
	} else {
		fmt.Println("Test passed, expected result was\nexpected:", expected, "\nresult:", result)
	}
}
