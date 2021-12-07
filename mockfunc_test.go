package mockfunc_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sharefull/mockfunc"
)

type testingT struct {
	fatal []interface{}
}

func (t *testingT) Helper() {}

func (t *testingT) Fatal(args ...interface{}) {
	t.fatal = args
}

func (t *testingT) Fatalf(format string, args ...interface{}) {
	t.fatal = append([]interface{}{format}, args...)
}

func TestSet(t *testing.T) {
	t.Parallel()

	type __ = mockfunc.Unused
	type T struct{ F func(a, b int) (c, d int) }
	V := func(vs ...interface{}) []interface{} { return vs }

	// default value of test data
	_func := func(_, _ int) (_, _ int) { return }
	_args := V(0, 0)
	_result := V(0, 0)

	cases := []struct {
		name        string
		dst         interface{}
		fun         interface{}
		args        []interface{}
		want        []interface{}
		expectFatal bool
	}{
		{"ok", &(new(T).F), _func, _args, _result, false},
		{"ng:dstnil", nil, _func, _args, _result, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var _t testingT
			mockfunc.Set(&_t, tt.dst, tt.fun)

			switch {
			case tt.expectFatal && _t.fatal == nil:
				t.Fatal("expected fatal did not occur")
			case !tt.expectFatal && _t.fatal != nil:
				t.Fatal("unexpected fatal", _t.fatal)
			case _t.fatal != nil:
				return
			}

			args := make([]reflect.Value, len(tt.args))
			for i := range tt.args {
				args[i] = reflect.ValueOf(tt.args[i])
			}
			results := reflect.ValueOf(tt.dst).Elem().Call(args)

			got := make([]interface{}, len(results))
			for i := range results {
				got[i] = results[i].Interface()
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func ExampleSet() {
	type __ = mockfunc.Unused
	___ := mockfunc.UnusedValue

	type Mock struct {
		DoFunc func(ctx context.Context, id int) (n int, err error)
	}

	var m Mock
	t := new(testing.T) // dummy
	mockfunc.Set(t, &m.DoFunc, func(_ __, id int) (__, error) {
		if id%2 == 0 {
			return ___, errors.New("error")
		}
		return ___, nil
	})

	fmt.Println(m.DoFunc(context.Background(), 1))
	fmt.Println(m.DoFunc(context.Background(), 2))

	// Output:
	// 0 <nil>
	// 0 error
}
