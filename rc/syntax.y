%{
package rc

func (l *lexer) Lex(lval *yySymType) int {
	yyDebug = 0 // 4 is most
	t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.ident = t.val
    return int(t.typ)
}

%}

%union {
	ident  string
	tree   tree
	word    word
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
%token <tree> kywd_for
%token <tree> kywd_in
%token <tree> kywd_while
%token <tree> kywd_if
%token <tree> kywd_not
%token <tree> kywd_twiddle
%token <tree> kywd_subshell
%token <tree> kywd_switch
%token <tree> kywd_fn

%type <tree> cmd cmdsa line for brace fn body cmdsan assign
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
		$$ = addCodeBlock($1, $2)
		yylex.(*lexer).tree = $$
	}

assign : first token_eq word {
 	$$ = &assignOper{lhs: $1, rhs: $3}
	yylex.(*lexer).tree = $$
}

first :
	comword
	| first token_caret word

cmd :
	for
	| fn
	| kywd_while { $$ = &whileOper{} }
	| kywd_if { $$ = &ifOper{} }
	| kywd_not { $$ = &ifOper{not:true} }
	| kywd_twiddle { $$ = &twiddleOper{} }
	| token_bang { $$ = &bangOper{} }
	| kywd_subshell { $$ = &subshellOper{} }
	| kywd_switch { $$ = &switchOper{} }
	| assign

brace :
	token_open_brace body token_closed_brace {
		$$ = $2
	}

body :
	cmd
	| cmdsan body { $$ = addCodeBlock($1, $2) }

fn :
	kywd_fn words brace {
		$$ = &fnOper{name: $2}
		yylex.(*lexer).tree = $$
	}
	| kywd_fn words {
		$$ = &fnOper{name: $2}
		yylex.(*lexer).tree = $$
	}

word :
	comword

comword :
	token_word { $$ = ident($1) }

for:
	kywd_for token_open_paren word token_closed_paren cmd {
		$$ = &forOper{iter:$3, code: $5}
		yylex.(*lexer).tree = $$
	}
	| kywd_for token_open_paren word kywd_in words token_closed_paren cmd {
		$$ = &forOper{iter: $3, in: $5, code: $7}
		yylex.(*lexer).tree = $$
	}

words :
	word
	| words word {
		$$ = addWord($1, $2)
	}

cmdsa :
	cmd token_semi
	| cmd token_amp { $$ = parallel($1) }

cmdsan :
	cmdsa
	| cmd token_eol

%%