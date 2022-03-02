package errorz

import (
	"fmt"
	"io"
	"runtime"
)

type wrappedError struct {
	err      error
	id       ID
	status   Status
	metadata Metadata
	prefix   string
	callers  []uintptr
}

// Error implements the error interface.
func (e *wrappedError) Error() string {
	return e.prefix + e.err.Error()
}

// Wrap wraps the given error, applying the given options.
func Wrap(err error, options ...Option) error {
	if err == nil {
		panic("nil error")
	}

	e, ok := err.(*wrappedError)
	if !ok {
		e = &wrappedError{
			err:      err,
			metadata: Metadata{},
			callers:  make([]uintptr, 1024),
		}
		e.callers = e.callers[:runtime.Callers(2, e.callers[:])]
	}

	for _, option := range options {
		option.Apply(e)
	}

	return e
}

// MaybeWrap is like Wrap, but returns nil if called with a nil error.
func MaybeWrap(err error, options ...Option) error {
	if err == nil {
		return nil
	}

	options = append(options, Skip())
	return Wrap(err, options...)
}

// MustWrap is like Wrap, but panics if the given error is non-nil.
func MustWrap(err error, options ...Option) {
	if err == nil {
		panic("nil error")
	}

	options = append(options, Skip())
	panic(Wrap(err, options...))
}

// MaybeMustWrap is like MustWrap, but does nothing if called with a nil error.
func MaybeMustWrap(err error, options ...Option) {
	if err == nil {
		return
	}

	options = append(options, Skip())
	MustWrap(err, options...)
}

// WrapRecover takes a recovered interface{} and converts it to a wrapped error.
func WrapRecover(r interface{}, options ...Option) error {
	if r == nil {
		panic("nil recover")
	}

	options = append(options, Skip())

	switch r := r.(type) {
	case *wrappedError:
		return r
	case error:
		return Wrap(r, options...)
	default:
		return Wrap(fmt.Errorf("%v", r), options...)
	}
}

// MaybeWrapRecover is like WrapRecover but returns nil if called with a nil recover.
func MaybeWrapRecover(r interface{}, options ...Option) error {
	if r == nil {
		return nil
	}

	options = append(options, Skip())
	return WrapRecover(r, options...)
}

// Errorf formats a new error and wraps it.
// Note: arguments implementing Option are applied on wrapping, the others are passed to fmt.Errorf().
func Errorf(format string, options ...Option) error {
	var mergedArgs []interface{}
	for _, option := range options {
		if args, ok := option.(Args); ok {
			mergedArgs = append(mergedArgs, args...)
		}
	}

	options = append(options, Skip())
	return Wrap(fmt.Errorf(format, mergedArgs...), options...)
}

// MustErrorf is like Errorf but panics instead of returning the error.
func MustErrorf(format string, options ...Option) {
	options = append(options, Skip())
	panic(Errorf(format, options...))
}

// Assertf is like MustErrorf if cond is false, does nothing otherwise.
func Assertf(cond bool, format string, options ...Option) {
	if cond {
		return
	}

	options = append(options, Skip())
	MustErrorf(format, options...)
}

// IgnoreClose calls Close on the given io.Closer, ignoring the returned error. Handy for the defer Close pattern.
func IgnoreClose(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

// MustClose calls Close on the given io.Closer, panicking in case of error. Handy for the defer Close pattern.
func MustClose(c io.Closer) {
	if c != nil {
		MaybeMustWrap(c.Close(), Skip())
	}
}

// Unwrap undoes Wrap, returning the original error.
func Unwrap(err error) error {
	if wErr, ok := err.(*wrappedError); ok {
		return wErr.err
	}
	return err
}

// Safe calls the function catching any panic and returning it as error.
func Safe(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if rErr := MaybeWrapRecover(recover()); rErr != nil {
				err = rErr
			}
		}()

		return MaybeWrap(f())
	}
}
