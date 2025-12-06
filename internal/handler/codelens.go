package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
	"github.com/sqls-server/sqls/ast/astutil"
	"github.com/sqls-server/sqls/internal/lsp"
	"github.com/sqls-server/sqls/token"
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
	emptySpaceMatcher := astutil.NodeMatcher{
		ExpectTokens: []token.Kind{
			token.Whitespace,
			token.Comment,
			token.MultilineComment,
		},
	}
	for _, stmt := range stmts {
		ss := lsp.Position{
			Line:      stmt.Pos().Line,
			Character: stmt.Pos().Col,
		}
		ee := lsp.Position{
			Line:      stmt.End().Line,
			Character: stmt.End().Col,
		}

		tokens := stmt.GetTokens()
		hasValidToken := false
		for _, token := range tokens {
			if emptySpaceMatcher.IsMatch(token) {
				ss = lsp.Position{
					Line:      token.End().Line,
					Character: token.End().Col,
				}
			} else {
				hasValidToken = true
				break
			}
		}

		if !hasValidToken {
			continue
		}

		r := lsp.Range{
			Start: ss,
			End:   ee,
		}
		codeLens = append(codeLens, lsp.CodeLens{
			Range: r,
			Command: &lsp.Command{
				Title:     "Execute",
				Command:   CommandExecuteQuery,
				Arguments: []interface{}{params.TextDocument.URI, "", r},
			},
		})
	}

	return codeLens, nil
}
