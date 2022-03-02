package errorz

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// GetCallers returns the raw stack trace from the error, or the current raw stack trace if not found.
func GetCallers(err error) []uintptr {
	if e, ok := err.(*wrappedError); ok {
		if e.callers != nil {
			return e.callers
		}
	}

	callers := make([]uintptr, 1024)
	return callers[:runtime.Callers(2, callers[:])]
}

// Skip skips the caller from the stack trace.
func Skip() OptionFunc {
	var callerFunc *runtime.Func
	if caller, _, _, ok := runtime.Caller(1); ok {
		callerFunc = runtime.FuncForPC(caller)
	}

	return func(err error) {
		if e, ok := err.(*wrappedError); callerFunc != nil && ok && e.callers != nil {
			for i, caller := range e.callers {
				if callerFunc == runtime.FuncForPC(caller) {
					e.callers = append(e.callers[:i], e.callers[i+1:]...)
					return
				}
			}
		}
	}
}

// SkipAll skips the caller and any lower frames from the stack trace.
func SkipAll() OptionFunc {
	var callerFunc *runtime.Func
	if caller, _, _, ok := runtime.Caller(1); ok {
		callerFunc = runtime.FuncForPC(caller)
	}

	return func(err error) {
		if e, ok := err.(*wrappedError); callerFunc != nil && ok && e.callers != nil {
			for i, caller := range e.callers {
				if callerFunc == runtime.FuncForPC(caller) {
					e.callers = e.callers[i+1:]
					return
				}
			}
		}
	}
}

// FormatStackTrace formats the given raw stack trace.
func FormatStackTrace(callers []uintptr) []string {
	frames := runtime.CallersFrames(callers)
	stackTrace := make([]string, 0, len(callers))

	for {
		frame, more := frames.Next()

		stackTrace = append(
			stackTrace,
			fmt.Sprintf("%v (%v:%v)", filepath.Base(frame.Function), frame.File, frame.Line))

		if !more {
			break
		}
	}

	return stackTrace
}
