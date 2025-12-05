package handler

import (
	"reflect"
	"testing"

	"github.com/sqls-server/sqls/internal/lsp"
)

type codeLensTestCase struct {
	name string
	text string
	want []lsp.CodeLens
}

var codeLensTestCases = []codeLensTestCase{
	{
		name: "code lens with multiple statements",
		text: `SELECT
	*
FROM
	USER;
SELECT
	*
FROM
	USER
WHERE
	id = 1;
`,
		want: []lsp.CodeLens{
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 0},
					End:   lsp.Position{Line: 3, Character: 9},
				},
			},
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 4, Character: 0},
					End:   lsp.Position{Line: 9, Character: 11},
				},
			},
		},
	},
	{
		name: "code lens with empty statement",
		text: `SELECT
	*
FROM
	USER;
`,
		want: []lsp.CodeLens{
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 0},
					End:   lsp.Position{Line: 3, Character: 9},
				},
			},
		},
	}, {
		name: "code lens with one line",
		text: `SELECT * from user; select * from user1;`,
		want: []lsp.CodeLens{
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 0},
					End:   lsp.Position{Line: 0, Character: 19},
				},
			},
			{
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 20},
					End:   lsp.Position{Line: 0, Character: 40},
				},
			},
		},
	},
}

func TestHandleTextDocumentCodeLens(t *testing.T) {
	for _, tt := range codeLensTestCases {
		t.Run(tt.name, func(t *testing.T) {
			codeLens, err := getCodeLens(tt.text, lsp.CodeActionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test.sql",
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(codeLens) != len(tt.want) {
				t.Fatalf("[%s] expected %d code lenses, got %d", tt.name, len(tt.want), len(codeLens))
			}
			for i, lens := range codeLens {
				if !reflect.DeepEqual(lens.Range, tt.want[i].Range) {
					t.Fatalf("[%s] expected code lens %d to be %v, got %v", tt.name, i, tt.want[i].Range, lens.Range)
				}
			}
		})
	}
}
