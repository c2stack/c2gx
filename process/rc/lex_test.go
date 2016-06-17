package rc

import (
	"testing"
	"fmt"
)

func TestRcLex(t *testing.T) {
	tests := []struct{
		label string
		script string
		expected []int
	} {
		{
			"empty",
			"",
			[]int{},
		},
		{
			"one word",
			"switch",
			[]int{
				token_switch,
			},
		},
		{
			"quote",
			"'hello world'",
			[]int{
				token_string,
			},
		},
		{
			"whitespace",
			"  switch  ",
			[]int{
				token_switch,
			},
		},
		{
			"multi-line",
			"switch\nswitch",
			[]int{
				token_switch, token_eol, token_switch,
			},
		},
		{
			"commands",
			"for ( x )",
			[]int{
				token_for, token_open_paren, token_ident, token_closed_paren,
			},
		},
		{
			"fn",
			"fn x { }",
			[]int{
				token_fn, token_ident, token_open_brace, token_closed_brace,
			},
		},
	}
	for _, test := range tests {
		label := fmt.Sprintf("%s - '%s'", test.label, test.script)
		l := lex(test.script)
		for j, tok := range test.expected {
			actual, err := l.nextToken()
			if err != nil {
				t.Error(err)
			}
			if actual.typ != tok {
				t.Errorf("%s, token #%d - expected '%s' but got '%s'(%s)", label, j,
					l.keyword(tok), l.keyword(actual.typ), actual.val)
			}
		}
		actualEof, err := l.nextToken()
		if err != nil {
			t.Error(err)
		}
		if actualEof.typ != ParseEof {
			t.Errorf("%s\nExpected:eof\n  Actual: %s-%s", label, l.keyword(actualEof.typ), actualEof.val)
		}
	}
}
