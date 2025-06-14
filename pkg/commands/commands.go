package commands

// Command represents a parsed command
type Command struct {
	Type  string
	Key   string
	Value string // Only used for SET commands
}

// CommandType constants
const (
	TypeSet = "SET"
	TypeGet = "GET"
	TypeDel = "DEL"
)

// ParseCommandType determines the type of command from the input
func ParseCommandType(input string) string {
	// Split the input by spaces to get the first word
	words := splitWords(input)
	if len(words) == 0 {
		return ""
	}

	switch words[0] {
	case TypeSet:
		return TypeSet
	case TypeGet:
		return TypeGet
	case TypeDel:
		return TypeDel
	default:
		return ""
	}
}

// splitWords splits the input string into words, preserving quoted strings
func splitWords(input string) []string {
	var words []string
	var currentWord string
	var inQuotes bool

	for i := 0; i < len(input); i++ {
		char := input[i]

		if char == '"' {
			inQuotes = !inQuotes
			continue
		}

		if char == ' ' && !inQuotes {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
			continue
		}

		currentWord += string(char)
	}

	if currentWord != "" {
		words = append(words, currentWord)
	}

	return words
}
