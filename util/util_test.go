package util

import "testing"

func TestCheckFileExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "CheckFileExists",
			args: args{
				path: "./hahahaha.hahahaha",
			},
			wantResult: false,
		},
		{
			name: "CheckFileExists",
			args: args{
				path: "../main.go",
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := CheckFileExists(tt.args.path); gotResult != tt.wantResult {
				t.Errorf("CheckFileExists() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
