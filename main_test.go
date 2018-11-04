package main_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a simple parser which is able to parse adding two integers.
func ExampleParse() {
	sum := ast.InterpreterFunc(func(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		nodes := node.Children()
		value0, _ := nodes[0].Value(ctx)
		value1, _ := nodes[2].Value(ctx)
		return value0.(int64) + value1.(int64), nil
	})

	p := combinator.SeqOf(
		terminal.Integer(),
		terminal.Rune('+'),
		terminal.Integer(),
	).Bind(sum)

	r := text.NewReader(text.NewFile("example.file", []byte("1+2")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	node, err := parsley.Parse(ctx, combinator.Sentence(p))
	if err != nil {
		panic(err)
	}

	value, err2 := node.Value(nil)
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Result: %d\n", value.(int64))
	// Output: Result: 3
}

// Let's define a simple parser which is able to parse adding two integers.
func ExampleEvaluate() {
	sum := ast.InterpreterFunc(func(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		nodes := node.Children()
		value0, _ := nodes[0].Value(ctx)
		value1, _ := nodes[2].Value(ctx)
		return value0.(int64) + value1.(int64), nil
	})

	p := combinator.SeqOf(
		terminal.Integer(),
		terminal.Rune('+'),
		terminal.Integer(),
	).Bind(sum)

	r := text.NewReader(text.NewFile("example.file", []byte("1+2")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, err := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %d\n", value.(int64))
	// Output: Result: 3
}

var _ = Describe("Parsley", func() {

	It("should handle direct left recursion", func() {
		input := "abbbbbbbbbbbbbbbbbbb"

		concat := ast.InterpreterFunc(func(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
			s := ""
			for _, node := range node.Children() {
				val, err := node.Value(ctx)
				if err != nil {
					return nil, err
				}
				if runeVal, ok := val.(rune); ok {
					s = s + string(runeVal)
				} else {
					s = s + val.(string)
				}
			}
			return s, nil
		})

		var p parser.Func
		p = combinator.Memoize(combinator.Any(
			combinator.SeqOf(
				&p,
				terminal.Rune('b'),
			).Bind(concat),
			terminal.Rune('a'),
		))

		f := text.NewFile("testfile", []byte(input))
		r := text.NewReader(f)
		ctx := parsley.NewContext(parsley.NewFileSet(), r)
		result, err := parsley.Evaluate(ctx, combinator.Sentence(&p), nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(input))
		Expect(ctx.CallCount()).To(Equal(298))
	})

	It("should handle highly ambiguous left-recursive grammar", func() {
		input := "1+2+3+4+5+6+7+8+9+10"

		add := ast.InterpreterFunc(func(ctx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
			nodes := node.Children()
			value0, _ := nodes[0].Value(ctx)
			value1, _ := nodes[2].Value(ctx)
			return value0.(int64) + value1.(int64), nil
		})

		var p parser.Func
		value := combinator.Memoize(combinator.Any(
			terminal.Integer(),
			&p,
		))

		p = combinator.Memoize(combinator.SeqOf(
			value,
			terminal.Rune('+'),
			value,
		).Bind(add))

		f := text.NewFile("testfile", []byte(input))
		r := text.NewReader(f)
		ctx := parsley.NewContext(parsley.NewFileSet(), r)
		result, err := parsley.Evaluate(ctx, combinator.Sentence(&p), nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(int64(55)))
		Expect(ctx.CallCount()).To(Equal(237769))

	})
})
