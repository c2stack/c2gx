package process

import (
	"fmt"
)

type Op interface {
	fmt.Stringer
}

type ForOp struct{
	Iter Ident
	In   Ident
	Code Op
}

func (self *ForOp) String() string {
	return fmt.Sprintf("for (%v in %v) %v", self.Iter, self.In, self.Code)
}

type WhileOp struct{
	Condition Ident
	Code      Op
}

func (self *WhileOp) String() string {
	return "while"
}

type IfOp struct{
	Inverse bool
	Code    Op
}
func (self *IfOp) String() string {
	return "if"
}

type TwiddleOp struct{
}
func (self *TwiddleOp) String() string {
	return "~"
}

type BangOp struct{
}
func (self *BangOp) String() string {
	return "!"
}

type SubShellOp struct{
}
func (self *SubShellOp) String() string {
	return "@"
}

type SwitchOp struct{
	Cases []Op
}
func (self *SwitchOp) String() string {
	return "switch"
}

type FuncOp struct{
	Name Ident
}
func (self *FuncOp) String() string {
	return fmt.Sprintf("fn %v", self.Name)
}

type AssignOp struct {
	Lhs Ident
	Rhs Ident
}

func (self *AssignOp) String() string {
	return fmt.Sprintf("%v = %v", self.Lhs, self.Rhs)
}

type ForkOp struct {
	A Op
	B Op
}

func (self *ForkOp) String() string {
	return "fork"
}

type CodeBlock struct {
	Lines []Op
}

func (self *CodeBlock) String() string {
	return fmt.Sprintf("code %v", self.Lines)
}

func ForkCode(code Op) (*ForkOp) {
	return &ForkOp{
		A : code,
	}
}

func AddCode(a Op, b Op) Op {
	if code, isCode := a.(*CodeBlock); isCode {
		code.Lines = append(code.Lines, b)
		return code
	} else if fork, isFork := a.(*ForkOp); isFork {
		fork.B = b
		return &CodeBlock{
			Lines : []Op{ fork },
		}
	}
	return &CodeBlock{
		Lines : []Op{ a, b },
	}
}

type Ident interface {
	fmt.Stringer
}

type Word string

func (self Word) String() string {
	return string(self)
}

type IdentExpression struct {
	Idents []Ident
}

func (self *IdentExpression) String() string {
	return fmt.Sprintf("%v", self.Idents)
}

func AddIdent(a Ident, b Ident) (*IdentExpression) {
	if multi, isMulti := a.(*IdentExpression); isMulti {
		multi.Idents = append(multi.Idents, b)
		return multi
	}
	return &IdentExpression{
		Idents : []Ident{b},
	}
}
