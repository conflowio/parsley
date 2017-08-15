## 1.0.5

IMPROVEMENTS:

* Add cause to parser.Error, add parser.WrapError constructor

## 1.0.4

IMPROVEMENTS:

* All combinators will return with a parser.Func type (not the parser.Parser interface)
* Change internal sepBy to a factory
* Define dependencies with Go Dep
* Generate mocks for testing

## 1.0.3

DEPRECATED:

* combinator.Memoize will be removed in version 2. Use the Memoize method on the history object instead.

IMPROVEMENTS:

* You don't need to use meaningless names for memoization anymore, but you have to be careful to call Memoize
  for your parsers only once.

CHANGES:

* History.GetParserIndex was removed as only the original combinator.Memoize needed it

## 1.0.2

IMPROVEMENTS:

* Add copyright and license headers to all .go files
* Improve code coverage, add notes for (hopefully) impossible panics
* Add codecov.io integration

## 1.0.1

BUG FIXES:

* IsEOF() in text.Reader was not ignoring whitespaces when ignoreWhitespaces was set to true.

IMPROVEMENTS:

* Add Reset() method to History to allow to reuse the parsers for multiple inputs

## 1.0.0

First release
