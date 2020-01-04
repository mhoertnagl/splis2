package compiler

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/mhoertnagl/splis2/internal/vm"
)

type Compiler interface {
	Compile(node Node) vm.Ins
}

type compiler struct {
	hg hash.Hash64
}

func NewCompiler() Compiler {
	return &compiler{
		hg: fnv.New64(),
	}
}

func (c *compiler) Compile(node Node) vm.Ins {
	switch n := node.(type) {
	case bool:
		return c.compileBooleanLiteral(n)
	case int64:
		return c.compileIntegerLiteral(n)
	case *SymbolNode:
		return c.compileSymbol(n)
	case *ListNode:
		return c.compileList(n)
	}
	return nil
}

func (c *compiler) compileBooleanLiteral(n bool) vm.Ins {
	if n {
		return vm.Instr(vm.OpTrue)
	}
	return vm.Instr(vm.OpFalse)
}

func (c *compiler) compileIntegerLiteral(n int64) vm.Ins {
	return vm.Instr(vm.OpConst, uint64(n))
}

func (c *compiler) compileSymbol(n *SymbolNode) vm.Ins {
	return vm.Instr(vm.OpGetLocal, c.hashSymbol(n))
}

func (c *compiler) compileList(n *ListNode) vm.Ins {
	items := n.Items
	if len(items) == 0 {
		panic("Empty list")
	}
	args := items[1:]
	switch sym := items[0].(type) {
	case *SymbolNode:
		switch sym.Name {
		case "+":
			return c.compileAdd(args)
		case "-":
			return c.compileSub(args)
		case "*":
			return c.compileMul(args)
		case "/":
			return c.compileDiv(args)
		case "let":
			return c.compileLet(args)
		case "def":
			return c.compileDef(args)
		case "if":
			return c.compileIf(args)
		case "do":
			return c.compileDo(args)
		default:
			panic(fmt.Sprintf("Cannot compile core function [%v]", sym))
		}
	default:
		panic(fmt.Sprintf("Cannot compile list head [%v]", sym))
	}
}

func (c *compiler) compileAdd(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty sum (+) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton sum (+ x) yields x.
		return c.Compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// addition operations except for the first pair. The resulting sequence of
		// instructions is then:
		//
		//   <(+ x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpAdd, <x3>, OpAdd, <x4>, OpAdd, ...
		//
		code := vm.NewCodeGen()
		code.Append(c.Compile(args[0]))
		for _, arg := range args[1:] {
			code.Append(c.Compile(arg))
			code.Instr(vm.OpAdd)
		}
		return code.Emit()
	}
}

func (c *compiler) compileSub(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty difference (-) yields 0.
		return vm.Instr(vm.OpConst, 0)
	case 1:
		// Singleton difference (- x) yields (- 0 x) which if effectively -x.
		code := vm.NewCodeGen()
		code.Instr(vm.OpConst, 0)
		code.Append(c.Compile(args[0]))
		code.Instr(vm.OpSub)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their difference.
		//
		//   <(- x1 x2)> := <x1>, <x2>, OpSub
		//
		code := vm.NewCodeGen()
		code.Append(c.Compile(args[0]))
		code.Append(c.Compile(args[1]))
		code.Instr(vm.OpSub)
		return code.Emit()
	default:
		panic("Too many arguments")
	}
}

func (c *compiler) compileMul(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty product (*) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton product (* x) yields x.
		return c.Compile(args[0])
	default:
		// Will compile this expression to a sequence of compiled subexpressions and
		// multiplication operations except for the first pair. The resulting
		// sequence of instructions is then:
		//
		//   <(* x1 x2 x3 x4 ...)> :=
		//     <x1>, <x2>, OpMul, <x3>, OpMul, <x4>, OpMul, ...
		//
		code := vm.NewCodeGen()
		code.Append(c.Compile(args[0]))
		for _, arg := range args[1:] {
			code.Append(c.Compile(arg))
			code.Instr(vm.OpMul)
		}
		return code.Emit()
	}
}

func (c *compiler) compileDiv(args []Node) vm.Ins {
	switch len(args) {
	case 0:
		// Empty division (/) yields 1.
		return vm.Instr(vm.OpConst, 1)
	case 1:
		// Singleton difference (/ x) yields (/ 1 x) which if effectively 1/x.
		code := vm.NewCodeGen()
		code.Instr(vm.OpConst, 1)
		code.Append(c.Compile(args[0]))
		code.Instr(vm.OpDiv)
		return code.Emit()
	case 2:
		// Only supports at most two operands and computes their quotient.
		//
		//   <(/ x1 x2)> := <x1>, <x2>, OpDiv
		//
		code := vm.NewCodeGen()
		code.Append(c.Compile(args[0]))
		code.Append(c.Compile(args[1]))
		code.Instr(vm.OpDiv)
		return code.Emit()
	default:
		panic("Too many arguments")
	}
}

func (c *compiler) compileLet(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[let] requires exactly two arguments")
	}
	if bs, ok := args[0].(*ListNode); ok {
		if len(bs.Items)%2 == 1 {
			panic("[let] reqires an even number of bindings")
		}
		code := vm.NewCodeGen()
		code.Instr(vm.OpNewEnv)
		for i := 0; i < len(bs.Items); i += 2 {
			if sym, ok2 := bs.Items[i].(*SymbolNode); ok2 {
				code.Append(c.Compile(bs.Items[i+1]))
				hsh := c.hashSymbol(sym)
				code.Instr(vm.OpSetLocal, hsh)
			} else {
				panic(fmt.Sprintf("[let] cannot bind to [%v]", sym))
			}
		}
		code.Append(c.Compile(args[1]))
		code.Instr(vm.OpPopEnv)
		return code.Emit()
	}
	panic("[let] requires first argument to be a list of bindings")
}

func (c *compiler) compileDef(args []Node) vm.Ins {
	if len(args) != 2 {
		panic("[def] requires exactly two arguments")
	}
	if sym, ok := args[0].(*SymbolNode); ok {
		code := vm.NewCodeGen()
		code.Append(c.Compile(args[1]))
		hsh := c.hashSymbol(sym)
		code.Instr(vm.OpSetGlobal, hsh)
		return code.Emit()
	}
	panic("[def] requires first argument to be a symbol")
}

func (c *compiler) compileIf(args []Node) vm.Ins {
	if len(args) != 2 && len(args) != 3 {
		panic("[if] requires exactly two or three arguments")
	}
	code := vm.NewCodeGen()
	switch len(args) {
	case 2:
		cnd := c.Compile(args[0])
		cns := c.Compile(args[1])
		cnsLen := uint64(len(cns))
		code.Append(cnd)
		code.Instr(vm.OpJumpIfNot, cnsLen)
		code.Append(cns)
	case 3:
		cnd := c.Compile(args[0])
		cns := c.Compile(args[1])
		alt := c.Compile(args[2])
		cnsLen := uint64(len(cns)) + 9 // Add the length of the jmp instruction.
		altLen := uint64(len(alt))
		code.Append(cnd)
		code.Instr(vm.OpJumpIfNot, cnsLen)
		code.Append(cns)
		code.Instr(vm.OpJump, altLen)
		code.Append(alt)
	}
	return code.Emit()
}

func (c *compiler) compileDo(args []Node) vm.Ins {
	code := vm.NewCodeGen()
	for _, arg := range args {
		code.Append(c.Compile(arg))
	}
	return code.Emit()
}

// func (c *compiler) compileCond(args []Node) vm.Ins {
//
// }

func (c *compiler) hashSymbol(sym *SymbolNode) uint64 {
	c.hg.Reset()
	c.hg.Write([]byte(sym.Name))
	return c.hg.Sum64()
}
