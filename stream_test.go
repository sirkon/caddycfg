package caddycfg

import (
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStreamer(t *testing.T) {
	c := caddy.NewTestController("http", `
         root {
             a 1234
         }`)
	s := newStream(c)

	require.True(t, s.NextArg())
	require.True(t, s.NextArg())
	require.Equal(t, s.Token().Value, "root")
	s.Confirm()
	require.True(t, s.NextArg())
	require.Equal(t, s.Token().Value, "{")
	s.Confirm()
	require.False(t, s.NextArg())
	require.Equal(t, s.Token().Value, "{")
	require.True(t, s.Next())
	require.True(t, s.Next())
	require.Equal(t, s.Token().Value, "a")
	s.Confirm()
	require.True(t, s.Next())
	require.Equal(t, s.Token().Value, "1234")
	s.Confirm()
	require.False(t, s.NextArg())
	require.True(t, s.Next())
	require.Equal(t, s.Token().Value, "}")
	s.Confirm()
	require.False(t, s.NextArg())
	require.False(t, s.Next())
	require.False(t, s.Next())
	require.False(t, s.NextArg())
}
