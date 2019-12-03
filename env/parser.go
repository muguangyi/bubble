// Copyright 2019 Bubble. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package env

type NodeType int

const (
	NODE_EMPTY NodeType = iota
	NODE_VALUE
	NODE_VARIABLE
	NODE_METHOD
	NODE_PARAM
)

type INode interface {
	Type() NodeType
	Execute(env IEnv) string
	Prev() INode
	Next() INode
	Push(n INode) INode
	Pop() INode
}

func NewParser() IParser {
	p := &parser{root: newEmpty()}
	p.cur = p.root

	return p
}

type IParser interface {
	Interpret(t TokenType, v string)
	Result() INode
}

type parser struct {
	root INode
	cur  INode
}

func (p *parser) Interpret(t TokenType, v string) {
	switch t {
	case TOKEN_VALUE:
		p.cur = p.cur.Push(newValue(v))
	case TOKEN_VARIABLE:
		p.cur = p.cur.Push(newVariable(v))
	case TOKEN_BEGIN_METHOD:
		p.cur = p.cur.Push(newMethod(v))
	case TOKEN_END_PARAM:
		p.cur = p.cur.Pop()
	case TOKEN_END_METHOD:
		for p.cur.Type() != NODE_METHOD {
			p.cur = p.cur.Pop()
		}
		p.cur = p.cur.Pop()
	}
}

func (p *parser) Result() INode {
	return p.root
}

// --- Node ---

type node struct {
	t    NodeType
	prev INode
	next INode
}

func (n *node) Type() NodeType {
	return n.t
}

func (n *node) Prev() INode {
	return n.prev
}

func (n *node) Next() INode {
	return n.next
}

func AsNext(n INode, next INode) {
	switch n.Type() {
	case NODE_EMPTY:
		n.(*empty).next = next
	case NODE_VALUE:
		n.(*value).next = next
	case NODE_VARIABLE:
		n.(*variable).next = next
	case NODE_METHOD:
		n.(*method).next = next
	case NODE_PARAM:
		n.(*param).next = next
	}
}

func AsPrev(n INode, prev INode) {
	switch n.Type() {
	case NODE_EMPTY:
		n.(*empty).prev = prev
	case NODE_VALUE:
		n.(*value).prev = prev
	case NODE_VARIABLE:
		n.(*variable).prev = prev
	case NODE_METHOD:
		n.(*method).prev = prev
	case NODE_PARAM:
		n.(*param).prev = prev
	}
}

// --- Empty ---

func newEmpty() INode {
	n := &empty{}
	n.t = NODE_EMPTY

	return n
}

type empty struct {
	node
}

func (n *empty) Execute(env IEnv) string {
	return ""
}

func (n *empty) Push(x INode) INode {
	AsNext(n, x)
	AsPrev(x, n)

	return x
}

func (n *empty) Pop() INode {
	return n
}

// --- Value ---

func newValue(v string) INode {
	n := &value{value: v}
	n.t = NODE_VALUE

	return n
}

type value struct {
	node
	value string
}

func (v *value) Execute(env IEnv) string {
	return v.value
}

func (v *value) Push(x INode) INode {
	AsNext(v, x)
	AsPrev(x, v)

	return x
}

func (v *value) Pop() INode {
	return v
}

// --- Variable ---

func newVariable(name string) INode {
	n := &variable{name: name}
	n.t = NODE_VARIABLE

	return n
}

type variable struct {
	node
	name string
}

func (v *variable) Execute(env IEnv) string {
	ret, err := env.Get(v.name)
	if err != nil {
		return v.name
	}

	return ret.ToString()
}

func (v *variable) Push(x INode) INode {
	AsNext(v, x)
	AsPrev(x, v)

	return x
}

func (v *variable) Pop() INode {
	return v
}

// --- Function ---

func newMethod(name string) INode {
	n := &method{name: name, params: make([]INode, 0), end: false}
	n.t = NODE_METHOD

	return n
}

type method struct {
	empty
	name   string
	params []INode
	end    bool
}

func (m *method) Execute(env IEnv) string {
	args := make([]IAny, len(m.params))
	for i, p := range m.params {
		args[i] = NewAny(p.Execute(env))
	}

	fn, err := env.GetFunc(m.name)
	if err != nil {
		return ""
	}

	ret, err := fn(args...)
	if err != nil {
		return ""
	}

	return ret.ToString()
}

func (m *method) Push(x INode) INode {
	if m.end {
		AsNext(m, x)
		AsPrev(x, m)
		return x
	} else {
		p := newParam()
		m.params = append(m.params, p)
		AsPrev(p, m)

		return p.Push(x)
	}
}

func (m *method) Pop() INode {
	m.end = true

	if m.prev.Type() == NODE_METHOD || m.prev.Type() == NODE_PARAM {
		return m.prev
	}

	return m
}

// --- Param ---

func newParam() INode {
	n := &param{subs: make([]INode, 0)}
	n.t = NODE_PARAM

	return n
}

type param struct {
	node
	subs []INode
}

func (p *param) Execute(env IEnv) string {
	compose := ""
	for _, s := range p.subs {
		compose += s.Execute(env)
	}

	return compose
}

func (p *param) Push(x INode) INode {
	AsPrev(x, p)

	p.subs = append(p.subs, x)

	if x.Type() == NODE_METHOD {
		return x
	} else {
		return p
	}
}

func (p *param) Pop() INode {
	return p.prev
}
