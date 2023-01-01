package generic

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestReduce(t *testing.T) {
	type args struct {
		s    []int
		init int
		f    func(int, int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "001",
			args: args{
				s:    []int{1, 2, 3},
				init: 0,
				f: func(cur int, next int) int {
					return cur + next
				},
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.s, tt.args.init, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterSame(t *testing.T) {
	type args struct {
		s []int
		f func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "001",
			args: args{
				s: []int{1, 2, 3},
				f: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.args.s, tt.args.f)
			assert.Exactly(t, tt.want, got)
			assert.Exactly(t, fmt.Sprintf("%p", tt.args.s), fmt.Sprintf("%p", got))
		})
	}
}

func TestFilterNew(t *testing.T) {
	type args struct {
		s []float32
		f func(float32) bool
	}
	tests := []struct {
		name string
		args args
		want []float32
	}{
		{
			name: "001",
			args: args{
				s: []float32{1, 2, 3},
				f: func(i float32) bool {
					return i/2 == 1
				},
			},
			want: []float32{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterNew(tt.args.s, tt.args.f)
			assert.Exactly(t, tt.want, got)
			assert.NotEqual(t, fmt.Sprintf("%p", tt.args.s), fmt.Sprintf("%p", got))
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		s []int8
		v int8
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "true", args: args{s: []int8{1, 2, 3}, v: 2}, want: true},
		{name: "false", args: args{s: []int8{1, 2, 3}, v: 4}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.s, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	type args struct {
		s []rune
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{name: "001", args: args{[]rune{'x', 'y', 'z'}}, want: []rune{'z', 'y', 'x'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reverse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapSame(t *testing.T) {
	type args struct {
		s []string
		f func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "001",
			args: args{
				s: []string{"a", "b", "c"},
				f: func(s string) string {
					return s + "|"
				},
			},
			want: []string{"a|", "b|", "c|"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.args.s, tt.args.f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
			assert.Exactly(t, fmt.Sprintf("%p", tt.args.s), fmt.Sprintf("%p", got))
		})
	}
}

func TestMapNew(t *testing.T) {
	type args struct {
		s []string
		f func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "001",
			args: args{
				s: []string{"a", "b", "c"},
				f: func(s string) string {
					return s + "|"
				},
			},
			want: []string{"a|", "b|", "c|"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapNew(tt.args.s, tt.args.f)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
			assert.NotEqual(t, fmt.Sprintf("%p", tt.args.s), fmt.Sprintf("%p", got))
		})
	}
}
