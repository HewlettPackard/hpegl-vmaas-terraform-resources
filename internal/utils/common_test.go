// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import "testing"

type testStruct struct {
	val int
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		n    interface{}
		want bool
	}{
		{
			name: "Test case 1: int 0",
			n:    0,
			want: true,
		},
		{
			name: "Test case 2: int 1",
			n:    1,
			want: false,
		},
		{
			name: "Test case 3: string ''",
			n:    "",
			want: true,
		},
		{
			name: "Test case 4: string 'abc'",
			n:    "abc",
			want: false,
		},
		{
			name: "Test case 5: struct {}",
			n:    testStruct{},
			want: true,
		},
		{
			name: "Test case 6: struct {1}",
			n:    testStruct{1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.n); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
