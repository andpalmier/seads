package internal

import "testing"

func TestResolveBingAdURL(t *testing.T) {

	// Mock input data
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"https://www.bing.com/aclick?ld=e8HzM08O4403De41WVRmwKUTVUCUwik1hnDEBnTTC-fZfwnGkVGTQaUWoSnagykADrMpks0uR0_x74jKWOQCGOM0-y0yYZqy5scQnip0o5WEvaeQGVhBoOVbIhKXZxDuimrWJE9QxKQkJ4wUyhi5KFuUU-LEWxkm-Qn7MNUDb_P5rcp5k55fr8kt3L1EUSMgmwzslUww&u=aHR0cHMlM2ElMmYlMmZ3d3cuZ2FsYXh1cy5jaCUyZmRlJTJmczElMmZwcm9kdWN0JTJmYXBwbGUtaXBhZC0yMDIxLTktZ2VuLW51ci13bGFuLTEwMjAtNjQtZ2Itc3BhY2UtZ3JleS10YWJsZXQtMTY2NDQ2ODYlM2ZjYW1wYWlnbmlkJTNkNDQ2MDMzNTcyJTI2YWRncm91cGlkJTNkMTE3MDk4MTA0NTkxMTk1MSUyNmFkaWQlM2QlMjZkZ0NpZGclM2QlMjZnY2xpZCUzZGJlMjg0YTJmMTgwNjFkNGIxYmNlMjIwNmE2ZGNmZDMyJTI2Z2Nsc3JjJTNkM3AuZHMlMjZtc2Nsa2lkJTNkYmUyODRhMmYxODA2MWQ0YjFiY2UyMjA2YTZkY2ZkMzIlMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2RTRUFfREVfQ0hfRFNBX0ZlZWRfU2NvcmVfMTAtNyUyNnV0bV90ZXJtJTNkOSUyNnV0bV9jb250ZW50JTNkRFNBX0ZlZWRfU2NvcmVfOQ&rlid=be284a2f18061d4b1bce2206a6dcfd32", "https://www.galaxus.ch/de/s1/product/apple-ipad-2021-9-gen-nur-wlan-1020-64-gb-space-grey-tablet-16644686?campaignid=446033572&adgroupid=1170981045911951&adid=&dgCidg=&gclid=be284a2f18061d4b1bce2206a6dcfd32&gclsrc=3p.ds&msclkid=be284a2f18061d4b1bce2206a6dcfd32&utm_source=bing&utm_medium=cpc&utm_campaign=SEA_DE_CH_DSA_Feed_Score_10-7&utm_term=9&utm_content=DSA_Feed_Score_9", false},
		{"https://www.bing.com/aclick?ld=e8HzM08O4403De41WVRmwKUTVUCUwik1hnDEBnTTC-fZfwnGkVGTQaUWoSnagykADrMpks0uR0_x74jKWOQCGOM0-y0yYZqy5scQnip0o5WEvaeQGVhBoOVbIhKXZxDuimrWJE9QxKQkJ4wUyhi5KFuUU-LEWxkm-Qn7MNUDb_P5rcp5k55fr8kt3L1EUSMgmwzslUww&u=aHRaaaaacHMlM2ElMmYlMmZ3d3cuZ2FsYXh1cy5jaCUyZmRlJTJmczElMmZwcm9kdWN0JTJmYXBwbGUtaXBhZC0yMDIxLTktZ2VuLW51ci13bGFuLTEwMjAtNjQtZ2Itc3BhY2UtZ3JleS10YWJsZXQtMTY2NDQ2ODYlM2ZjYW1wYWlnbmlkJTNkNDQ2MDMzNTcyJTI2YWRncm91cGlkJTNkMTE3MDk4MTA0NTkxMTk1MSUyNmFkaWQlM2QlMjZkZ0NpZGclM2QlMjZnY2xpZCUzZGJlMjg0YTJmMTgwNjFkNGIxYmNlMjIwNmE2ZGNmZDMyJTI2Z2Nsc3JjJTNkM3AuZHMlMjZtc2Nsa2lkJTNkYmUyODRhMmYxODA2MWQ0YjFiY2UyMjA2YTZkY2ZkMzIlMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2RTRUFfREVfQ0hfRFNBX0ZlZWRfU2NvcmVfMTAtNyUyNnV0bV90ZXJtJTNkOSUyNnV0bV9jb250ZW50JTNkRFNBX0ZlZWRfU2NvcmVfOQ&rlid=be284a2f18061d4b1bce2206a6dcfd32", "", true},
		{"", "", true}, // Adding a case with an empty string
	}

	// Call the function
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ResolveBingAdURL(tt.input)
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
