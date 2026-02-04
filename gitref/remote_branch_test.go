package gitref

import "testing"

func TestParseRemoteBranch(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   RemoteBranch
		wantOK bool
	}{
		{
			name:   "origin style",
			input:  "origin/main",
			want:   RemoteBranch{Remote: "origin", Branch: "main"},
			wantOK: true,
		},
		{
			name:   "full ref",
			input:  "refs/remotes/upstream/release",
			want:   RemoteBranch{Remote: "upstream", Branch: "release"},
			wantOK: true,
		},
		{
			name:   "invalid",
			input:  "main",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseRemoteBranch(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got != tt.want {
				t.Fatalf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
