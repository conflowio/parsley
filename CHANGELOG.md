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
