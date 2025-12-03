package handler

import (
	"context"
	"encoding/json"

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

	codeLens := []lsp.CodeLens{
		{
			Range: params.Range,
			Command: &lsp.Command{
				Title:   "Execute Query",
				Command: CommandExecuteQuery,
			},
		},
	}
	return codeLens, nil
}
