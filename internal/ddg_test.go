package internal

import "testing"

func TestResolveDuckDuckGoAdURL(t *testing.T) {
	// Mock input data

	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"https://duckduckgo.com/y.js?ad_domain=amazon.de&ad_provider=bingv7aa&ad_type=txad&click_metadata=82egrwr67c_dPFSLF7ujxyJ6t5rVdL2d8znAksMjbwy9992dNhJTTCm8lZz4ziAi0Iuh0cWxuzVptjG9ysIqsChKEoJcJWCkMJWQvtcXbdlaRvtvZsQb1lQcHLg8xCGuaONlyYZG5BaJvQoXneYjPQ.1HMiwHE7jMz%2DSWAQ_ZC19Q&eddgt=eQDAWPYPcgW_s8f5RZ5%2Drg%3D%3D&rut=df718eaf133e92b365529a1861a2f41cfcf174f46df95f8ece429a817d7a6acf&u3=https%3A%2F%2Fwww.bing.com%2Faclick%3Fld%3De8%2D7JNgCaN__Q3DK5letIKHDVUCUxNeUeMitBN9iBGBaDRBrO3N4U%2D_G5AqWJng7soJo087qdr8pOHmnHb8YrHjluXq3%2Dv4R_djJYHb314e75DnyWbchF9NvEfXijv3ds86lHLWap3qfsC9g8bL8mQ3%2DtDL454u8Mq8_LHbPE6WwOi3uljXudXI2Ps3ETrCIaURK1Zdw%26u%3DaHR0cHMlM2ElMmYlMmZ3d3cuYW1hem9uLmRlJTJmcyUyZiUzZmllJTNkVVRGOCUyNmtleXdvcmRzJTNkaXBhZCUyNmluZGV4JTNkYXBzJTI2dGFnJTNkaHlkZGVtc24tMjElMjZyZWYlM2RwZF9zbF84NzZieTcyaDBzX2UlMjZhZGdycGlkJTNkMTE4OTY3MjIzNjM5MzE1MiUyNmh2YWRpZCUzZDc0MzU0NjMxMzUzNjE0JTI2aHZuZXR3JTNkbyUyNmh2cW10JTNkZSUyNmh2Ym10JTNkYmUlMjZodmRldiUzZGMlMjZodmxvY2ludCUzZCUyNmh2bG9jcGh5JTNkMjE3ODg1JTI2aHZ0YXJnaWQlM2Rrd2QtNzQzNTQ3MjA3NjA3MDAlM2Fsb2MtMTc1JTI2aHlkYWRjciUzZDI5MjI1XzIzNjgyNzAlMjZtY2lkJTNkNjk0YTViNDA4YjhlMzBjODg2ZTYzYzkwYWIwY2YyYTUlMjZtc2Nsa2lkJTNkYjAxYzUwMDcxY2I1MTFlMzBjZTkzODJlYTI2OGIxODY%26rlid%3Db01c50071cb511e30ce9382ea268b186&vqd=4-161279935822117090085869941001314414437&iurl=%7B1%7DIG%3D04AD3191FC424E439C64638B02AFE9EB%26CID%3D219157009B056A13365542F09A526B4C%26ID%3DDevEx%2C5048.1", "https://www.amazon.de/s/?ie=UTF8&keywords=ipad&index=aps&tag=hyddemsn-21&ref=pd_sl_876by72h0s_e&adgrpid=1189672236393152&hvadid=74354631353614&hvnetw=o&hvqmt=e&hvbmt=be&hvdev=c&hvlocint=&hvlocphy=217885&hvtargid=kwd-74354720760700:loc-175&hydadcr=29225_2368270&mcid=694a5b408b8e30c886e63c90ab0cf2a5&msclkid=b01c50071cb511e30ce9382ea268b186", false},
		{"https://duckduckgo.com/y.js?ad_domain=amazon.de&ad_provider=bingv7aa&ad_type=txad&click_metadata=82egrwr67c_dPFSLF7ujxyJ6t5rVdL2d8znAksMjbwy9992dNhJTTCm8lZz4ziAi0Iuh0cWxuzVptjG9ysIqsChKEoJcJWCkMJWQvtcXbdlaRvtvZsQb1lQcHLg8xCGuaONlyYZG5BaJvQoXneYjPQ.1HMiwHE7jMz%2DSWAQ_ZC19Q&eddgt=eQDAWPYPcgW_s8f5RZ5%2Drg%3D%3D&rut=df718eaf133e92b365529a1861a2f41cfcf174f46df95f8ece429a817d7a6acf&u3=https%3A%2F%2Fwww.bing.com%2Faclick%sdfsf3Fld%3De8%2D7JNgCaN__Q3DK5letIasdKHDVUCUxNeUeMitBN9iBGBaDRBrO3N4U%2D_G5AqWJng7soJo087qdr8pOHmnHb8YrHjluXq3%2Dv4R_djJYHb314e75DnyWbchF9NvEfXijv3ds86lHLWap3qfsC9g8bL8mQ3%2DtDL454u8Mq8_LHbPE6WwOi3uljXudXI2Ps3ETrCIaURK1Zdw%26u%3DaHR0cHMlM2ElMmYlMmZ3d3cuYW1hem9uLmRlJTJmcyUyZiUzZmllJTNkVVRGOCUyNmtleXdvcmRzJTNkaXBhZCUyNmluZGV4JTNkYXBzJTI2dGFnJTNkaHlkZGVtc24tMjElMjZyZWYlM2RwZF9zbF84NzZieTcyaDBzX2UlMjZhZGdycGlkJTNkMTE4OTY3MjIzNjM5MzE1MiUyNmh2YWRpZCUzZDc0MzU0NjMxMzUzNjE0JTI2aHZuZXR3JTNkbyUyNmh2cW10JTNkZSUyNmh2Ym10JTNkYmUlMjZodmRldiUzZGMlMjZodmxvY2ludCUzZCUyNmh2bG9jcGh5JTNkMjE3ODg1JTI2aHZ0YXJnaWQlM2Rrd2QtNzQzNTQ3MjA3NjA3MDAlM2Fsb2MtMTc1JTI2aHlkYWRjciUzZDI5MjI1XzIzNjgyNzAlMjZtY2lkJTNkNjk0YTViNDA4YjhlMzBjODg2ZTYzYzkwYWIwY2YyYTUlMjZtc2Nsa2lkJTNkYjAxYzUwMDcxY2I1MTFlMzBjZTkzODJlYTI2OGIxODY%26rlid%3Db01c50071cb511e30ce9382ea268b186&vqd=4-161279935822117090085869941001314414437&iurl=%7B1%7DIG%3D04AD3191FC424E439C64638B02AFE9EB%26CID%3D219157009B056A13365542F09A526B4C%26ID%3DDevEx%2C5048.1", "", true},
		{"", "", true}, // Adding a case with an empty string
	}

	// Call the function
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ResolveDuckDuckGoAdURL(tt.input)
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
