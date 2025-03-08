package main

import (
	"testing"
)

func TestNewConverter(t *testing.T) {
	// given
	input := "test string"

	// when
	converter := NewConverter(input)

	// then
	if converter == nil {
		t.Error("Expected non-nil Converter")
	}
	if len(converter.input) != len(input) {
		t.Errorf("Expected input length %d, got %d", len(input), len(converter.input))
	}
	if converter.pos != 0 {
		t.Errorf("Expected initial position 0, got %d", converter.pos)
	}
}

func TestConverter_StackOperations(t *testing.T) {
	// given
	c := NewConverter("")

	// then
	c.push("div")
	if len(c.stack) != 1 || c.stack[0] != "div" {
		t.Error("Push failed")
	}

	tag := c.peek()
	if tag != "div" {
		t.Errorf("Expected peek to return 'div', got '%s'", tag)
	}

	tag, ok := c.pop()
	if !ok || tag != "div" {
		t.Error("Pop failed")
	}

	// Test empty stack
	_, ok = c.pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}
}

func TestConverter_Write(t *testing.T) {
	testCases := []struct {
		input    rune
		expected string
	}{
		{'<', "&lt;"},
		{'>', "&gt;"},
		{'&', "&amp;"},
		{'ф', "ф"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.input), func(t *testing.T) {
			// given
			c := NewConverter("")

			// when
			c.write(tc.input)

			// then
			if c.output.String() != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, c.output.String())
			}
		})
	}
}

func TestConverter_CheckAhead(t *testing.T) {
	testCases := []struct {
		input    string
		pattern  string
		pos      int
		expected bool
	}{
		{"**bold**", "**", 0, true},
		{"**bold**", "**", 1, false},
		{"*italic*", "**", 0, false},
		{"```code```", "```", 0, true},
		{"", "**", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// given
			c := NewConverter(tc.input)
			c.pos = tc.pos

			// when
			result := c.checkAhead(tc.pattern)

			// then
			if result != tc.expected {
				t.Errorf("Expected %v for input '%s' at pos %d with pattern '%s', got %v",
					tc.expected, tc.input, tc.pos, tc.pattern, result)
			}
		})
	}
}

func TestConverter_ProcessLink(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"[text](http://example.com)",
			"<a href=\"http://example.com\">text</a>",
		},
		{
			"[text](invalid)",
			"",
		},
		{
			"[text]invalid",
			"",
		},
		{
			"[text]()",
			"",
		},
		{
			"[](http://example.com)",
			"<a href=\"http://example.com\"></a>",
		},
		{
			"[text with <special> chars](http://example.com)",
			"<a href=\"http://example.com\">text with &lt;special&gt; chars</a>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// given
			c := NewConverter(tc.input)

			// when
			result := c.processLink()

			// then
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestConverter_ConvertSpecialCases(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Unclosed tags",
			"**bold**",
			"<b>bold</b>",
		},
		{
			"Nested formatting",
			"**bold *italic* text**",
			"<b>bold <i>italic</i> text</b>",
		},
		{
			"Escape sequences",
			"\\* \\` \\[ \\\\",
			"* ` [ \\",
		},
		{
			"Empty input",
			"",
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			c := NewConverter(tc.input)

			// when
			result := c.Convert()

			// then
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}
