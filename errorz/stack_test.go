package errorz

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPackageFromFuncName(t *testing.T) {
	require.Equal(t, "", getPackageFromFuncName(""))
	require.Equal(t, "a", getPackageFromFuncName("a.b"))
	require.Equal(t, "a", getPackageFromFuncName("a.b.c"))
	require.Equal(t, "a", getPackageFromFuncName("a.(*b).c"))
	require.Equal(t, "a/b/c", getPackageFromFuncName("a/b/c.d"))
	require.Equal(t, "a/b/c", getPackageFromFuncName("a/b/c.(*d).e"))
	require.Equal(t, "a.com/b.c/d", getPackageFromFuncName("a.com/b.c/d.e.f"))
	require.Equal(t, "a.com/b.c/d", getPackageFromFuncName("a.com/b.c/d.(*e).f"))
}
