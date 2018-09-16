package caddycfg

import (
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/require"
	"reflect"
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

func TestSlices(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
		wantErr  bool
	}

	var stringSlice []string
	var intSlice []int
	samples := []sample{
		{
			name:     "success-simple-string-slice",
			input:    "root data1 data2",
			target:   &stringSlice,
			expected: []string{"data1", "data2"},
			wantErr:  false,
		},
		{
			name:     "success-missing-data",
			input:    "root",
			target:   &stringSlice,
			expected: []string(nil),
			wantErr:  false,
		},
		{
			name:     "sucess-simple-int-slice",
			input:    "root 1234 4321",
			target:   &intSlice,
			expected: []int{1234, 4321},
			wantErr:  false,
		},
		{
			name:    "error-simple-int-slice-not-a-number",
			input:   "root 1234 not-a-number",
			target:  &intSlice,
			wantErr: true,
		},
		{
			name: "success-complex-string-slice",
			input: `
root { a
  b
  c
  d
}
`,
			target:   &stringSlice,
			expected: []string{"a", "b", "c", "d"},
			wantErr:  false,
		},
		{
			name:    "error-complex-unclosed-block",
			input:   "root {",
			target:  &stringSlice,
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
			require.Equal(t, s.expected, reflect.ValueOf(s.target).Elem().Interface())
		})
	}
}
