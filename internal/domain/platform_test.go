package domain

import "testing"

func TestPlatformValidate(t *testing.T) {
	tests := []struct {
		name     string
		platform Platform
		wantErr  error
	}{
		{
			name:     "valid platform",
			platform: PlatformVinted,
		},
		{
			name:     "invalid platform",
			platform: "HUI",
			wantErr:  ErrInvalidPlatform,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.platform.Validate()

			if err != tt.wantErr {
				t.Fatalf("got error %v want %v", err, tt.wantErr)
			}
		})
	}
}
