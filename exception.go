package exception

import (
	"fmt"
	"io"
	"runtime"
)

// Exception implements Throwable.
type Exception struct {
	cause   Throwable
	message string
	stackTrace []uintptr
}

var _ Throwable = (*Exception)(nil)

func (e *Exception) Cause() Throwable {
	return e.cause
}

func (e *Exception) Message() string {
	return e.message
}

func (e *Exception) StackTrace() []uintptr {
	return e.stackTrace
}

func (e *Exception) SetStackTrace(stackTrace []uintptr) {
	e.stackTrace = stackTrace
}

func (e *Exception) PrintStackTrace(w io.Writer) {
	var frames *runtime.Frames = runtime.CallersFrames(e.StackTrace())
	fmt.Fprint(w, frames)
}

func (e *Exception) Error() string {
	return "Exception: " + e.Message()
}

func (e *Exception) EliminatedStackTrace() []uintptr {
	var st []uintptr = e.StackTrace()
	var stDepth int = len(st)
	eliminated := make([]uintptr, 0, stDepth)
	for i := 0; i < stDepth-1; i++ {
		j := i+1
		for j < stDepth && st[i] == st[j] {
			j++
		}
		if j == stDepth {
			eliminated = append(eliminated, st[:i+1]...)
		} else if j < stDepth {
			if j > i+1 {
				eliminated = append(eliminated, st[:i+1]...)
				eliminated = append(eliminated, st[j:]...)
			}
		} else {
			panic(fmt.Sprintf(
				"BUG: EliminatedStackTrace(): i (%d) j (%d) exceed bounds (stDepth %d)",
				i, j, stDepth))
		}
	}
	return eliminated
}

// New constructs a new Exception with the specified detail message and cause.
func New(message string, cause Throwable) *Exception {
	var pc []uintptr = callers()
	return &Exception{
		cause:   cause,
		message: message,
		stackTrace: pc,
	}
}

// NewWithStackTraceDepthLimit is like New, but with a specified depth limit of stacktrace.
// It is recommended only using NewWithStackTraceDepthLimit
// when
//     0. the stacktrace will be very deep but only first levels are interested, or
//     1. the needed depth of callers is known or can be estimated, and
//        performance issues with New are encountered.
func NewWithStackTraceDepthLimit(message string, cause Throwable, limit int) *Exception {
	var pc []uintptr
	_, pc = implorer(limit, -1)
	return &Exception{
		cause:   cause,
		message: message,
		stackTrace: pc,
	}
}

func callers() []uintptr {
	var pc []uintptr
	// 16 is deep enough under most conditions, not considering recursive calls.
	// JavaScript V8 has a default stacktrace limit of 10.
	// See also comments in implorer definition
	// (in else branch, above recursive calls of implorer).
	_, pc = implorer(16, 0)
	return pc
}

// If recursion is negative, then implorer will call runtime.Callers just once,
// not taking care of depth limit exceeding.
func implorer(depthLimit int, recursion int) (int, []uintptr) {
	// 10 is the stacktrace depth limit of JavaScript V8.
	// We use it as the initial capacity of the slice.
	// The stacktrace depth is unlimited though.
	var pc = make([]uintptr, depthLimit)

	// Suppose a function f throw an exception, i.e. calling exception.New,
	// then the implorer are (assuming depthLimit not exceeded)
	//
	//     0: runtime.Callers
	//     1: implorer (this function)
	//     2: callers (wraps implorer to expose a simple api)
	//     3: newException (returns an exception)
	//     4: exception.New (returns a Throwable)
	//     5: f
	//
	// Thus we skip first five implorer.
	//
	// When depth exceeds depthLimit,
	// we increase depthLimit and call implorer recursively,
	// thus we skip additional implorer recursion
	depth := runtime.Callers(5+recursion, pc)

	if recursion < 0 {
		return depth, pc[0:depth]
	} else if isFull(depth, depthLimit) {
		// uintptr size on 64-bit machine is 8.
		// Given the initial depthLimit 16,
		// the second call will allocate an a 16 KiB array,
		// and the third call will allocate a 2 MiB array,
		// within the size of L1 and L3 cache of today's commodity machines.
		const depth_limit_times = 1024/8
		return implorer(depthLimit*(depth_limit_times), recursion+1)
	} else {
		// Let gc free up the underlying array filled up with zero values.
		var result = make([]uintptr, depth)
		copy(result, pc)
		return depth, result
	}
}

func isFull(length int, limit int) bool {
	if length < limit {
		return false
	} else if length == limit {
		return true
	} else {
		message := fmt.Sprintf(
			"isFull(length: %d, limit: %d): length should not exceed limit",
			length, limit)
		panic(message)
	}
}


// FromError converts an error (err) to an exception,
// only utilizing the `err.Error()` method.
func FromError(err error) *Exception {
	return New(err.Error(), nil)
}

