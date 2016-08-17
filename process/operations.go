package process

import (
	"fmt"
	"github.com/c2stack/c2g/c2"
	"bytes"
)

const dumpIndent = "  "

type Op interface {
	Dump(indent string, buf *bytes.Buffer)
	ParentNode() CodeNode
	Exec(stack *Stack) error
}

type CodeNode interface {
	Op

	// If full, then all subsequent code lines should go to parent
	AddOp(op Op) (full bool)
}

// F O R
///////////
type ForOp struct {
	Parent CodeNode
	Iter   string
	Range  []*Value
	code   *CodeBlock
}

func (self *ForOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%sfor %v in %v\n", indent, self.Iter, self.Range))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *ForOp) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}

func (self *ForOp) Exec(stack *Stack) error {
	panic("TODO")
}

func (self *ForOp) ParentNode() CodeNode {
	return self.Parent
}

// W H I L E
///////////
type WhileOp struct {
	Parent    CodeNode
	Condition *Value
	code      *CodeBlock
}

func (self *WhileOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%swhile %s\n", indent, self.Condition))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *WhileOp) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}


func (self *WhileOp) Exec(stack *Stack) error {
	panic("TODO")
}

func (self *WhileOp) ParentNode() CodeNode {
	return self.Parent
}

// I F
///////////
type IfOp struct {
	Parent    CodeNode
	Inverse   bool
	Condition Op
	code      *CodeBlock
}

func (self *IfOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%sif %s\n", indent, self.Condition))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *IfOp) AddOp(op Op) bool {
	if self.Condition == nil {
		self.Condition = op
	} else {
		self.code = AddCode(self.code, op)
	}
	return true
}

func (self *IfOp) Exec(stack *Stack) error {
	panic("TODO")
}
func (self *IfOp) ParentNode() CodeNode {
	return self.Parent
}

// S U B S H E L L
///////////
type SubShell struct {
	Parent CodeNode
	code   *CodeBlock
}

func (self *SubShell) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprint(indent, "@\n"))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *SubShell) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}

func (self *SubShell) Exec(stack *Stack) error {
	if self.code != nil {
		return self.code.Exec(stack)
	}
	return nil
}
func (self *SubShell) ParentNode() CodeNode {
	return self.Parent
}

// S W I T C H
///////////
type SwitchOp struct {
	Parent CodeNode
	code   *CodeBlock
}

func (self *SwitchOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprint(indent, "switch\n"))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *SwitchOp) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}

func (self *SwitchOp) Exec(stack *Stack) error {
	panic("TODO")
}
func (self *SwitchOp) ParentNode() CodeNode {
	return self.Parent
}

// C A S E
///////////
type CaseOp struct {
	Parent     CodeNode
	Conditions []string
	code       *CodeBlock
}

func (self *CaseOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%scase %v\n", indent, self.Conditions))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *CaseOp) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}

func (self *CaseOp) Exec(stack *Stack) error {
	panic("TODO")
}
func (self *CaseOp) ParentNode() CodeNode {
	return self.Parent
}


// F U N C  C A L L
///////////
type FuncCallOp struct {
	Parent CodeNode
	Name   string
	Args   []*Value
}

func (self *FuncCallOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%s%s(%v)\n", indent, self.Name, self.Args))
}

func (self *FuncCallOp) AddOp(op Op) bool {
	return false
}

func (self *FuncCallOp) Exec(stack *Stack) error {
	fdef, found := stack.Func(self.Name)
	if ! found {
		return c2.NewErr("Function not found")
	}
	return fdef.Resolve(stack)
}

func (self *FuncCallOp) ParentNode() CodeNode {
	return self.Parent
}


// F U N C  D E F
///////////
type FuncDefOp struct {
	Parent CodeNode
	Name   string
	code   *CodeBlock
}

func (self *FuncDefOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%sfn %s\n", indent, self.Name))
	self.code.Dump(indent + dumpIndent, buf)
}

func (self *FuncDefOp) AddOp(op Op) bool {
	self.code = AddCode(self.code, op)
	return true
}

func (self *FuncDefOp) Exec(stack *Stack) error {
	stack.SetFunc(self.Name, self)
	return nil
}

func (self *FuncDefOp) Resolve(stack *Stack) error {
	return nil
}

func (self *FuncDefOp) ParentNode() CodeNode {
	return self.Parent
}

// A S S I G N
///////////
type AssignOp struct {
	Parent CodeNode
	Var    string
	Val   *Value
}

func (self *AssignOp) Dump(indent string, buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("%s%s = %v\n", indent, self.Var, self.Val))
}

//func (self *AssignOp) AddOp(op Op) bool {
//	self.code = AddCode(self.code, op)
//	return true
//}
//
func (self *AssignOp) Exec(stack *Stack) error {
c2.Debug.Printf("here")
	obj, err := self.Val.Eval(stack)
	if err != nil {
		return err
	}
	stack.Assign(self.Var, obj)
	return nil
}

func (self *AssignOp) ParentNode() CodeNode {
	return self.Parent
}

// C O D E  B L O C K
type CodeBlock struct {
	code []Op
}

func (self *CodeBlock) Dump(indent string, buf *bytes.Buffer) {
	if self == nil {
		return
	}
	for _, c := range self.code {
		c.Dump(indent, buf)
	}
}

func (self *CodeBlock) AddOp(op Op) bool {
	self.code = append(self.code, op)
	return true
}

func (self *CodeBlock) Exec(stack *Stack) error {
	for _, op := range self.code {
		if err := op.Exec(stack); err != nil {
			return err
		}
	}
	return nil
}

func AddCode(block *CodeBlock, op Op) *CodeBlock {
	if block == nil {
		return &CodeBlock{code: []Op{op}}
	}
	block.code = append(block.code, op)
	return block
}
