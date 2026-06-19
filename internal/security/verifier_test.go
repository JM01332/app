package security

import "testing"

func TestMatchesClient(t *testing.T) {
	testCases := []struct {
		name            string
		audience        []string
		authorizedParty string
		want            bool
	}{
		{
			name:     "matching audience",
			audience: []string{"account", "python-client"},
			want:     true,
		},
		{
			name:            "matching authorized party",
			audience:        []string{"account"},
			authorizedParty: "python-client",
			want:            true,
		},
		{
			name:            "different client",
			audience:        []string{"account"},
			authorizedParty: "other-client",
			want:            false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := matchesClient(testCase.audience, testCase.authorizedParty, "python-client")
			if got != testCase.want {
				t.Errorf("matchesClient() = %t, want %t", got, testCase.want)
			}
		})
	}
}
