package rc

import (
	"fmt"
)

type tree interface {
	fmt.Stringer
}

type forOper struct{
	iter word
	in word
	code tree
}

func (self *forOper) String() string {
	return fmt.Sprintf("for (%v in %v) %v", self.iter, self.in, self.code)
}

type whileOper struct{
	condition word
	code tree
}

func (self *whileOper) String() string {
	return "while"
}

type ifOper struct{
	not bool
	code tree
}
func (self *ifOper) String() string {
	return "if"
}

type twiddleOper struct{
}
func (self *twiddleOper) String() string {
	return "~"
}

type bangOper struct{
}
func (self *bangOper) String() string {
	return "!"
}

type subshellOper struct{
}
func (self *subshellOper) String() string {
	return "@"
}

type switchOper struct{
	cases []tree
}
func (self *switchOper) String() string {
	return "switch"
}

type fnOper struct{
	name word
}
func (self *fnOper) String() string {
	return fmt.Sprintf("fn %v", self.name)
}

type assignOper struct {
	lhs word
	rhs word
}

func (self *assignOper) String() string {
	return fmt.Sprintf("%v = %v", self.lhs, self.rhs)
}

type forkOper struct {
	a tree
	b tree
}

func (self *forkOper) String() string {
	return "fork"
}

type codeBlock struct {
	lines []tree
}

func (self *codeBlock) String() string {
	return fmt.Sprintf("code %v", self.lines)
}

func parallel(code tree) (*forkOper) {
	return &forkOper{
		a : code,
	}
}

func addCodeBlock(a tree, b tree) (*codeBlock) {
	if code, isCode := a.(*codeBlock); isCode {
		code.lines = append(code.lines, b)
		return code
	} else if fork, isFork := a.(*forkOper); isFork {
		fork.b = b
		return &codeBlock{
			lines : []tree { fork },
		}
	}
	return &codeBlock{
		lines : []tree { a, b },
	}
}

type word interface {
	fmt.Stringer
}

type ident string

func (self ident) String() string {
	return string(self)
}

type multiIdent struct {
	idents []word
}

func (self *multiIdent) String() string {
	return fmt.Sprintf("%v", self.idents)
}

func addWord(a word, b word) (*multiIdent) {
	if multi, isMulti := a.(*multiIdent); isMulti {
		multi.idents = append(multi.idents, b)
		return multi
	}
	return &multiIdent{
		idents : []word{b},
	}
}
