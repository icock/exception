// Package exception defines and implements a Throwable interface
// with stacktrace information.
//
// Throwable is not meant to replace error in Go.
// Throwable is meant to represent unchecked exceptions.
//
// Exception
//
// The Exception struct implements Throwable interface.
// New exceptions can be defined by embedding Exception,
// if customizing methods of Throwable is not required.
//
// For example, a NullReference exception can be defined as following:
//
//     type NullReference struct {
//       *exception.Exception
//     }
//     func NewNullReference(message string, cause Throwable) *NullReference {
//       return &NullReference{exception.New(message, cause)}
//     }
//
// NullReference
//
// The above NullReference exception is included in this package.
//
package exception

import (
	"io"
)

// Throwable is like an error, but with stacktrace information.
type Throwable interface {
	// Throwable returns the cause of this throwable
	// or nil if the cause is nonexistent or unknown.
	Cause() Throwable
	// Message returns the detail message string of this throwable.
	Message() string
	// StackTrace provides the stack trace information.
	StackTrace() []uintptr
	// EliminatedStackTrace is like StackTrace but with recursive calls eliminated.
	// It does not eliminates cross recursions like `f(g(f(g(f()))))`.
	EliminatedStackTrace() []uintptr
	// SetStackTrace sets the stack trace information.
	SetStackTrace([]uintptr)
	// PrintStackTrace prints this throwable and its backtrace
	// to the specified print writer.
	PrintStackTrace(writer io.Writer)
	// Error returns a string representation of this throwable,
	// which is used by functions like fmt.Print and fmt.Printf("`%v`").
	// Any type satisfies Throwable also satisfies error.
	Error() string
}

