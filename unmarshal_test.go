package caddycfg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mholt/caddy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type jsonUnmarshalerNeedEncoding struct {
	a string
}

func (v *jsonUnmarshalerNeedEncoding) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &v.a); err != nil {
		return err
	}
	return nil
}

func TestJSONUnmarshaler(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   interface{}
		expected interface{}
		wantErr  bool
	}

	var dest *tmpType
	var dest2 *jsonUnmarshalerNeedEncoding

	samples := []sample{
		{
			name:   "success",
			input:  "root abcdefgh",
			target: &dest,
			expected: tmpType{
				a: "abcdefgh",
			},
			wantErr: false,
		},
		{
			name:   "success-need-encoding",
			input:  "root abcdefgh",
			target: &dest2,
			expected: jsonUnmarshalerNeedEncoding{
				a: "abcdefgh",
			},
			wantErr: false,
		},
		{
			name:    "error-junk-data",
			input:   "root abcdefgh junk",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-missing-data",
			input:   "root",
			target:  &dest,
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
			// reduce expected and data to values
			require.Equal(t, reduceToValue(s.expected), reduceToValue(s.target))
		})
	}
}

func reduceToValue(e interface{}) interface{} {
	v := reflect.ValueOf(e)
	for v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface()
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
			name:    "error-args-block-mix",
			input:   "root a b c d { e }",
			target:  &stringSlice,
			wantErr: true,
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
			name:    "error-no-data",
			input:   "root",
			target:  &strint,
			wantErr: true,
		},
		{
			name:    "error-args-not-allowed",
			input:   "root a { a 1 }",
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

type argumentsImpl struct {
	data []string
}

func (a *argumentsImpl) AppendArgument(arg Token) error {
	a.data = append(a.data, arg.Value)
	return nil
}

func (a *argumentsImpl) Arguments() []string {
	return a.data
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
		customArgFriendly struct {
			argumentsImpl
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
		jsonCase struct {
			A tmpType `json:"a"`
		}
		jsonDirectCase struct {
			A *tmpType `json:"a"`
		}
		jsonEncodingCase struct {
			A jsonUnmarshalerNeedEncoding `json:"a"`
		}
	)

	var (
		friendly        argFriendly
		customFriendly  customArgFriendly
		nCustomFriendly customArgFriendly
		unfriendly      argUnfriendly
		complexStruct   head
		option          optional
		jcase           jsonCase
		jdircase        jsonDirectCase
		jenccase        jsonEncodingCase
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
			name: "success-custom-no-args-friendly",
			input: ` root {
                         a 23
                     }`,
			target: &customFriendly,
			expected: customArgFriendly{
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
			name: "success-custom-args-friendly",
			input: `
                root a b c {
                   a 777
                }`,
			target: &customFriendly,
			expected: customArgFriendly{
				argumentsImpl: argumentsImpl{
					data: []string{"a", "b", "c"},
				},
				A: 777,
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
			name:   "args-custom-friendly-no-block",
			input:  `root a b c`,
			target: &nCustomFriendly,
			expected: customArgFriendly{
				argumentsImpl: argumentsImpl{
					data: []string{"a", "b", "c"},
				},
			},
			wantErr: false,
		},
		{
			name: "error-unknown-field",
			input: `
                root {
                    field 1
                }`,
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
		{
			// this case covers a sitation when the type T itself doesn't implement JSONUnmarshaler and *T does
			// it is expected the unmarshaler should switch to type pointer in this case
			name: "success-json-unmarshaler-field",
			input: `
                root {
                    a abcdef
                }`,
			target: &jcase,
			expected: jsonCase{
				A: tmpType{
					a: "abcdef",
				},
			},
			wantErr: false,
		},
		{
			name: "success-json-unmarshaler-direct-field",
			input: `
                root {
                    a abcdef
                }`,
			target: &jdircase,
			expected: jsonDirectCase{
				A: &tmpType{
					a: "abcdef",
				},
			},
			wantErr: false,
		},
		{
			name: "success-json-unmarshaler-field-with-encoding",
			input: `
                root {
                    a abcdef
                }`,
			target: &jenccase,
			expected: jsonEncodingCase{
				A: jsonUnmarshalerNeedEncoding{
					a: "abcdef",
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

func TestUnmarshalForbiddenType(t *testing.T) {
	type sample struct {
		name   string
		target interface{}
	}

	var (
		channel = make(chan int)
		str     string
	)
	samples := []sample{
		{
			name:   "channel",
			target: &channel,
		},
		{
			name:   "string",
			target: str,
		},
	}

	for _, s := range samples {
		t.Run(s.name, func(t *testing.T) {
			c := caddy.NewTestController("http", `root a b c`)
			require.Error(t, Unmarshal(c, s.target))
		})
	}
}

func TestUnmarshalHeadInfo(t *testing.T) {
	c := caddy.NewTestController("http", "root a")
	var dest string
	head, err := UnmarshalHeadInfo(c, &dest)
	require.NoError(t, err)
	require.Equal(t, head.Value, "root")
	require.Equal(t, dest, "a")
}

type argAcc struct {
	data []string
}

func (argAcc) Err(head Token) error {
	return fmt.Errorf(head.Value)
}

func (argAcc) AppendArgument(arg Token) error {
	return nil
}

func (argAcc) Arguments() []string {
	return nil
}

func TestValidation(t *testing.T) {
	dest := argAcc{}
	c := caddy.NewTestController("http", "root a b c")
	err := Unmarshal(c, &dest)
	require.Equal(t, err.Error(), "root")
}
