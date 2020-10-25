// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"fmt"

	"github.com/opsidian/parsley/data"
)

// Parse parses the given input and returns with the root node of the AST.
// If a transformer is set on the context then the result will be transformed using it.
// If there are multiple possible parse trees only the first one is returned.
func Parse(ctx *Context, p Parser) (Node, error) {
	var err Error
	var node Node

	if node, _, err = p.Parse(ctx, data.EmptyIntMap, ctx.Reader().Pos(0)); err != nil {
		if ctxErr := ctx.Error(); ctxErr != nil && ctxErr.Pos() > err.Pos() {
			err = ctxErr
		}

		return nil, fmt.Errorf("failed to parse the input: %w", ctx.FileSet().ErrorWithPosition(err))
	}

	if ctx.TransformationEnabled() {
		node, err = Transform(ctx.UserContext(), node)
		if err != nil {
			return nil, ctx.FileSet().ErrorWithPosition(err)
		}
	}

	if ctx.StaticCheckEnabled() {
		if err = StaticCheck(ctx.UserContext(), node); err != nil {
			return nil, ctx.FileSet().ErrorWithPosition(err)
		}
	}

	return node, nil
}
