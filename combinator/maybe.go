package combinator

import (
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/terminal"
)

// Maybe matches the given parser and also returns with a result containing an empty node
func Maybe(parser parser.Parser) parser.Parser {
	return Or(parser, terminal.Empty())
}
