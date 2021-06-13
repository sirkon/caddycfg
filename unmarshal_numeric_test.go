package caddycfg

import (
	"github.com/caddyserver/caddy"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt8(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *int8
		expected int8
		wantErr  bool
	}

	var dest int8

	samples := []sample{
		{
			name:     "success",
			input:    "root 12",
			target:   &dest,
			expected: 12,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 1234",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestInt16(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *int16
		expected int16
		wantErr  bool
	}

	var dest int16

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233",
			target:   &dest,
			expected: 1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 123433",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestInt32(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *int32
		expected int32
		wantErr  bool
	}

	var dest int32

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233",
			target:   &dest,
			expected: 1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestInt64(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *int64
		expected int64
		wantErr  bool
	}

	var dest int64

	samples := []sample{
		{
			name:     "success",
			input:    "root -1233",
			target:   &dest,
			expected: -1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestInt(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *int
		expected int
		wantErr  bool
	}

	var dest int

	samples := []sample{
		{
			name:     "success",
			input:    "root -1233",
			target:   &dest,
			expected: -1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestUint8(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *uint8
		expected uint8
		wantErr  bool
	}

	var dest uint8

	samples := []sample{
		{
			name:     "success",
			input:    "root 12",
			target:   &dest,
			expected: 12,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 1234",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestUint16(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *uint16
		expected uint16
		wantErr  bool
	}

	var dest uint16

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233",
			target:   &dest,
			expected: 1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 123433",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestUint32(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *uint32
		expected uint32
		wantErr  bool
	}

	var dest uint32

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233",
			target:   &dest,
			expected: 1233,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestUint64(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *uint64
		expected uint64
		wantErr  bool
	}

	var dest uint64

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233333",
			target:   &dest,
			expected: 1233333,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestUint(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *uint
		expected uint
		wantErr  bool
	}

	var dest uint

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233333",
			target:   &dest,
			expected: 1233333,
			wantErr:  false,
		},
		{
			name:    "error-overflow",
			input:   "root 12343300000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			target:  &dest,
			wantErr: true,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestFloat32(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *float32
		expected float32
		wantErr  bool
	}

	var dest float32

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233",
			target:   &dest,
			expected: 1233,
			wantErr:  false,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}

func TestFloat64(t *testing.T) {
	type sample struct {
		name     string
		input    string
		target   *float64
		expected float64
		wantErr  bool
	}

	var dest float64

	samples := []sample{
		{
			name:     "success",
			input:    "root 1233333",
			target:   &dest,
			expected: 1233333,
			wantErr:  false,
		},
		{
			name:    "error-junk-data",
			input:   "root 12 junk",
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
			require.Equal(t, s.expected, *s.target)
		})
	}
}
