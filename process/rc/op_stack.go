package rc

import (
	"github.com/c2stack/c2gx/process"
)

type opStack struct {
	code      []process.CodeNode
	head      int
	depth     int
	nextDepth int
}

func newStack() *opStack {
	stack := &opStack{
		head : -1,
		code : make([]process.CodeNode, 128),
	}
	stack.pushOp(&process.SubShell{})
	return stack
}

func (self *opStack) headNode() process.CodeNode {
	if self.head < 0 {
		return nil
	}
	return self.code[0]
}

func (self *opStack) pushOp(op process.Op) {
	// add op to head assuming it will accept, if not go to parent
	head := self.peekNode()
	for head != nil {
		if head.AddOp(op) {
			break
		}
		self.popNode()
		head = self.peekNode()
	}

	// if node, add to stack for next insertion point
	if node, isNode := op.(process.CodeNode); isNode {
		// TODO: Performance, could be fancy and grow array by some heuristic
		self.head++
		self.code[self.head] = node
	}
}

func (self *opStack) peekNode() process.CodeNode {
	if self.head < 0 {
		return nil
	}
	return self.code[self.head]
}

func (self *opStack) popNode() process.CodeNode {
	head := self.code[self.head]
	self.head--
	return head
}
