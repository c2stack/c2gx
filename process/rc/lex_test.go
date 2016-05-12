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
				kywd_switch,
			},
		},
		{
			"whitespace",
			"  switch  ",
			[]int{
				kywd_switch,
			},
		},
		{
			"multi-line",
			"switch\nswitch",
			[]int{
				kywd_switch, token_eol, kywd_switch,
			},
		},
		{
			"commands",
			"for ( x )",
			[]int{
				kywd_for, token_open_paren, token_word, token_closed_paren,
			},
		},
	}
	only := ""//commands"
	for _, test := range tests {
		if len(only) > 0 && test.label != only {
			continue
		}
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
