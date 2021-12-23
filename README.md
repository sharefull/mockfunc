# mockfunc

`mockfunc` creates a function and set to a variable or field.
The function can use place holders instead of unused parameters.
It helps to reduce unnecessary code from your tests.

```go
type __ = mockfunc.Unused
___ := mockfunc.UnusedValue

type Mock struct {
	DoFunc func(ctx context.Context, id int) (n int, err error)
}

var m Mock
t := new(testing.T) // dummy
mockfunc.Set(t, &m.DoFunc, func(_ __, id int) (__, error) {
	if id%2 == 0 {
		return ___, errors.New("error") // return a zero value (0) and an error
	}
	return ___, nil // return a zero value (0) and nil
})

fmt.Println(m.DoFunc(context.Background(), 1))
fmt.Println(m.DoFunc(context.Background(), 2))

// Output:
// 0 <nil>
// 0 error
```
