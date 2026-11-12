package validparentheses

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "simple pair",
			input:    "()",
			expected: true,
		},
		{
			name:     "multiple pairs",
			input:    "()[]{}",
			expected: true,
		},
		{
			name:     "mismatched pair",
			input:    "(]",
			expected: false,
		},
		{
			name:     "nested matching",
			input:    "([{}])",
			expected: true,
		},
		{
			name:     "unmatched open bracket",
			input:    "([",
			expected: false,
		},
		{
			name:     "unmatched close bracket",
			input:    "]",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "incorrect ordering",
			input:    "([)]",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.input)
			if got != tt.expected {
				t.Errorf("IsValid(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func BenchmarkIsValid(b *testing.B) {
	s := "((([[{{}}]])))((([[[[]]]])))"
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = IsValid(s)
	}
}
