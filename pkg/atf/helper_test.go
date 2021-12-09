package atf

import "testing"

func Test_path(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want string
	}{
		{
			name: "2 value",
			args: []interface{}{"acc", 0},
			want: "acc.0",
		},
		{
			name: "1 value",
			args: []interface{}{"acc"},
			want: "acc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := path(tt.args...); got != tt.want {
				t.Errorf("path() = %v, want %v", got, tt.want)
			}
		})
	}
}
