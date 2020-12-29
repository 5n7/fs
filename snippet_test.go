package fs

import (
	"os/user"
	"path/filepath"
	"testing"
)

func Test_expandTilde(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("test failed: %v", err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				path: "~/foo/bar/baz",
			},
			want:    filepath.Join(u.HomeDir, "foo/bar/baz"),
			wantErr: false,
		},
		{
			name: "",
			args: args{
				path: "~",
			},
			want:    u.HomeDir,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTilde(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandTilde() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("expandTilde() = %v, want %v", got, tt.want)
			}
		})
	}
}
