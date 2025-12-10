package handler

import (
	"testing"

	"github.com/sqls-server/sqls/internal/config"
	"github.com/sqls-server/sqls/internal/database"
	"github.com/sqls-server/sqls/internal/lsp"
)

func Test_executeQuery(t *testing.T) {
	tx := newTestContext()
	tx.setup(t)
	defer tx.tearDown()

	didChangeConfigurationParams := lsp.DidChangeConfigurationParams{
		Settings: struct {
			SQLS *config.Config "json:\"sqls\""
		}{
			SQLS: &config.Config{
				Connections: []*database.DBConfig{
					{
						Driver:         "mock",
						DataSourceName: "",
					},
				},
			},
		},
	}
	if err := tx.conn.Call(tx.ctx, "workspace/didChangeConfiguration", didChangeConfigurationParams, nil); err != nil {
		t.Fatal("conn.Call workspace/didChangeConfiguration:", err)
	}

	uri := "file:///test.sql"
	text := "SELECT 1; SELECT 2;"
	didOpenParams := lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "sql",
			Version:    0,
			Text:       text,
		},
	}
	if err := tx.conn.Call(tx.ctx, "textDocument/didOpen", didOpenParams, nil); err != nil {
		t.Fatal("conn.Call textDocument/didOpen:", err)
	}
	tx.testFile(t, didOpenParams.TextDocument.URI, didOpenParams.TextDocument.Text)

	// executeCommandParams := lsp.ExecuteCommandParams{
	// 	Command:   CommandExecuteQuery,
	// 	Arguments: []interface{}{uri},
	// }
	// var got interface{}
	// tx.conn.Call(tx.ctx, "workspace/executeCommand", executeCommandParams, &got)
	// pass error
}

func Test_extractRangeText(t *testing.T) {
	type args struct {
		text      string
		startLine int
		startChar int
		endLine   int
		endChar   int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "extract single line",
			args: args{
				text:      "select * from city",
				startLine: 0,
				startChar: 0,
				endLine:   0,
				endChar:   8,
			},
			want: "select *",
		},
		{
			name: "extract multi line with not equal start end",
			args: args{
				text:      "select 1;\nselect 2;\nselect 3;",
				startLine: 0,
				startChar: 7,
				endLine:   2,
				endChar:   8,
			},
			want: "1;\nselect 2;\nselect 3",
		},
		{
			name: "extract multi line with equal start end",
			args: args{
				text:      "select 1;\nselect 2;\nselect 3;",
				startLine: 1,
				startChar: 2,
				endLine:   1,
				endChar:   6,
			},
			want: "lect",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractRangeText(tt.args.text, tt.args.startLine, tt.args.startChar, tt.args.endLine, tt.args.endChar); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_getStatements(t *testing.T) {

	tests := []struct {
		name string
		text string
		want []string
	}{
		{
			name: "single statement",
			text: "SELECT * FROM `user` a;",
			want: []string{
				"SELECT * FROM `user` a;",
			},
		},
		{
			name: "multiple statements",
			text: "SELECT * FROM `user` a;\nSELECT * FROM `user` b;",
			want: []string{
				"SELECT * FROM `user` a;",
				"\nSELECT * FROM `user` b;",
			},
		},
		{
			name: "multiple statements with comments",
			text: "SELECT * FROM `user` a;\n-- this is a comment\nSELECT * FROM `user` b;",
			want: []string{
				"SELECT * FROM `user` a;",
				"\n-- this is a comment\nSELECT * FROM `user` b;",
			},
		},
		{
			name: "multiple statements with comments and empty lines",
			text: "SELECT * FROM `user` a;\n-- this is a comment\n\nSELECT * FROM `user` b;",
			want: []string{
				"SELECT * FROM `user` a;",
				"\n-- this is a comment\n\nSELECT * FROM `user` b;",
			},
		},
		{
			name: "multiple statements with first statement comment",
			text: "-- this is a test\nSELECT * FROM `user` a;\n\n-- this is a test\n-- this is a test\ninsert into `user` (id, name, age) values (6, 'test', 10);",
			want: []string{
				"-- this is a test\nSELECT * FROM `user` a;",
				"\n\n-- this is a test\n-- this is a test\ninsert into `user` (id, name, age) values (6, 'test', 10);",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := getStatements(tt.text)
			if err != nil {
				t.Errorf("getStatements() error = %v", err)
			}
			if len(stmts) != len(tt.want) {
				t.Errorf("getStatements() = %v, want %v", len(stmts), len(tt.want))
			}
			for i, s := range stmts {
				if s.String() != tt.want[i] {
					t.Errorf("getStatements() = %v, want %v", s.String(), tt.want[i])
				}
			}
		})
	}
}
