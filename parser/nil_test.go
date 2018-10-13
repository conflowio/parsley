package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Nil", func() {

	var (
		p   = parser.Nil()
		ctx *parsley.Context
	)

	BeforeEach(func() {
		ctx = parsley.NewContext(parsley.NewFileSet(), &parsleyfakes.FakeReader{})
	})

	It("should return with an empty node", func() {
		res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, 1)
		Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
		Expect(res).To(Equal(ast.NilNode(1)))
		Expect(err).ToNot(HaveOccurred())
	})

})
