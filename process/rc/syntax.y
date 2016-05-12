%{
package rc

import "github.com/c2gx/process"

func (l *lexer) Lex(lval *yySymType) int {
	yyDebug = 0 // 0 is off, 4 is most
	t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.ident = t.val
    return int(t.typ)
}

%}

%union {
	ident   string
	op 		process.Op
	word	process.Ident
}

%token token_eol
%token token_semi
%token token_open_paren
%token token_closed_paren
%token token_open_brace
%token token_closed_brace
%token token_amp
%token token_eq
%token token_caret
%token token_bang
%token <ident> token_word
%token <op> kywd_for
%token <op> kywd_in
%token <op> kywd_while
%token <op> kywd_if
%token <op> kywd_not
%token <op> kywd_twiddle
%token <op> kywd_subshell
%token <op> kywd_switch
%token <op> kywd_fn

%type <op> cmd cmdsa line for brace fn body cmdsan assign
%type <word> comword word words first

%%

/*
  Based on
    http://plan9.bell-labs.com/sources/plan9/sys/src/cmd/rc/syn.y
*/
rc :
	/* empty */ { return 1 }
	| line token_eol {
		yylex.(*lexer).tree = $1
	};

line :
	cmd
	| cmdsa line {
		$$ = process.AddCode($1, $2)
		yylex.(*lexer).tree = $$
	}

assign : first token_eq word {
 	$$ = &process.AssignOp{Lhs: $1, Rhs: $3}
	yylex.(*lexer).tree = $$
}

first :
	comword
	| first token_caret word

cmd :
	for
	| fn
	| kywd_while { $$ = &process.WhileOp{} }
	| kywd_if { $$ = &process.IfOp{} }
	| kywd_not { $$ = &process.IfOp{Inverse:true} }
	| kywd_twiddle { $$ = &process.TwiddleOp{} }
	| token_bang { $$ = &process.BangOp{} }
	| kywd_subshell { $$ = &process.SubShellOp{} }
	| kywd_switch { $$ = &process.SwitchOp{} }
	| assign

brace :
	token_open_brace body token_closed_brace {
		$$ = $2
	}

body :
	cmd
	| cmdsan body { $$ = process.AddCode($1, $2) }

fn :
	kywd_fn words brace {
		$$ = &process.FuncOp{Name: $2}
		yylex.(*lexer).tree = $$
	}
	| kywd_fn words {
		$$ = &process.FuncOp{Name: $2}
		yylex.(*lexer).tree = $$
	}

word :
	comword

comword :
	token_word { $$ = process.Word($1) }

for:
	kywd_for token_open_paren word token_closed_paren cmd {
		$$ = &process.ForOp{Iter:$3, Code: $5}
		yylex.(*lexer).tree = $$
	}
	| kywd_for token_open_paren word kywd_in words token_closed_paren cmd {
		$$ = &process.ForOp{Iter: $3, In: $5, Code: $7}
		yylex.(*lexer).tree = $$
	}

words :
	word
	| words word {
		$$ = process.AddIdent($1, $2)
	}

cmdsa :
	cmd token_semi
	| cmd token_amp { $$ = process.ForkCode($1) }

cmdsan :
	cmdsa
	| cmd token_eol

%%