package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/sqls-server/sqls/internal/lsp"
)

func (s *Server) handleTextDocumentCodeLens(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	if req.Params == nil {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams}
	}

	var params lsp.CodeActionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, err
	}

	f, ok := s.files[params.TextDocument.URI]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", params.TextDocument.URI)
	}
	text := f.Text
	stmts, err := getStatements(text)
	if err != nil {
		return nil, err
	}

	codeLens := []lsp.CodeLens{}
	for _, stmt := range stmts {
		r := lsp.Range{
			Start: lsp.Position{
				Line:      stmt.Pos().Line,
				Character: stmt.Pos().Col,
			},
			End: lsp.Position{
				Line:      stmt.End().Line,
				Character: stmt.End().Col,
			},
		}
		codeLens = append(codeLens, lsp.CodeLens{
			Range: r,
			Command: &lsp.Command{
				Title:   "Execute Query",
				Command: CommandExecuteQuery,
			},
		})
	}

	return codeLens, nil
}
