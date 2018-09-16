package caddycfg

import (
	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
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

func TestMap(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
		wantErr  bool
	}

	type customType string

	var strint map[string]int
	var intstr map[int]string
	var forbid map[customType]string
	samples := []sample{
		{
			name: "success-string-int",
			input: `
                 root {
                    a 1
                    b 2
                    c 3
                 }
`,
			target: &strint,
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			wantErr: false,
		},
		{
			name:     "success-string-empty-no-data",
			input:    "root { }",
			target:   &strint,
			expected: map[string]int(nil),
			wantErr:  false,
		},
		{
			name: "sucess-int-string",
			input: `
                    root {
                       1 a
                       2 b
                    }
                    `,
			target: &intstr,
			expected: map[int]string{
				1: "a",
				2: "b",
			},
			wantErr: false,
		},
		{
			name:    "error-unclosed-block",
			input:   "root { a 1",
			target:  &strint,
			wantErr: true,
		},
		{
			name:    "error-forbidden-type",
			input:   "root { }",
			target:  &forbid,
			wantErr: true,
		},
		{
			name: "error-duplicate-keys",
			input: `
                 root {
                     a 1
                     a 1
                 }`,
			target:  &strint,
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
			assert.Equal(t, s.expected, reflect.ValueOf(s.target).Elem().Interface())
		})
	}
}

func TestStruct(t *testing.T) {
	type (
		sample struct {
			name     string
			input    string
			target   interface{}
			expected interface{}
			wantErr  bool
		}

		argFriendly struct {
			Args
			A int `json:"a"`
		}

		argUnfriendly struct {
			A int `json:"a"`
		}

		sub struct {
			A int `json:"a"`
		}
		head struct {
			sub
			B struct {
				A string `json:"a"`
			} `json:"b"`
		}
		optional struct {
			A *sub   `json:"a"`
			B string `json:"b"`
		}
	)

	var (
		friendly      argFriendly
		unfriendly    argUnfriendly
		complexStruct head
		option        optional
	)

	samples := []sample{
		{
			name: "success-no-args-friendly",
			input: `
                root {
                   a 23
                }`,
			target: &friendly,
			expected: argFriendly{
				A: 23,
			},
			wantErr: false,
		},
		{
			name: "success-no-args-unfriendly",
			input: `
                root {
                   a 1
                }`,
			target: &unfriendly,
			expected: argUnfriendly{
				A: 1,
			},
			wantErr: false,
		},
		{
			name: "success-args-friendly",
			input: `
                root a b c {
                   a 1
                }`,
			target: &friendly,
			expected: argFriendly{
				Args: Args{
					data: []string{"a", "b", "c"},
				},
				A: 1,
			},
			wantErr: false,
		},
		{
			name: "error-args-unfriendly",
			input: `
                root a b c {
                   a 1
                }`,
			target:   &unfriendly,
			expected: argUnfriendly{},
			wantErr:  true,
		},
		{
			name:    "error-args-friendly-no-block",
			input:   `root a b c`,
			target:  &friendly,
			wantErr: true,
		},
		{
			name: "success-complex-struct",
			input: `
                root {
                    a 12
                    b {
                        a "text lol"
                    }
                }`,
			target: &complexStruct,
			expected: head{
				sub: sub{
					A: 12,
				},
				B: struct {
					A string `json:"a"`
				}{
					A: "text lol",
				},
			},
			wantErr: false,
		},
		{
			name: "success-optional-field-absent",
			input: `
                root {
                    b "1111"
                }`,
			target: &option,
			expected: optional{
				A: nil,
				B: "1111",
			},
			wantErr: false,
		},
		{
			name: "success-optional-field-exist",
			input: `
                root {
                    a {
                        a 12
                    }
                }`,
			target: &option,
			expected: optional{
				A: &sub{
					A: 12,
				},
			},
			wantErr: false,
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
			assert.Equal(t, s.expected, reflect.ValueOf(s.target).Elem().Interface())
		})
	}
}
