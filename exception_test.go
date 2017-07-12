package exception

import (
	"testing"
	"fmt"
	"github.com/pkg/errors"
)

func TestFromError(t *testing.T) {
	err := errors.New("a random error")
	var e Throwable = FromError(err)
	if e.Cause() != nil {
		t.Error("Exceptions converted via FromError should have a nil cause.")
	}
	if e.Message() != err.Error() {
		t.Error("Exceptions converted via FromError should use original err.Error() as its message.")
	}
	expected := "Exception: a random error"
	if e.Error() != expected {
		t.Error(
			"e.Error() FAIL\n  Actual Value: %v\nExpected Value: %s",
			e.Error(), expected)
	}
}

func TestNew(t *testing.T) {
	var e Throwable = New("an exception", nil)
	t.Log(e.Error())
}

func g(n int) Throwable {
	if n == 0 {
		return f()
	} else {
		return g(n-1)
	}
}

func f() Throwable {
	return FromError(errors.New("from f"))
}

func ExampleNewWithStackTraceDepthLimit() {
	var e Throwable = NewWithStackTraceDepthLimit("an exception", nil, 10)
	fmt.Println(e.StackTrace())
	// Output: FIXME
}

func ExampleNew() {
	var e Throwable = New("an exception", nil)
	fmt.Println(e.StackTrace())
	// Output: FIXME
}

func ExampleRecursiveCalls() {
	var e Throwable = g(3)
	fmt.Println(e.StackTrace())
	// Output: FIXME
}



