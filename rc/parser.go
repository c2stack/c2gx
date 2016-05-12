//line syntax.y:2
package rc

import __yyfmt__ "fmt"

//line syntax.y:2
func (l *lexer) Lex(lval *yySymType) int {
	yyDebug = 0 // 4 is most
	t, _ := l.nextToken()
	if t.typ == ParseEof {
		return 0
	}
	lval.ident = t.val
	return int(t.typ)
}

//line syntax.y:16
type yySymType struct {
	yys   int
	ident string
	tree  tree
	word  word
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
const token_word = 57356
const kywd_for = 57357
const kywd_in = 57358
const kywd_while = 57359
const kywd_if = 57360
const kywd_not = 57361
const kywd_twiddle = 57362
const kywd_subshell = 57363
const kywd_switch = 57364
const kywd_fn = 57365

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
	"token_word",
	"kywd_for",
	"kywd_in",
	"kywd_while",
	"kywd_if",
	"kywd_not",
	"kywd_twiddle",
	"kywd_subshell",
	"kywd_switch",
	"kywd_fn",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line syntax.y:135

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 33
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 61

var yyAct = [...]int{

	27, 18, 3, 26, 4, 18, 25, 19, 38, 11,
	19, 15, 36, 7, 8, 9, 10, 12, 13, 16,
	47, 37, 28, 29, 33, 45, 21, 19, 30, 32,
	19, 22, 34, 35, 18, 21, 39, 18, 41, 42,
	22, 18, 44, 39, 43, 41, 24, 32, 18, 46,
	48, 2, 20, 1, 17, 14, 23, 40, 6, 31,
	5,
}
var yyPact = [...]int{

	-4, -1000, 48, 30, -4, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 40, -7, 11, -1000, -1000,
	-1000, -1000, -1000, -1000, -7, 16, -1000, -1000, -7, -7,
	5, -1000, -1000, -4, -1000, -1000, -4, -7, 33, 21,
	-4, -1000, -1000, 13, -1000, -1000, -1000, -4, -1000,
}
var yyPgo = [...]int{

	0, 2, 4, 51, 60, 59, 58, 8, 57, 55,
	0, 3, 6, 54, 53,
}
var yyR1 = [...]int{

	0, 14, 14, 3, 3, 9, 13, 13, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 5, 7,
	7, 6, 6, 11, 10, 4, 4, 12, 12, 2,
	2, 8, 8,
}
var yyR2 = [...]int{

	0, 0, 2, 1, 2, 3, 1, 3, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 3, 1,
	2, 3, 2, 1, 1, 5, 7, 1, 2, 2,
	2, 1, 2,
}
var yyChk = [...]int{

	-1000, -14, -3, -1, -2, -4, -6, 17, 18, 19,
	20, 13, 21, 22, -9, 15, 23, -13, -10, 14,
	4, 5, 10, -3, 6, -12, -11, -10, 11, 12,
	-11, -5, -11, 8, -11, -11, 7, 16, -7, -1,
	-8, -2, -1, -12, 9, 4, -7, 7, -1,
}
var yyDef = [...]int{

	1, -2, 0, 3, 0, 8, 9, 10, 11, 12,
	13, 14, 15, 16, 17, 0, 0, 0, 6, 24,
	2, 29, 30, 4, 0, 22, 27, 23, 0, 0,
	0, 21, 28, 0, 5, 7, 0, 0, 0, 19,
	0, 31, 25, 0, 18, 32, 20, 0, 26,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23,
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

	case 1:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line syntax.y:53
		{
			return 1
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:54
		{
			yylex.(*lexer).tree = yyDollar[1].tree
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:60
		{
			yyVAL.tree = addCodeBlock(yyDollar[1].tree, yyDollar[2].tree)
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:65
		{
			yyVAL.tree = &assignOper{lhs: yyDollar[1].word, rhs: yyDollar[3].word}
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:77
		{
			yyVAL.tree = &whileOper{}
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:78
		{
			yyVAL.tree = &ifOper{}
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:79
		{
			yyVAL.tree = &ifOper{not: true}
		}
	case 13:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:80
		{
			yyVAL.tree = &twiddleOper{}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:81
		{
			yyVAL.tree = &bangOper{}
		}
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:82
		{
			yyVAL.tree = &subshellOper{}
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:83
		{
			yyVAL.tree = &switchOper{}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:87
		{
			yyVAL.tree = yyDollar[2].tree
		}
	case 20:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:93
		{
			yyVAL.tree = addCodeBlock(yyDollar[1].tree, yyDollar[2].tree)
		}
	case 21:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line syntax.y:96
		{
			yyVAL.tree = &fnOper{name: yyDollar[2].word}
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:100
		{
			yyVAL.tree = &fnOper{name: yyDollar[2].word}
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line syntax.y:109
		{
			yyVAL.word = ident(yyDollar[1].ident)
		}
	case 25:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line syntax.y:112
		{
			yyVAL.tree = &forOper{iter: yyDollar[3].word, code: yyDollar[5].tree}
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 26:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line syntax.y:116
		{
			yyVAL.tree = &forOper{iter: yyDollar[3].word, in: yyDollar[5].word, code: yyDollar[7].tree}
			yylex.(*lexer).tree = yyVAL.tree
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:123
		{
			yyVAL.word = addWord(yyDollar[1].word, yyDollar[2].word)
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line syntax.y:129
		{
			yyVAL.tree = parallel(yyDollar[1].tree)
		}
	}
	goto yystack /* stack new state and value */
}
