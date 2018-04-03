package parsley

// // Template is a Parsley template
// type Template struct {
// 	name      string
// 	s         *Sentence
// 	parseTree ast.Node
// }
//
// // New creates a new undefined template with a given name and parser
// func New(name string, p parser.Parser) *Template {
// 	return &Template{
// 		name: name,
// 		s:    NewSentence(p),
// 	}
// }
//
// // Must is a helper that wraps a call to a function returning (*Template, error)
// // and panics if the error is non-nil. It is intended for use in variable
// // initializations such as
// //	var t = template.Must(template.New("name", s).Parse("text"))
// func Must(t *Template, err error) *Template {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return t
// }
//
// // Name returns the name of the template.
// func (t *Template) Name() string {
// 	return t.name
// }
//
// func (t *Template) Parse(r reader.Reader) (*Template, error) {
// 	parseTree, _, err := t.s.Parse(r)
// 	if err != nil {
// 		return nil, err
// 	}
// 	t.parseTree = parseTree
// 	return t, nil
// }
