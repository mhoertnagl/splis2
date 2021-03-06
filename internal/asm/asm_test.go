package asm_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mhoertnagl/noodles/internal/asm"
	"github.com/mhoertnagl/noodles/internal/vm"
)

func TestAssembleBool1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleBool2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleInteger(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 42),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 42),
	)
	testa(t, i, e)
}

func TestAssembleRandom(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 100),
		asm.Instr(vm.OpRand),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 100),
		vm.Instr(vm.OpRand),
	)
	testa(t, i, e)
}

// --- IF ---

func TestAssembleIf1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Label("L0"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 19),
		vm.Instr(vm.OpConst, 1),
	)
	testa(t, i, e)
}

func TestAssembleIf2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 0),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 37),
		vm.Instr(vm.OpConst, 0),
	)
	testa(t, i, e)
}

func TestAssembleIf3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
	testa(t, i, e)
}

func TestAssembleIf4(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 42),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpConst, 21),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 46),
		vm.Instr(vm.OpConst, 42),
		vm.Instr(vm.OpJump, 55),
		vm.Instr(vm.OpConst, 21),
	)
	testa(t, i, e)
}

// --- COND ---

func TestAssembleCond1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L1"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 2),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L2"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpConst, 3),
		asm.Label("L0"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 28),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 56),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 75),
		vm.Instr(vm.OpConst, 3),
	)
	testa(t, i, e)
}

// --- AND ---

func TestAssembleAnd0(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpTrue),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleAnd1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleAnd2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleAnd3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIfNot, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpFalse),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIfNot, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

// --- OR ---

func TestAssembleOr0(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleOr1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
	)
	testa(t, i, e)
}

func TestAssembleOr2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 20),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJump, 21),
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

func TestAssembleOr3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpTrue),
		asm.Labeled(vm.OpJumpIf, "L0"),
		asm.Instr(vm.OpFalse),
		asm.Labeled(vm.OpJump, "L1"),
		asm.Label("L0"),
		asm.Instr(vm.OpTrue),
		asm.Label("L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpTrue),
		vm.Instr(vm.OpJumpIf, 30),
		vm.Instr(vm.OpFalse),
		vm.Instr(vm.OpJump, 31),
		vm.Instr(vm.OpTrue),
	)
	testa(t, i, e)
}

// --- VECTORS ---

func TestAssembleVectorConcat(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpConcat),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpConcat),
	)
	testa(t, i, e)
}

// --- SET ---

func TestAssembleSet2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpJump, 58),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 27),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

// --- LET ---

func TestAssembleLet5(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		// BEGIN FN
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCall),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		// END FN
		asm.Label("L0"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L1"),
		// BEGIN LET
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpDropArgs, 1),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 96),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 65),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 95),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 9),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpDropArgs, 1),
	)
	testa(t, i, e)
}

// --- FN ---

func TestAssembleAnonymousFun1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 39),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun11(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Ref(0, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 67),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 49),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 19),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 50),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 19),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun3(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpAdd),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 60),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 29),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpAdd),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousFun4(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Ref(0, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpJump, 79),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 61),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 30),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 20),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleAnonymousNestedFun(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Labeled(vm.OpJump, "L2"),
		// (fn [m] (/ n m))
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpReturn),
		// end
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpJump, 106),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpJump, 78),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 48),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 19),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleLeafFunDef(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpAdd),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 40),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpAdd),
	)
	testa(t, i, e)
}

func TestAssembleVariadicFun1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 0),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpPushArgs, 0),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 46),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleVariadicFun2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 4),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpConst, 2),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpList),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpCons),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 4),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 85),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpList),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpCons),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 46),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleFac(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpConst, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpMul),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 107),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 65),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpJump, 106),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleTailFac(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 0),
		asm.Instr(vm.OpEQ),
		asm.Labeled(vm.OpJumpIfNot, "L2"),
		asm.Instr(vm.OpGetArg, 1),
		asm.Labeled(vm.OpJump, "L3"),
		asm.Label("L2"),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpMul),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpSub),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpRecCall),
		asm.Label("L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Labeled(vm.OpJump, "L4"),
		asm.Label("L5"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L4"),
		asm.Ref(0, "L5"),
		asm.Instr(vm.OpSetGlobal, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 5),
		asm.Instr(vm.OpGetGlobal, 1),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 116),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 0),
		vm.Instr(vm.OpEQ),
		vm.Instr(vm.OpJumpIfNot, 65),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpJump, 115),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpSub),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpRecCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpJump, 191),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 151),
		vm.Instr(vm.OpSetGlobal, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 5),
		vm.Instr(vm.OpGetGlobal, 1),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleClosure1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 75),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 48),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 28),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleClosure2(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpConst, 1),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpAdd),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 1),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpCall),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 6),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 115),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpJump, 87),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpRef, 1, 56),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 6),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

func TestAssembleClosure(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Labeled(vm.OpJump, "L0"),
		asm.Label("L1"),
		asm.Instr(vm.OpPushArgs, 1),
		asm.Instr(vm.OpPop),
		asm.Labeled(vm.OpJump, "L2"),
		asm.Label("L3"),
		asm.Instr(vm.OpPushArgs, 2),
		asm.Instr(vm.OpPop),
		asm.Instr(vm.OpGetArg, 1),
		asm.Instr(vm.OpGetArg, 0),
		asm.Instr(vm.OpDiv),
		asm.Instr(vm.OpReturn),
		asm.Label("L2"),
		asm.Instr(vm.OpGetArg, 0),
		asm.Ref(1, "L3"),
		asm.Instr(vm.OpReturn),
		asm.Label("L0"),
		// def divN
		asm.Ref(0, "L1"),
		asm.Instr(vm.OpSetGlobal, 0),
		// ((divN 3) 9 )
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 9),
		// (divN 3)
		asm.Instr(vm.OpEnd),
		asm.Instr(vm.OpConst, 3),
		asm.Instr(vm.OpGetGlobal, 0),
		asm.Instr(vm.OpCall),
		// --- (divN 3)
		asm.Instr(vm.OpCall),
		// --- ((divN 3) 9 )
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpJump, 85),
		vm.Instr(vm.OpPushArgs, 1),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpJump, 58),
		vm.Instr(vm.OpPushArgs, 2),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetArg, 1),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpDiv),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpGetArg, 0),
		vm.Instr(vm.OpRef, 1, 28),
		vm.Instr(vm.OpReturn),
		vm.Instr(vm.OpRef, 0, 9),
		vm.Instr(vm.OpSetGlobal, 0),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 9),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 3),
		vm.Instr(vm.OpGetGlobal, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpCall),
	)
	testa(t, i, e)
}

// --- STRING ---

func TestAssembleString(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Str("Hello, World!"),
	}
	e := vm.ConcatVar(
		vm.Str("Hello, World!"),
	)
	testa(t, i, e)
}

// --- WRITE ---

func TestAssembleWrite1(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpEnd),
		asm.Str("Hello, World!\n"),
		asm.Instr(vm.OpGetGlobal, 1),
		asm.Instr(vm.OpWrite),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Str("Hello, World!\n"),
		vm.Instr(vm.OpGetGlobal, 1),
		vm.Instr(vm.OpWrite),
	)
	testa(t, i, e)
}

// --- HALT ---

func TestAssembleHalt(t *testing.T) {
	i := []asm.AsmCmd{
		asm.Instr(vm.OpHalt),
	}
	e := vm.ConcatVar(
		vm.Instr(vm.OpHalt),
	)
	testa(t, i, e)
}

func testa(t *testing.T, i asm.AsmCode, e []byte) {
	t.Helper()
	a := asm.NewAssembler()
	ib := a.Assemble(i)
	compareAssembly(t, ib, e)
}

func compareAssembly(t *testing.T, a []byte, e []byte) {
	t.Helper()

	err := false

	d := asm.NewDisassembler()
	da := d.Disassemble(a)
	de := d.Disassemble(e)
	la := len(da)
	le := len(de)

	lm := la
	if le > la {
		lm = le
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("     %-20s%s\n", "Actual", "Expecting"))
	buf.WriteString(fmt.Sprintf("     %-20s%s\n", "------", "---------"))

	for i := 0; i < lm; i++ {
		sa := ""
		if i < la {
			sa = da[i]
		}
		se := ""
		if i < le {
			se = de[i]
		}

		buf.WriteString(fmt.Sprintf("%3d: %-20s%-20s", i, sa, se))

		if sa != se {
			err = true
			buf.WriteString("<--")
		}

		buf.WriteString("\n")
	}

	if err {
		t.Errorf("\n%s\n", buf.String())
	}
}
