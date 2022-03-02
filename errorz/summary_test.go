package errorz_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-errors/errorz"
)

func TestToSummary(t *testing.T) {
	s := errorz.ToSummary(errorz.Errorf("some error", errorz.Prefix("prefix"), errorz.ID("id"), errorz.Status(http.StatusUnauthorized)))
	require.Equal(t, errorz.Status(http.StatusUnauthorized), s.Status)
	require.Equal(t, errorz.ID("id"), s.ID)
	require.Equal(t, "prefix: some error", s.Message)
	require.NotEmpty(t, s.StackTrace)
	require.True(t, strings.HasPrefix(s.StackTrace[0], "errorz_test.TestToSummary"))

	s = errorz.ToSummary(errorz.Errorf("some error"))
	require.Equal(t, errorz.Status(0), s.Status)
	require.Equal(t, errorz.ID(""), s.ID)
	require.Equal(t, "some error", s.Message)
	require.NotEmpty(t, s.StackTrace)
	require.True(t, strings.HasPrefix(s.StackTrace[0], "errorz_test.TestToSummary"))

	s = errorz.ToSummary(fmt.Errorf("some error"))
	require.Equal(t, errorz.Status(0), s.Status)
	require.Equal(t, errorz.ID(""), s.ID)
	require.Equal(t, "some error", s.Message)
	require.NotEmpty(t, s.StackTrace)
	require.True(t, strings.HasPrefix(s.StackTrace[0], "errorz_test.TestToSummary"))
}
