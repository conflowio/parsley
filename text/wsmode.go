// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text

// WsMode is a type for definining how to handle whitespaces after the tokens
type WsMode uint8

// Whitespace modes
// WsNone means no whitespaces will read and skipped after a token
// WsSpaces means spaces and tabs will be read and skipped automatically after a match
// WsSpacesNl means spaces, tabs and new lines will be read and skipped automatically after a match
const (
	WsNone WsMode = iota
	WsSpaces
	WsSpacesNl
	WsSpacesForceNl
)
