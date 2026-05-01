package generator

import "testing"

func TestDomainRegexToURLRegex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOut string
	}{
		{
			name:    "anchored",
			input:   "^ads\\..*\\.example\\.org$",
			wantOut: "^https?://ads\\..*\\.example\\.org(?::\\d+)?(?:/|$)",
		},
		{
			name:    "not anchored",
			input:   "google",
			wantOut: "^https?://[^/:]*google[^/:]*(?::\\d+)?(?:/|$)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := domainRegexToURLRegex(test.input)
			if got != test.wantOut {
				t.Fatalf("domainRegexToURLRegex(%q) = %q, want %q", test.input, got, test.wantOut)
			}
		})
	}
}
