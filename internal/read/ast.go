package read

// TODO: Move to separate namespace?

type Node interface {
}

type ErrorNode struct {
	Msg string
}

func IsError(n Node) bool {
	_, ok := n.(*ErrorNode)
	return ok
}

// TODO: Make it a varargs version with fmt.
func NewError(msg string) *ErrorNode {
	return &ErrorNode{Msg: msg}
}

func IsNil(n Node) bool {
	return n == nil
}

func IsBool(n Node) bool {
	_, ok := n.(bool)
	return ok
}

func IsNumber(n Node) bool {
	_, ok := n.(float64)
	return ok
}

func IsString(n Node) bool {
	_, ok := n.(string)
	return ok
}

type SymbolNode struct {
	Name string
}

func IsSymbol(n Node) bool {
	_, ok := n.(*SymbolNode)
	return ok
}

func NewSymbol(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

type ListNode struct {
	Items []Node
}

func IsList(n Node) bool {
	_, ok := n.(*ListNode)
	return ok
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

func IsVector(n Node) bool {
	_, ok := n.(*VectorNode)
	return ok
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

func IsHashMap(n Node) bool {
	_, ok := n.(*HashMapNode)
	return ok
}

func NewHashMap(items map[Node]Node) *HashMapNode {
	return &HashMapNode{Items: items}
}

// TODO: Rename to NewEmptyHashMap()
func NewHashMap2() *HashMapNode {
	return NewHashMap(make(map[Node]Node))
}
