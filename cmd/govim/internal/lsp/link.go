// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"
	"strconv"

	"github.com/myitcv/govim/cmd/govim/internal/lsp/protocol"
	"github.com/myitcv/govim/cmd/govim/internal/span"
)

func (s *Server) documentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	uri := span.NewURI(params.TextDocument.URI)
	view := s.findView(ctx, uri)
	f, m, err := newColumnMap(ctx, view, uri)
	if err != nil {
		return nil, err
	}
	// find the import block
	ast := f.GetAST(ctx)
	var result []protocol.DocumentLink
	for _, imp := range ast.Imports {
		spn, err := span.NewRange(f.GetFileSet(ctx), imp.Pos(), imp.End()).Span()
		if err != nil {
			return nil, err
		}
		rng, err := m.Range(spn)
		if err != nil {
			return nil, err
		}
		target, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		target = "https://godoc.org/" + target
		result = append(result, protocol.DocumentLink{
			Range:  rng,
			Target: target,
		})
	}
	return result, nil
}