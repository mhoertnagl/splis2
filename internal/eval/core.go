package eval

import (
	"reflect"

	"github.com/mhoertnagl/splis2/internal/data"
)

func InitCore(e Evaluator) {
	e.AddCoreFun("list", list)
	e.AddCoreFun("list?", isList)
	e.AddCoreFun("count", count)
	e.AddCoreFun("empty?", isEmpty)
	e.AddCoreFun("+", evalxf("+", sum))
	e.AddCoreFun("-", eval12f("-", neg, diff))
	e.AddCoreFun("*", evalxf("*", prod))
	e.AddCoreFun("/", eval12f("/", reciproc, div))
	e.AddCoreFun("=", eq)
	e.AddCoreFun("<", eval2f("<", lt))
	e.AddCoreFun(">", eval2f(">", gt))
	e.AddCoreFun("<=", eval2f("<=", le))
	e.AddCoreFun(">=", eval2f(">=", ge))
}

func sum(acc, v float64) float64   { return acc + v }
func neg(n float64) data.Node      { return -n }
func diff(a, b float64) data.Node  { return a - b }
func prod(acc, v float64) float64  { return acc * v }
func reciproc(n float64) data.Node { return 1 / n }
func div(a, b float64) data.Node   { return a / b }

func lt(a, b float64) data.Node { return a < b }
func gt(a, b float64) data.Node { return a > b }
func le(a, b float64) data.Node { return a <= b }
func ge(a, b float64) data.Node { return a >= b }

func list(e Evaluator, env data.Env, args []data.Node) data.Node {
	return data.NewList(args)
}

func isList(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 1 {
		return e.Error("[list?] expects 1 argument.")
	}
	return data.IsList(args[0])
}

func count(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 1 {
		return e.Error("[count] expects 1 argument.")
	}
	switch x := args[0].(type) {
	case *data.ListNode:
		return float64(len(x.Items))
	case *data.VectorNode:
		return float64(len(x.Items))
	case *data.HashMapNode:
		return float64(len(x.Items))
	default:
		return e.Error("[%s] cannot be an argument to [count].", "")
	}
}

func isEmpty(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 1 {
		return e.Error("[empty?] expects 1 argument.")
	}
	switch x := args[0].(type) {
	case *data.ListNode:
		return len(x.Items) == 0
	case *data.VectorNode:
		return len(x.Items) == 0
	case *data.HashMapNode:
		return len(x.Items) == 0
	default:
		return e.Error("[%s] cannot be an argument to [empty?].", "")
	}
}

func eq(e Evaluator, env data.Env, args []data.Node) data.Node {
	if len(args) != 2 {
		return e.Error("[=] expects 2 arguments.")
	}
	return eq2(e, env, args[0], args[1])
}

func eq2(e Evaluator, env data.Env, a, b data.Node) data.Node {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch x := a.(type) {
	case *data.ListNode:
		y := b.(*data.ListNode)
		return eqSeq(e, env, x.Items, y.Items)
	case *data.VectorNode:
		y := b.(*data.VectorNode)
		return eqSeq(e, env, x.Items, y.Items)
	case *data.HashMapNode:
		y := b.(*data.HashMapNode)
		return eqHashMap(e, env, x.Items, y.Items)
	default:
		return a == b
	}
}

func eqSeq(e Evaluator, env data.Env, as, bs []data.Node) data.Node {
	if len(as) != len(bs) {
		return false
	}
	for i := 0; i < len(as); i++ {
		if eq2(e, env, as[i], bs[i]) == false {
			return false
		}
	}
	return true
}

func eqHashMap(e Evaluator, env data.Env, as, bs data.Map) data.Node {
	if len(as) != len(bs) {
		return false
	}
	for k, va := range as {
		vb, ok := bs[k]
		if !ok {
			return false
		}
		if eq2(e, env, va, vb) == false {
			return false
		}
	}
	return true
}

func evalxf(name string, f func(float64, float64) float64) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		var acc float64
		for i, arg := range args {
			if v, ok := arg.(float64); ok {
				acc = f(acc, v)
			} else {
				// TODO: Add a printer instance to the evaluator to print expressions.
				return e.Error("[%d]. argument [%s] is not a number.", i+1, "")
			}
		}
		return acc
	}
}

func eval12f(name string, f func(float64) data.Node, g func(float64, float64) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		switch len(args) {
		case 1:
			if n, ok := args[0].(float64); ok {
				return f(n)
			}
			// TODO: Add a printer instance to the evaluator to print expressions.
			return e.Error("Argument [%s] is not a number.", "")
		case 2:
			n1, ok1 := args[0].(float64)
			n2, ok2 := args[1].(float64)
			if ok1 && ok2 {
				return g(n1, n2)
			}
			if !ok1 {
				return e.Error("First argument [%s] is not a number.", "")
			}
			if !ok2 {
				return e.Error("Second argument [%s] is not a number.", "")
			}
		}
		return e.Error("[%s] requires either 1 or 2 arguments.", name)
	}
}

func eval2f(name string, f func(float64, float64) data.Node) CoreFun {
	return func(e Evaluator, env data.Env, args []data.Node) data.Node {
		switch len(args) {
		case 2:
			n0, ok0 := args[0].(float64)
			n1, ok1 := args[1].(float64)
			if ok0 && ok1 {
				return f(n0, n1)
			}
			if !ok0 {
				return e.Error("First argument [%s] is not a number.", "")
			}
			if !ok1 {
				return e.Error("Second argument [%s] is not a number.", "")
			}
		}
		return e.Error("[%s] expects 2 arguments.", name)
	}
}