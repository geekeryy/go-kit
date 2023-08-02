package util

func MaybeRemoveQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	switch s[0] {
	case '"':
		if s[len(s)-1] != '"' {
			return s
		}
	case '\'':
		if s[len(s)-1] != '\'' {
			return s
		}
	default:
		return s
	}
	return s[1 : len(s)-1]
}
