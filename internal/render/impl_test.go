package render

import (
	"strings"
	"testing"
)

func TestAllFormats(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Bold text",
			input:    "**Bold Text**",
			expected: "*Bold Text*",
		},
		{
			name:     "Italic text",
			input:    "*Italic Text*",
			expected: "_Italic Text_",
		},
		{
			name:     "Bold and Italic",
			input:    "***Bold and Italic***",
			expected: "*_Bold and Italic_*",
		},
		{
			name:     "Strikethrough",
			input:    "~~Strikethrough~~",
			expected: "~Strikethrough~",
		},
		{
			name:     "Inline code",
			input:    "`inline code`",
			expected: "`inline code`",
		},
		{
			name:     "Heading 1",
			input:    "# Heading 1",
			expected: "*Heading 1*",
		},
		{
			name:     "Heading 2",
			input:    "## Heading 2",
			expected: "*Heading 2*",
		},
		{
			name:     "Heading 3",
			input:    "### Heading 3",
			expected: "*Heading 3*",
		},
		{
			name:  "Bullet points",
			input: "- Bullet point\n- Another bullet style",
			expected: `• Bullet point
• Another bullet style`,
		},
		{
			name:  "Numbered list",
			input: "1. Numbered list\n2. Second item",
			expected: `1\. Numbered list
2\. Second item`,
		},
		{
			name:     "Blockquote",
			input:    "> Blockquote",
			expected: "> Blockquote",
		},
		{
			name:     "Link",
			input:    "[Hyperlink text](URL)",
			expected: "[Hyperlink text](URL)",
		},
		{
			name:     "Horizontal line",
			input:    "---",
			expected: "---",
		},
		{
			name:     "Code block",
			input:    "```python\nprint('hello')\n```",
			expected: "```python\nprint('hello')\n```",
		},
		{
			name:     "Code block with language",
			input:    "```go\nfmt.Println(\"hello\")\n```",
			expected: "```go\nfmt.Println(\"hello\")\n```",
		},
		{
			name:     "Special characters escaping",
			input:    "Text with dots... and other chars: ! # + - = | { }",
			expected: "Text with dots\\.\\.\\. and other chars: \\! \\# \\+ \\- \\= \\| \\{ \\}",
		},
		{
			name:     "Parentheses and brackets escaping",
			input:    "Version 2.0 (latest) [stable]",
			expected: "Version 2\\.0 \\(latest\\) \\[stable\\]",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			adjusted, err := AdjustMdToTelegramFormat(tc.input)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}
			result := strings.TrimSpace(adjusted)
			expected := strings.TrimSpace(tc.expected)

			if result != expected {
				t.Errorf("\nInput:    %q\nExpected: %q\nGot:      %q", tc.input, expected, result)
			}
		})
	}
}
