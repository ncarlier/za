package helper

// Match finds whether the string matches the pattern.
// Supports '*' and '?' wildcards in the pattern string.
func Match(pattern, str string) bool {
	if pattern == "" {
		return str == pattern
	}
	if pattern == "*" {
		return true
	}
	return deepMatchRune([]rune(str), []rune(pattern))
}

func deepMatchRune(str, pattern []rune) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		default:
			if len(str) == 0 || str[0] != pattern[0] {
				return false
			}
		case '?':
			if len(str) == 0 {
				return false
			}
		case '*':
			return deepMatchRune(str, pattern[1:]) ||
				(len(str) > 0 && deepMatchRune(str[1:], pattern))
		}
		str = str[1:]
		pattern = pattern[1:]
	}
	return len(pattern) == 0
}
