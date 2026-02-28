package repospec

import (
	"testing"
)

func TestValidationErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		wantCode int
	}{
		{
			name:     "ErrInvalidHost",
			err:      &ValidationError{Code: ErrInvalidHost, Field: "host", Message: "host is required"},
			wantCode: 1001,
		},
		{
			name:     "ErrInvalidGroup",
			err:      &ValidationError{Code: ErrInvalidGroup, Field: "group", Message: "group is required"},
			wantCode: 1002,
		},
		{
			name:     "ErrInvalidRepo",
			err:      &ValidationError{Code: ErrInvalidRepo, Field: "repo", Message: "repo is required"},
			wantCode: 1003,
		},
		{
			name:     "ErrInvalidPort",
			err:      &ValidationError{Code: ErrInvalidPort, Field: "port", Message: "port must be numeric"},
			wantCode: 1004,
		},
		{
			name:     "ErrInvalidBasePath",
			err:      &ValidationError{Code: ErrInvalidBasePath, Field: "basePath", Message: "must start with /"},
			wantCode: 1005,
		},
		{
			name:     "ErrInvalidSubGroup",
			err:      &ValidationError{Code: ErrInvalidSubGroup, Field: "subGroups", Message: "not supported"},
			wantCode: 1006,
		},
		{
			name:     "ErrInvalidProvider",
			err:      &ValidationError{Code: ErrInvalidProvider, Field: "provider", Message: "invalid"},
			wantCode: 1007,
		},
		{
			name:     "ErrInvalidRepoKey",
			err:      &ValidationError{Code: ErrInvalidRepoKey, Field: "repoKey", Message: "invalid format"},
			wantCode: 1008,
		},
		{
			name:     "ErrInvalidURL",
			err:      &ValidationError{Code: ErrInvalidURL, Field: "url", Message: "invalid"},
			wantCode: 1009,
		},
		{
			name:     "ErrInvalidSpec",
			err:      &ValidationError{Code: ErrInvalidSpec, Field: "spec", Message: "invalid"},
			wantCode: 1010,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("got code %d, want %d", tt.err.Code, tt.wantCode)
			}
		})
	}
}

func TestValidationErrorError(t *testing.T) {
	tests := []struct {
		name    string
		err     ValidationError
		wantMsg string
	}{
		{
			name:    "with field",
			err:     ValidationError{Field: "host", Message: "host is required"},
			wantMsg: "host: host is required",
		},
		{
			name:    "without field",
			err:     ValidationError{Message: "some error"},
			wantMsg: "some error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("got %q, want %q", got, tt.wantMsg)
			}
		})
	}
}
