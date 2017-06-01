package combinator

import (
	"github.com/opsidian/parsley/parser"
)

// Maybe matches the given parser and also returns with a result containing an empty node
func Maybe(p parser.Parser) parser.Parser {
	return Or(p, parser.Empty())
}
