package bin_test

import (
	"hash/fnv"
	"reflect"
	"testing"

	"github.com/mhoertnagl/splis2/internal/bin"
	"github.com/mhoertnagl/splis2/internal/vm"
)

func TestLinkLib(t *testing.T) {
	lb1 := bin.NewLib()
	lb1.Entries = []uint64{0}
	lb1.Fns = vm.ConcatVar(
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
	lb1.Code = vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 0),
		vm.Instr(vm.OpCall),
	)

	lb2 := bin.NewLib()
	lb2.Entries = []uint64{0}
	lb2.Fns = vm.ConcatVar(
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
	lb2.Code = vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpRef, 0),
		vm.Instr(vm.OpCall),
	)

	exp := bin.NewLib()
	exp.Entries = []uint64{0, 32}
	exp.Fns = vm.ConcatVar(
		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpAdd),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),

		vm.Instr(vm.OpNewEnv),
		vm.Instr(vm.OpSetLocal, hash("x")),
		vm.Instr(vm.OpPop),
		vm.Instr(vm.OpGetLocal, hash("x")),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpMul),
		vm.Instr(vm.OpPopEnv),
		vm.Instr(vm.OpReturn),
	)
	exp.Code = vm.ConcatVar(
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 1),
		vm.Instr(vm.OpRef, 0),
		vm.Instr(vm.OpCall),
		vm.Instr(vm.OpEnd),
		vm.Instr(vm.OpConst, 2),
		vm.Instr(vm.OpRef, 1),
		vm.Instr(vm.OpCall),
	)

	lnk := bin.NewLinker()
	lnk.Add(lb1)
	lnk.Add(lb2)
	assertLibsEqual(t, lnk.Lib(), exp)
}

func assertLibsEqual(t *testing.T, a *bin.Lib, e *bin.Lib) {
	t.Helper()
	if reflect.DeepEqual(a.Code, e.Code) == false {
		t.Errorf(
			"\nCode mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Code,
			e.Code,
		)
	}
	if reflect.DeepEqual(a.Entries, e.Entries) == false {
		t.Errorf(
			"\nEntries mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Entries,
			e.Entries,
		)
	}
	if reflect.DeepEqual(a.Fns, e.Fns) == false {
		t.Errorf(
			"\nFns mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Fns,
			e.Fns,
		)
	}
	if reflect.DeepEqual(a.Macros, e.Macros) == false {
		t.Errorf(
			"\nMacros mismatch:\n  Actual: %v\n  Expect: %v\n",
			a.Macros,
			e.Macros,
		)
	}
}

func hash(sym string) uint64 {
	hg := fnv.New64()
	hg.Reset()
	hg.Write([]byte(sym))
	return hg.Sum64()
}
