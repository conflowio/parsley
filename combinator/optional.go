package combinator

import (
	"github.com/opsidian/parsley/parser"
)

// Optional returns the parser's matches or an empty match if it fails
func Optional(p parser.Parser) parser.Parser {
	return Choice("", p, parser.Empty())
}
