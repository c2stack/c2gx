%{
package rc

import "github.com/c2gx/process"

// Parser for Rc
//  https://en.wikipedia.org/wiki/Rc

func (l *lexer) Lex(lval *yySymType) int {
	t, _ := l.nextToken()
    if t.typ == ParseEof {
        return 0
    }
    lval.token = t.val
    lval.stack = l.stack
    return int(t.typ)
}

func tokenString(s string) string {
    return s[1:len(s) - 1]
}

%}

%union {
	token   string
	value   *process.Value
	values  []*process.Value
	fcall 	*process.FuncCallOp
	stack   *opStack
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
%token <token> token_ident
%token <token> token_string
%token <token> token_number
%token token_for
%token token_in
%token token_while
%token token_if
%token token_not
%token token_twiddle
%token token_subshell
%token token_switch
%token token_fn
%token token_backtick
%type <value> value
%type <values> values opt_values
%type <fcall> fn_call

%%

/*
  Loosely based on
    http://plan9.bell-labs.com/sources/plan9/sys/src/cmd/rc/syn.y
  but should have no known incompatibilities with Rc syntax
*/

rc:
	/* empty */
	| line token_eol

line :
	cmd
	| inline_cmds line

body :
	cmd
	| cmds body

inline_cmds :
	cmd token_semi
	| cmd token_amp {
	 	// TODO: Fork
	}

cmds :
	inline_cmds
	| cmd token_eol

assign :
	token_ident token_eq value {
		yyVAL.stack.pushOp(&process.AssignOp{Parent:yyVAL.stack.peekNode(), Var: $1, Val: $3})
	}

opt_values :
	/* empty */ { $$ = []*process.Value{} }
	| values

values :
	value { $$ = []*process.Value{$1} }
	| values value { $$ = append($1, $2) }
	/* | function call */

value:
	token_ident { $$ = &process.Value{Var:$1} }
	| token_string { $$ = &process.Value{Str:tokenString($1)} }
	| token_number { $$ = &process.Value{Num:$1} }
	| token_backtick fn_call {
		$$ = &process.Value{Func:$2}
	}

cmd :
	for cmd
	| fn_def
	| fn_call
	| token_while {
		yyVAL.stack.pushOp(&process.WhileOp{Parent:yyVAL.stack.peekNode()})
	}
	| token_if {
		yyVAL.stack.pushOp(&process.IfOp{Parent:yyVAL.stack.peekNode()})
	}
	| token_not {
		yyVAL.stack.pushOp(&process.IfOp{Parent:yyVAL.stack.peekNode(), Inverse:true})
	}
	| token_bang {
		panic("TODO")
	}
	| token_subshell {
		yyVAL.stack.pushOp(&process.SubShell{Parent:yyVAL.stack.peekNode()})
	}
	| token_switch {
		yyVAL.stack.pushOp(&process.SwitchOp{Parent:yyVAL.stack.peekNode()})
	}
	| assign
	| token_open_brace cmds token_closed_brace {

	}
	| token_eol

fn_call :
 	token_ident opt_values token_eol {
		$$ = &process.FuncCallOp{Parent:yyVAL.stack.peekNode(), Name: $1, Args: $2}
		yyVAL.stack.pushOp($$)
	}

fn_def_def :
	token_fn token_ident token_open_brace token_eol {
		yyVAL.stack.pushOp(&process.FuncDefOp{Parent:yyVAL.stack.peekNode(), Name: $2})
	}

brace :
	token_open_brace body token_closed_brace {
		yyVAL.stack.pushOp(&process.SubShell{Parent:yyVAL.stack.peekNode()})
	}

fn_def :
	fn_def_def brace {
		yyVAL.stack.popNode()
	}

for:
	token_for token_open_paren token_ident token_closed_paren {
		// lex.skipnl()
		yyVAL.stack.pushOp(&process.ForOp{Parent:yyVAL.stack.peekNode(), Iter:$3})
	}
	| token_for token_open_paren token_ident token_in values token_closed_paren {
		// lex.skipnl()
		yyVAL.stack.pushOp(&process.ForOp{Parent:yyVAL.stack.peekNode(), Iter: $3, Range: $5})
	}

%%