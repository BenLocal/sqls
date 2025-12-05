package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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
	return getCodeLens(f.Text, params)
}

func getCodeLens(text string, params lsp.CodeActionParams) ([]lsp.CodeLens, error) {
	stmts, err := getStatements(text)
	if err != nil {
		return nil, err
	}

	codeLens := []lsp.CodeLens{}
	for _, stmt := range stmts {
		if strings.TrimSpace(stmt.String()) == "" {
			continue
		}

		ss := lsp.Position{
			Line:      stmt.Pos().Line,
			Character: stmt.Pos().Col,
		}
		ee := lsp.Position{
			Line:      stmt.End().Line,
			Character: stmt.End().Col,
		}

		tokens := stmt.GetTokens()
		for _, token := range tokens {
			if strings.TrimSpace(token.String()) == "" {
				ss = lsp.Position{
					Line:      token.End().Line,
					Character: token.End().Col,
				}
			} else {
				break
			}
		}

		r := lsp.Range{
			Start: ss,
			End:   ee,
		}
		codeLens = append(codeLens, lsp.CodeLens{
			Range: r,
			Command: &lsp.Command{
				Title:     "Execute Query",
				Command:   CommandExecuteQuery,
				Arguments: []interface{}{params.TextDocument.URI, "", r},
			},
		})
	}

	return codeLens, nil
}
