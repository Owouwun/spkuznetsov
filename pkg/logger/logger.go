package logger

import (
	"reflect"
	"testing"
)

func LogIfErr(t *testing.T, comment string, f interface{}, opts ...interface{}) {
	fValue := reflect.ValueOf(f)
	if fValue.Kind() != reflect.Func {
		t.Logf("Error: 'f' is not a function")
		return
	}

	args := make([]reflect.Value, len(opts))
	for i, arg := range opts {
		args[i] = reflect.ValueOf(arg)
	}

	results := fValue.Call(args)
	if len(results) != 1 {
		t.Logf("Error: must has 1 result value (error)")
		return
	}

	if err, ok := results[0].Interface().(error); ok && err != nil {
		t.Logf(comment, err)
	}
}
