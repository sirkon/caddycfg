package caddycfg

import (
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONUnmarshaler(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   **tmpType
		expected string
		wantErr  bool
	}

	var dest *tmpType

	samples := []sample{
		{
			name:     "success",
			input:    "root abcdefgh",
			target:   &dest,
			expected: "abcdefgh",
			wantErr:  false,
		},
		{
			name:     "error-junk-data",
			input:    "root abcdefgh junk",
			target:   &dest,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "error-missing-data",
			input:    "root",
			target:   &dest,
			expected: "",
			wantErr:  true,
		},
	}

	for _, s := range samples {
		t.Run(s.name, func(t *testing.T) {
			c := caddy.NewTestController("http", s.input)
			err := Unmarshal(c, s.target)
			if err != nil {
				if !s.wantErr {
					t.Error(err)
				}
				return
			}
			if err == nil && s.wantErr {
				t.Errorf("error expected")
				return
			}
			require.Equal(t, s.expected, (*s.target).a)
		})
	}
}

func TestBoolean(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *bool
		expected bool
		wantErr  bool
	}

	var target bool
	samples := []sample{
		{
			name:     "true",
			input:    "root true",
			target:   &target,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "false",
			input:    "root false",
			target:   &target,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "error-wrong-data",
			input:    "root error",
			target:   &target,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "error-missing-data",
			input:    "root",
			target:   &target,
			expected: false,
			wantErr:  true,
		},
		{
			name:     "error-junk-data",
			input:    "root true 1234",
			target:   &target,
			expected: false,
			wantErr:  true,
		},
	}

	for _, s := range samples {
		t.Run(s.name, func(t *testing.T) {
			c := caddy.NewTestController("http", s.input)
			err := Unmarshal(c, s.target)
			if err != nil {
				if !s.wantErr {
					t.Error(err)
				}
				return
			}
			if err == nil && s.wantErr {
				t.Errorf("error expected")
				return
			}
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestString(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *string
		expected string
		wantErr  bool
	}

	var target string
	samples := []sample{
		{
			name:     "success",
			input:    "root data",
			target:   &target,
			expected: "data",
			wantErr:  false,
		},
		{
			name:    "error-missing-data",
			input:   "root",
			target:  &target,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root data 1234",
			target:  &target,
			wantErr: true,
		},
	}

	for _, s := range samples {
		t.Run(s.name, func(t *testing.T) {
			c := caddy.NewTestController("http", s.input)
			err := Unmarshal(c, s.target)
			if err != nil {
				if !s.wantErr {
					t.Error(err)
				}
				return
			}
			if err == nil && s.wantErr {
				t.Errorf("error expected")
				return
			}
			require.Equal(t, s.expected, *s.target)
		})
	}
}
