package caddycfg

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type tmpType struct {
	a string
}

func (a *tmpType) UnmarshalJSON(data []byte) error {
	a.a = string(data)
	return nil
}

func TestReference(t *testing.T) {
	var dest ***int
	r := ref(reflect.ValueOf(&dest).Elem())
	r.Set(reflect.ValueOf(112))
	require.Equal(t, 112, ***dest)

	var input **tmpType
	r = ref(reflect.ValueOf(&input).Elem())
	r.Interface().(json.Unmarshaler).UnmarshalJSON([]byte("1234"))
	require.Equal(t, "1234", (*input).a)
}

func TestRefType(t *testing.T) {
	type ownTmpType struct {
		json.Unmarshaler
	}

	var dest ***ownTmpType
	whatType, isJSONUnmarshaler := refType(reflect.ValueOf(&dest).Elem().Type())
	require.True(t, isJSONUnmarshaler)
	require.IsType(t, &ownTmpType{}, reflect.Zero(whatType).Interface())
}

func TestLocerrf(t *testing.T) {
	tok := Token{
		File:  "caddyfile",
		Value: "abcdef",
		Lin:   1,
		Col:   1,
	}
	err := locErrf(tok, "error happened: %s", errors.New("error"))
	require.Equal(t, "caddyfile:1: error happened: error", err.Error())
}

func TestCreateStructIndex(t *testing.T) {
	type (
		sub2 struct {
			A int `json:"a"`
		}
		sub1 struct {
			sub2
			B int `json:"b"`
		}
		head struct {
			sub1
			C int `json:"c"`
		}
	)

	tests := []struct {
		name     string
		expected map[string][]int
		s        interface{}
		wantErr  bool
	}{
		{
			name: "direct-success",
			expected: map[string][]int{
				"a": {0},
				"b": {1},
				"c": {2},
			},
			s: struct {
				A int `json:"a"`
				B int `json:"b"`
				C int `json:"c"`
			}{},
			wantErr: false,
		},
		{
			name: "anonymous-success",
			expected: map[string][]int{
				"a": {0, 0, 0},
				"b": {0, 1},
				"c": {1},
			},
			s:       head{},
			wantErr: false,
		},
		{
			name: "error-no-json-tag",
			s: struct {
				A int
			}{},
			wantErr: true,
		},
		{
			name: "error-duplicate-json-tag",
			s: struct {
				sub2
				B int `json:"a"`
			}{},
			wantErr: true,
		},
		{
			name:     "success-empty-index",
			expected: map[string][]int{},
			s: struct {
				data string
			}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := map[string][]int{}
			if err := createStructIndex(index, reflect.ValueOf(tt.s), []int{}); (err != nil) != tt.wantErr {
				t.Errorf("createStructIndex() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				assert.Equal(t, tt.expected, index)
			}
		})
	}
}
