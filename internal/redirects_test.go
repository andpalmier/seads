package internal

import "testing"

func TestResolveDoubleClickAdURL(t *testing.T) {
	// Mock input data
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"https://ad.doubleclick.net/searchads/link/click?lid=43700075582814360&ds_s_kwgid=58700005073863721&ds_a_cid=407529252&ds_a_caid=9330499856&ds_a_agid=94333264717&ds_a_lid=kwd-26976553897&&ds_e_adid=72705331872109&ds_e_target_id=kwd-72705993882225:loc-175&&ds_e_network=o&ds_url_v=2&ds_dest_url=https://www.mediamarkt.ch/de/category/_apple-ipad-751034.html?utm_source=bing&utm_medium=cpc&utm_campaign=rt_search_performance_nsp_na_de-2-bk-marke-apple-2023-00-0015&gclid=b9b1a35b09ed1da60350cb72eeead8a6&gclsrc=3p.ds&ds_rl=1254804&msclkid=b9b1a35b09ed1da60350cb72eeead8a6", "https://www.mediamarkt.ch/de/category/_apple-ipad-751034.html?utm_source=bing", false},
		{"https://ad.doubleclick.net/searchads/link/click?lid=43700075582814360&ds_s_kwgid=58700005073863721&ds_a_cid=407529252&ds_a_caid=9330499856&ds_a_agid=94333264717&ds_a_lid=kwd-26976553897&&ds_e_adid=72705331872109&ds_e_target_id=kwd-72705993882225:loc-175&&ds_e_network=o&ds_url_v=2&ds_dsddest_url=https://www.mediamarkt.ch/de/category/_apple-ipad-751034.html?utm_source=bing&utm_medium=cpc&utm_campaign=rt_search_performance_nsp_na_de-2-bk-marke-apple-2023-00-0015&gclid=b9b1a35b09ed1da60350cb72eeead8a6&gclsrc=3p.ds&ds_rl=1254804&msclkid=b9b1a35b09ed1da60350cb72eeead8a6", "", true},
		{"", "", true}, // Adding a case with an empty string
	}

	// Call the function
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ResolveDoubleClickAdURL(tt.input)
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
