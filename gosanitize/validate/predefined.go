package validate

type PatternMatch map[string]string

/* Custom match patterns, email is not a proper validation pattern, just for testing. */
var Patterns = PatternMatch{
	"Email": "^([^ @]+@[^\\\\.]+\\\\.[^ ]+)$",
}

func AddPattern(id string, pattern string) {
	Patterns[id] = pattern
}
