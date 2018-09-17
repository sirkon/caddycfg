package caddycfg

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArgs_appendData(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name     string
		items    []string
		expected Args
	}{
		{
			name:  "empty-nil",
			items: nil,
			expected: Args{
				data: nil,
			},
		},
		{
			name:  "empty-empty",
			items: []string{},
			expected: Args{
				data: nil,
			},
		},
		{
			name:  "not-empty",
			items: []string{"1", "2"},
			expected: Args{
				data: []string{"1", "2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Args{}
			a.appendData(tt.items)
			require.Equal(t, tt.expected, a)
			if len(tt.items) > 0 {
				require.Equal(t, tt.items, a.Arguments())
			}
		})
	}
}
