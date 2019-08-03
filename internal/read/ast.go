package read

// TODO: Move to separate namespace?

type Node interface {
}

type ErrorNode struct {
	Msg string
}

func NewError(msg string) *ErrorNode {
	return &ErrorNode{Msg: msg}
}

type TrueNode struct{}
type FalseNode struct{}
type NilNode struct{}

var TrueObject *TrueNode = &TrueNode{}
var FalseObject *FalseNode = &FalseNode{}
var NilObject *NilNode = &NilNode{}

type StringNode struct {
	Val string
}

func NewString(val string) *StringNode {
	return &StringNode{Val: val}
}

type NumberNode struct {
	Val float64
}

func NewNumber(val float64) *NumberNode {
	return &NumberNode{Val: val}
}

type SymbolNode struct {
	Name string
}

func NewSymbol(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

type ListNode struct {
	Items []Node
}

func NewList(items []Node) *ListNode {
	return &ListNode{Items: items}
}

func NewList2(items ...Node) *ListNode {
	return &ListNode{Items: items}
}

type VectorNode struct {
	Items []Node
}

func NewVector(items []Node) *VectorNode {
	return &VectorNode{Items: items}
}

func NewVector2(items ...Node) *VectorNode {
	return &VectorNode{Items: items}
}

type HashMapNode struct {
	Items map[Node]Node
}

func NewHashMap(items map[Node]Node) *HashMapNode {
	return &HashMapNode{Items: items}
}

func NewHashMap2() *HashMapNode {
	return NewHashMap(make(map[Node]Node))
}
