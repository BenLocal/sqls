package handler

import (
	"testing"

	"github.com/sqls-server/sqls/internal/lsp"
)

func TestHandleTextDocumentCodeLens(t *testing.T) {
	s := `SELECT
    *
FROM
    USER;
SELECT
    *
FROM
    USER
WHERE
    id = 1;
`
	params := lsp.CodeActionParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: "file:///test.sql",
		},
	}
	codeLens, err := getCodeLens(s, params)
	if err != nil {
		t.Fatal(err)
	}

	if len(codeLens) != 2 {
		t.Fatalf("expected 2 code lenses, got %d", len(codeLens))
	}
	f := codeLens[0]
	if f.Range.Start.Line != 0 {
		t.Fatalf("expected start line 0, got %d", f.Range.Start.Line)
	}
	if f.Range.Start.Character != 0 {
		t.Fatalf("expected start character 0, got %d", f.Range.Start.Character)
	}
	if f.Range.End.Line != 3 {
		t.Fatalf("expected end line 3, got %d", f.Range.End.Line)
	}
	if f.Range.End.Character != 9 {
		t.Fatalf("expected end character 9, got %d", f.Range.End.Character)
	}

	second := codeLens[1]
	if second.Range.Start.Line != 4 {
		t.Fatalf("expected start line 4, got %d", second.Range.Start.Line)
	}
	if second.Range.Start.Character != 0 {
		t.Fatalf("expected start character 0, got %d", second.Range.Start.Character)
	}
	if second.Range.End.Line != 9 {
		t.Fatalf("expected end line 9, got %d", second.Range.End.Line)
	}
	if second.Range.End.Character != 11 {
		t.Fatalf("expected end character 11, got %d", second.Range.End.Character)
	}
}
