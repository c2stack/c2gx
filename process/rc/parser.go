//line rc.y:2
package rc

import __yyfmt__ "fmt"

//line rc.y:2
import "github.com/dhubler/c2gx/process"

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
	return s[1 : len(s)-1]
}

//line rc.y:25
type yySymType struct {
	yys    int
	token  string
	value  *process.Value
	values []*process.Value
	fcall  *process.FuncCallOp
	stack  *opStack
}

const token_eol = 57346
const token_semi = 57347
const token_open_paren = 57348
const token_closed_paren = 57349
const token_open_brace = 57350
const token_closed_brace = 57351
const token_amp = 57352
const token_eq = 57353
const token_caret = 57354
const token_bang = 57355
const token_ident = 57356
const token_string = 57357
const token_number = 57358
const token_for = 57359
const token_in = 57360
const token_while = 57361
const token_if = 57362
const token_not = 57363
const token_twiddle = 57364
const token_subshell = 57365
const token_switch = 57366
const token_fn = 57367
const token_backtick = 57368

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"token_eol",
	"token_semi",
	"token_open_paren",
	"token_closed_paren",
	"token_open_brace",
	"token_closed_brace",
	"token_amp",
	"token_eq",
	"token_caret",
	"token_bang",
	"token_ident",
	"token_string",
	"token_number",
	"token_for",
	"token_in",
	"token_while",
	"token_if",
	"token_not",
	"token_twiddle",
	"token_subshell",
	"token_switch",
	"token_fn",
	"token_backtick",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line rc.y:171

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 38
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 72

var yyAct = [...]int{

	35, 34, 3, 44, 7, 59, 16, 51, 25, 43,
	15, 40, 36, 37, 38, 11, 19, 53, 28, 17,
	55, 8, 9, 10, 39, 12, 13, 20, 54, 46,
	36, 37, 38, 47, 45, 49, 32, 42, 22, 36,
	37, 38, 39, 23, 50, 26, 41, 22, 52, 45,
	56, 39, 23, 31, 29, 2, 58, 27, 4, 49,
	24, 57, 4, 48, 21, 30, 18, 6, 5, 14,
	1, 33,
}
var yyPact = [...]int{

	2, -1000, 60, 42, 2, 2, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 2, -1000, 48, 45, 25,
	-3, -1000, -1000, -1000, -1000, -1000, 37, -1000, 33, -5,
	-1000, 2, 16, 59, 16, -1000, -1000, -1000, -1000, -7,
	40, -1000, -1000, 10, 11, 33, 2, -1000, -1000, -1000,
	-1000, 16, 57, -1000, 16, -1000, -1000, -1000, -2, -1000,
}
var yyPgo = [...]int{

	0, 0, 1, 71, 4, 70, 55, 2, 57, 3,
	29, 69, 68, 67, 66, 65,
}
var yyR1 = [...]int{

	0, 5, 5, 6, 6, 9, 9, 8, 8, 10,
	10, 11, 3, 3, 2, 2, 1, 1, 1, 1,
	7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
	7, 7, 4, 14, 15, 13, 12, 12,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 2, 1, 2, 2, 2, 1,
	2, 3, 0, 1, 1, 2, 1, 1, 1, 2,
	2, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	3, 1, 3, 4, 3, 2, 4, 6,
}
var yyChk = [...]int{

	-1000, -5, -6, -7, -8, -12, -13, -4, 19, 20,
	21, 13, 23, 24, -11, 8, 4, 17, -14, 14,
	25, 4, 5, 10, -6, -7, -10, -8, -7, 6,
	-15, 8, 11, -3, -2, -1, 14, 15, 16, 26,
	14, 9, 4, 14, -9, -7, -10, -1, 4, -1,
	-4, 14, 8, 7, 18, 9, -9, 4, -2, 7,
}
var yyDef = [...]int{

	1, -2, 0, 3, 0, 0, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 0, 31, 0, 0, 12,
	0, 2, 7, 8, 4, 20, 0, 9, 0, 0,
	35, 0, 0, 0, 13, 14, 16, 17, 18, 0,
	0, 30, 10, 0, 0, 5, 0, 11, 32, 15,
	19, 12, 0, 36, 0, 34, 6, 33, 0, 37,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 8:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rc.y:82
		{
			// TODO: Fork
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rc.y:91
		{
			yyVAL.stack.pushOp(&process.AssignOp{Parent: yyVAL.stack.peekNode(), Var: yyDollar[1].token, Val: yyDollar[3].value})
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line rc.y:96
		{
			yyVAL.values = []*process.Value{}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:100
		{
			yyVAL.values = []*process.Value{yyDollar[1].value}
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rc.y:101
		{
			yyVAL.values = append(yyDollar[1].values, yyDollar[2].value)
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:105
		{
			yyVAL.value = &process.Value{Var: yyDollar[1].token}
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:106
		{
			yyVAL.value = &process.Value{Str: tokenString(yyDollar[1].token)}
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:107
		{
			yyVAL.value = &process.Value{Num: yyDollar[1].token}
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rc.y:108
		{
			yyVAL.value = &process.Value{Func: yyDollar[2].fcall}
		}
	case 23:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:116
		{
			yyVAL.stack.pushOp(&process.WhileOp{Parent: yyVAL.stack.peekNode()})
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:119
		{
			yyVAL.stack.pushOp(&process.IfOp{Parent: yyVAL.stack.peekNode()})
		}
	case 25:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:122
		{
			yyVAL.stack.pushOp(&process.IfOp{Parent: yyVAL.stack.peekNode(), Inverse: true})
		}
	case 26:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:125
		{
			panic("TODO")
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:128
		{
			yyVAL.stack.pushOp(&process.SubShell{Parent: yyVAL.stack.peekNode()})
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line rc.y:131
		{
			yyVAL.stack.pushOp(&process.SwitchOp{Parent: yyVAL.stack.peekNode()})
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rc.y:135
		{

		}
	case 32:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rc.y:141
		{
			yyVAL.fcall = &process.FuncCallOp{Parent: yyVAL.stack.peekNode(), Name: yyDollar[1].token, Args: yyDollar[2].values}
			yyVAL.stack.pushOp(yyVAL.fcall)
		}
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rc.y:147
		{
			yyVAL.stack.pushOp(&process.FuncDefOp{Parent: yyVAL.stack.peekNode(), Name: yyDollar[2].token})
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line rc.y:152
		{
			yyVAL.stack.pushOp(&process.SubShell{Parent: yyVAL.stack.peekNode()})
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line rc.y:157
		{
			yyVAL.stack.popNode()
		}
	case 36:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line rc.y:162
		{
			// lex.skipnl()
			yyVAL.stack.pushOp(&process.ForOp{Parent: yyVAL.stack.peekNode(), Iter: yyDollar[3].token})
		}
	case 37:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line rc.y:166
		{
			// lex.skipnl()
			yyVAL.stack.pushOp(&process.ForOp{Parent: yyVAL.stack.peekNode(), Iter: yyDollar[3].token, Range: yyDollar[5].values})
		}
	}
	goto yystack /* stack new state and value */
}
