// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := parser.NewError(test.NewPosition(1), "ERR1")
	assert.Equal(t, test.NewPosition(1), err.Pos())
	assert.Equal(t, "ERR1", err.Msg())
	assert.Equal(t, "ERR1 at Pos{1}", err.Error())
}
