package util

import (
	"embed"
	"testing"

	"router-practice/internal/router"
)

//go:embed embed_test/*
var fncEMBED embed.FS

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
				path: "../../main.go",
			},
			wantResult: true,
		},
		{
			name: "CheckFileExists_embed",
			args: args{
				path: "embed_test/sample.txt",
			},
			wantResult: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmbed := false
			if tt.name == "CheckFileExists_embed" {
				isEmbed = true
				router.Content = fncEMBED
			}

			if gotResult := CheckFileExists(tt.args.path, isEmbed); gotResult != tt.wantResult {
				t.Errorf("CheckFileExists() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
