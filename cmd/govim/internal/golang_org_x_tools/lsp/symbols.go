// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"context"

	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/protocol"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/source"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/telemetry"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/telemetry/log"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/telemetry/trace"
)

func (s *Server) documentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]protocol.DocumentSymbol, error) {
	ctx, done := trace.StartSpan(ctx, "lsp.Server.documentSymbol")
	defer done()

	snapshot, fh, ok, err := s.beginFileRequest(params.TextDocument.URI, source.Go)
	if !ok {
		return []protocol.DocumentSymbol{}, err
	}
	symbols, err := source.DocumentSymbols(ctx, snapshot, fh)
	if err != nil {
		log.Error(ctx, "DocumentSymbols failed", err, telemetry.URI.Of(fh.Identity().URI))
		return []protocol.DocumentSymbol{}, nil
	}
	return symbols, nil
}
