package errorz_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ibrt/golang-errors/errorz"

	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	err := errorz.Wrap(fmt.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.TestWrap"))
	require.PanicsWithValue(t, "nil error", func() { _ = errorz.Wrap(nil) })
	require.Equal(t, err, errorz.Wrap(err))
	require.Equal(t, errorz.ID(""), errorz.GetID(err))
	require.Equal(t, errorz.Metadata{}, errorz.GetMetadata(err))
}

func TestMaybeWrap(t *testing.T) {
	err := errorz.MaybeWrap(fmt.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.TestMaybeWrap"))
	require.Nil(t, errorz.MaybeWrap(nil))
}

func TestMustWrap(t *testing.T) {
	require.PanicsWithError(t, "test error", func() { errorz.MustWrap(fmt.Errorf("test error")) })
	require.PanicsWithValue(t, "nil error", func() { errorz.MustWrap(nil) })
}

func TestMaybeMustWrap(t *testing.T) {
	require.PanicsWithError(t, "test error", func() { errorz.MaybeMustWrap(fmt.Errorf("test error")) })
	require.NotPanics(t, func() { errorz.MaybeMustWrap(nil) })
}

func TestWrapRecover(t *testing.T) {
	err := errorz.WrapRecover("test error")
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	err = errorz.WrapRecover(fmt.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.NotNil(t, err)
	err = errorz.WrapRecover(errorz.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.NotNil(t, err)
	require.PanicsWithValue(t, "nil recover", func() { _ = errorz.WrapRecover(nil) })
}

func TestMaybeWrapRecover(t *testing.T) {
	err := errorz.MaybeWrapRecover("test error")
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.NotNil(t, err)
	err = errorz.MaybeWrapRecover(fmt.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	err = errorz.MaybeWrapRecover(errorz.Errorf("test error"))
	require.NotNil(t, err)
	require.Equal(t, "test error", err.Error())
	require.NotNil(t, err)
	require.Nil(t, errorz.MaybeWrapRecover(nil))
}

func TestErrorf(t *testing.T) {
	err := errorz.Errorf("test error")
	require.NotNil(t, t, err)
	require.Equal(t, "test error", err.Error())
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.TestErrorf"))
	err = errorz.Errorf("format %s", errorz.Prefix("prefix"), errorz.A("xxx"), errorz.ID("id"), errorz.M("k", "v"))
	require.NotNil(t, t, err)
	require.Equal(t, "prefix: format xxx", err.Error())
	require.Equal(t, errorz.ID("id"), errorz.GetID(err))
	require.Equal(t, errorz.Metadata{"k": "v"}, errorz.GetMetadata(err))
}

func TestMustErrorf(t *testing.T) {
	require.PanicsWithError(t, "test error", func() { errorz.MustErrorf("test error") })
}

func TestAssert(t *testing.T) {
	require.NotPanics(t, func() { errorz.Assertf(true, "test error") })
	require.PanicsWithError(t, "test error: value", func() { errorz.Assertf(false, "test error: %v", errorz.Args{"value"}) })
}

type testCloser struct {
	fail   bool
	closed bool
}

// Close implements io.Closer.
func (c *testCloser) Close() error {
	c.closed = true
	if c.fail {
		return errorz.Errorf("close error")
	}
	return nil
}

func TestIgnoreClose(t *testing.T) {
	tc := &testCloser{}
	require.False(t, tc.closed)
	errorz.IgnoreClose(tc)
	require.True(t, tc.closed)

	tc = &testCloser{fail: true}
	require.False(t, tc.closed)
	errorz.IgnoreClose(tc)
	require.True(t, tc.closed)
}

func TestMustClose(t *testing.T) {
	tc := &testCloser{}
	require.False(t, tc.closed)
	require.NotPanics(t, func() {
		errorz.MustClose(tc)
	})
	require.True(t, tc.closed)

	tc = &testCloser{fail: true}
	require.False(t, tc.closed)
	require.PanicsWithError(t, "close error", func() {
		errorz.MustClose(tc)
	})
	require.True(t, tc.closed)
}

func TestUnwrap(t *testing.T) {
	require.Nil(t, errorz.Unwrap(nil))
	err := fmt.Errorf("test error")
	ret := errorz.Unwrap(err)
	require.Equal(t, ret, err)
	ret = errorz.Unwrap(errorz.Wrap(err))
	require.Equal(t, ret, err)
}

func TestSafe(t *testing.T) {
	require.EqualError(t, errorz.Safe(func() error { panic(errorz.Errorf("test error")) })(), "test error")
	require.EqualError(t, errorz.Safe(func() error { return errorz.Errorf("test error") })(), "test error")
}

func TestID(t *testing.T) {
	require.Equal(t, errorz.ID(""), errorz.GetID(errorz.Errorf("test error")))
	require.Equal(t, errorz.ID(""), errorz.GetID(fmt.Errorf("test error")))
	err := errorz.Errorf("test error", errorz.ID("id"))
	require.NotNil(t, err)
	require.Equal(t, errorz.ID("id"), errorz.GetID(err))
	require.Equal(t, "id", errorz.ID("id").String())
	err = errorz.Errorf("test error", errorz.ID("id"))
	require.NotNil(t, err)
	require.Equal(t, errorz.ID("id"), errorz.GetID(err))
}

func TestStatus(t *testing.T) {
	require.Equal(t, errorz.Status(0), errorz.GetStatus(errorz.Errorf("test error")))
	require.Equal(t, errorz.Status(0), errorz.GetStatus(fmt.Errorf("test error")))
	err := errorz.Errorf("test error", errorz.Status(50))
	require.NotNil(t, err)
	require.Equal(t, errorz.Status(50), errorz.GetStatus(err))
	err = errorz.Errorf("test error", errorz.Status(http.StatusNotFound))
	require.NotNil(t, err)
	require.Equal(t, errorz.Status(http.StatusNotFound), errorz.GetStatus(err))
	require.Equal(t, 50, errorz.Status(50).Int())
}

func TestMetadata(t *testing.T) {
	require.Equal(t, errorz.Metadata{}, errorz.GetMetadata(errorz.Errorf("test error")))
	require.Equal(t, errorz.Metadata{}, errorz.GetMetadata(fmt.Errorf("test error")))
	err := errorz.Errorf("test error",
		errorz.Metadata{"k1": "v1", "k2": 2},
		errorz.M("k3", "v3"),
		errorz.Metadata{"k4": "v4", "k5": "v5"},
		errorz.Metadata{},
		errorz.M("k5", "over"))
	require.Equal(t, errorz.Metadata{
		"k1": "v1",
		"k2": 2,
		"k3": "v3",
		"k4": "v4",
		"k5": "over",
	}, errorz.GetMetadata(err))
	require.Equal(t, "v1", errorz.GetMetadata(err).Get("k1"))
	require.Equal(t, nil, errorz.GetMetadata(err).Get("unknown"))
	require.Equal(t, "v1", errorz.GetMetadata(err).GetString("k1"))
	require.Equal(t, "", errorz.GetMetadata(err).GetString("k2"))
	require.Equal(t, "", errorz.GetMetadata(err).GetString("unknown"))
	require.Nil(t, errorz.Metadata(nil).Get("unknown"))
}

func TestPrefix(t *testing.T) {
	require.Equal(t, "p2 20: p1 10: test error",
		errorz.Errorf("test error",
			errorz.Prefix("p1 %v", 10),
			errorz.Prefix("p2 %v", 20)).Error())
}

func TestCallers(t *testing.T) {
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(errorz.Errorf("test error")))[0], "errorz_test.TestCallers"))
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(fmt.Errorf("test error")))[0], "errorz_test.TestCallers"))
}

func TestSkip(t *testing.T) {
	err := skipErr()
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.TestSkip"))

	err = skipNoSkipErr()
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.noSkipErr"))
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[1], "errorz_test.TestSkip"))
}

func TestSkipAll(t *testing.T) {
	err := skipAllNoSkipErr()
	require.True(t, strings.HasPrefix(errorz.FormatStackTrace(errorz.GetCallers(err))[0], "errorz_test.TestSkipAll"))
}

//go:noinline
func skipErr() error {
	return errorz.Errorf("test error", errorz.Skip())
}

//go:noinline
func noSkipErr() error {
	return errorz.Errorf("test error")
}

//go:noinline
func skipNoSkipErr() error {
	return errorz.Wrap(noSkipErr(), errorz.Skip())
}

//go:noinline
func skipAllNoSkipErr() error {
	return errorz.Wrap(noSkipErr(), errorz.SkipAll())
}
