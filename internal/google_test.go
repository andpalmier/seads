package internal

import "testing"

func TestResolveGoogleAdURL(t *testing.T) {
	// Mock input data
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"https://syndicatedsearch.goog/aclk?sa=L&ai=DChcSEwj3kfDn76yNAxXdpYMHHYbOGwoYABACGgJlZg&co=1&ase=2&gclid=EAIaIQobChMI95Hw5--sjQMV3aWDBx2GzhsKEAAYASAAEgKWU_D_BwE&category=acrcp_v1_0&sig=AOD64_0604Bi8LMNrgZPqKiZI8hLFgvHfQ&adurl=https://www.refurbed.ch/c/ipads/%3Futm_campaign%3DCH%2520-%2520PMax%2520-%2520Near%2520Index%252C%2520Index%252C%2520Over%2520Index%26utm_medium%3Dcpc%26utm_source%3Dgoogle%26gad_source%3D5%26gad_campaignid%3D21651695024&q=", "https://www.refurbed.ch/c/ipads/?utm_campaign=CH%20-%20PMax%20-%20Near%20Index%2C%20Index%2C%20Over%20Index&utm_medium=cpc&utm_source=google&gad_source=5&gad_campaignid=21651695024", false},
		{"https://syndicatedsearch.goog/aclk?sa=L&ai=DChcSEwj3kfDn76yNAxXdpYMHHYbOGwoYABACGgJlZg&co=1&ase=2&gclid=EAIaIQobChMI95Hw5--sjQMV3aWDBx2GzhsKEAAYASAAEgKWU_D_BwE&category=acrcp_v1_0&sig=AOD64_0604Bi8LMNrgZPqKiZI8hLFgvHfQ&aasdfdusrl=https://www.refurbasfdafsded.ch/c/ipads/%sfd%asdfasd3asdfDCH%2520-%2520PMax%2520-%2520Near%2520Index%252C%2520Index%252C%2520Over%2520Index%26utm_medium%3Dcpc%26utm_source%3Dgoogle%26gad_source%3D5%26gad_campaignid%3D21651695024&q=", "", true},
		{"", "", true}, // Adding a case with an empty string
	}

	// Call the function
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ResolveGoogleAdURL(tt.input)
			// Assertions
			if (err != nil) != tt.hasError {
				t.Errorf("expected error: %v\ngot: %v", tt.hasError, err)
			}
			if result != tt.expected {
				t.Errorf("expected %v\ngot %v", tt.expected, result)
			}
		})
	}
}
