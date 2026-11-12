package solution

// IsValid determines if the string of brackets is valid using a stack.
func IsValid(s string) bool {
	stack := make([]rune, 0, len(s))
	match := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, ch := range s {
		switch ch {
		case '(', '{', '[':
			stack = append(stack, ch)
		case ')', '}', ']':
			if len(stack) == 0 || stack[len(stack)-1] != match[ch] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
