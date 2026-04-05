//go:build unit

package tests

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func MustClose(t testing.TB, closer io.Closer) {
	require.NoError(t, closer.Close())
}
