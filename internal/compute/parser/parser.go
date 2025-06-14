package parser

import (
	"context"
	"strings"
	"unicode"

	"concurrency-go/internal"

	"go.uber.org/zap"
)

// Parser represents a command parser
type Parser struct {
	logger *zap.Logger
}

// New creates a new Parser instance
func New(logger *zap.Logger) *Parser {
	return &Parser{
		logger: logger,
	}
}

// ProcessCommand implements the ComputeLayer interface
func (p *Parser) ProcessCommand(ctx context.Context, query string) (*internal.Command, error) {
	p.logger.Info("parsing command", zap.String("query", query))
	return p.parseCommand(query)
}

// parseCommand parses the input command and returns the result
func (p *Parser) parseCommand(input string) (*internal.Command, error) {
	words := strings.Fields(input)
	if len(words) == 0 {
		p.logger.Error("empty query received")
		return nil, internal.ErrInvalidCommand
	}

	cmd := &internal.Command{}

	switch words[0] {
	case "SET":
		if len(words) != 3 {
			p.logger.Error("invalid SET command syntax",
				zap.Int("words_count", len(words)),
				zap.Strings("words", words))
			return nil, internal.ErrInvalidSyntax
		}
		cmd.Type = "SET"
		cmd.Key = words[1]
		cmd.Value = words[2]

	case "GET":
		if len(words) != 2 {
			p.logger.Error("invalid GET command syntax",
				zap.Int("words_count", len(words)),
				zap.Strings("words", words))
			return nil, internal.ErrInvalidSyntax
		}
		cmd.Type = "GET"
		cmd.Key = words[1]

	case "DEL":
		if len(words) != 2 {
			p.logger.Error("invalid DEL command syntax",
				zap.Int("words_count", len(words)),
				zap.Strings("words", words))
			return nil, internal.ErrInvalidSyntax
		}
		cmd.Type = "DEL"
		cmd.Key = words[1]

	default:
		p.logger.Error("unknown command type",
			zap.String("command", words[0]))
		return nil, internal.ErrInvalidCommand
	}

	// Validate arguments
	if !isValidArgument(cmd.Key) || (cmd.Type == "SET" && !isValidArgument(cmd.Value)) {
		p.logger.Error("invalid argument format",
			zap.String("type", cmd.Type),
			zap.String("key", cmd.Key),
			zap.String("value", cmd.Value))
		return nil, internal.ErrInvalidSyntax
	}

	p.logger.Info("command parsed successfully",
		zap.String("type", cmd.Type),
		zap.String("key", cmd.Key),
		zap.String("value", cmd.Value))
	return cmd, nil
}

// isValidArgument checks if the argument matches the grammar rules
func isValidArgument(arg string) bool {
	if len(arg) == 0 {
		return false
	}

	for _, r := range arg {
		if !isValidChar(r) {
			return false
		}
	}
	return true
}

// isValidChar checks if the character is valid according to the grammar
func isValidChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || isPunctuation(r)
}

// isPunctuation checks if the character is a valid punctuation mark
func isPunctuation(r rune) bool {
	validPunctuation := map[rune]bool{
		'*': true,
		'/': true,
		'_': true,
		'.': true,
		'-': true,
	}
	return validPunctuation[r]
}
