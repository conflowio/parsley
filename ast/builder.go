// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import "github.com/opsidian/parsley/parsley"

// NodeBuilderFunc defines a helper to implement the NodeBuilder interface with functions
type NodeBuilderFunc func([]parsley.Node) parsley.Node

// BuildNode combines the given nodes into a single node
func (f NodeBuilderFunc) BuildNode(nodes []parsley.Node) parsley.Node {
	return f(nodes)
}
