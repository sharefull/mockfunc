package mockfunc

import (
	"reflect"
	"testing"
)

type Unused struct{}

var typUnused = reflect.TypeOf(Unused{})

type TestingT interface {
	Helper()
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

var _ TestingT = (*testing.T)(nil)

func Set(t TestingT, dst, fun interface{}) {
	t.Helper()
	dstV := reflect.ValueOf(dst)
	funV := reflect.ValueOf(fun)

	if dstV.Kind() != reflect.Ptr ||
		dstV.Elem().Kind() != reflect.Func ||
		!dstV.Elem().CanSet() {
		t.Fatal("dst must be a pointer of function")
		return
	}

	if funV.Kind() != reflect.Func {
		t.Fatal("fun must be a function")
		return
	}

	dstT := dstV.Elem().Type()
	funT := funV.Type()
	switch {
	case dstT.NumIn() != funT.NumIn():
		t.Fatal("The number of arguments of dst and fun must be same")
	case dstT.NumOut() != funT.NumOut():
		t.Fatal("The number of results of dst and fun must be same")
	}

	fn := func(_args []reflect.Value) (_results []reflect.Value) {
		args := make([]reflect.Value, len(_args))
		for i := range _args {
			switch {
			case funT.In(i) == typUnused:
				args[i] = reflect.Zero(typUnused)
			case dstT.In(i) != funT.In(i):
				t.Fatalf("The %d-th argument is different between dst and fun: %v vs %v", i, dstT.In(i), funT.In(i))
			default:
				args[i] = _args[i]
			}
		}

		results := funV.Call(args)
		_results = make([]reflect.Value, len(results))
		for i := range results {
			switch {
			case results[i].Type() == typUnused:
				_results[i] = reflect.Zero(dstT.Out(i))
			case dstT.Out(i) != funT.Out(i):
				t.Fatalf("The %d-th result is different between dst and fun: %v vs %v", i, dstT.Out(i), funT.Out(i))
			default:
				_results[i] = results[i]
			}
		}

		return _results
	}

	v := reflect.MakeFunc(dstT, fn)
	dstV.Elem().Set(v)
}
