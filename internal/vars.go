package internal

import (
	"github.com/fatih/color"
	"github.com/go-rod/rod/lib/devices"
	"time"
)

// SearchEngineFunction holds the search engine name and its corresponding function
type SearchEngineFunction struct {
	EngineName     string
	SearchFunction func(string, string, string, bool) ([]AdResult, error)
}

// AdResult contains information regarding an ad found
type AdResult struct {
	Engine           string    `json:"engine"`
	Query            string    `json:"query"`
	OriginalAdURL    string    `json:"OriginalAdURL"`
	FinalDomainURL   string    `json:"final-domain-url"`
	FinalRedirectURL string    `json:"final-redirect-url"`
	RedirectChain    []string  `json:"redirect-chain"`
	Time             time.Time `json:"time"`
	Advertiser       string    `json:"advertiser"`
	Location         string    `json:"location"`
}

var (
	// command line args
	ConfigFilePath      = "config.yaml"
	ConcurrencyLevel    = 4
	ScreenshotPath      = ""
	PrintCleanLinks     = false
	EnableNotifications = false
	PrintRedirectChain  = false
	UserAgentString     = ""
	EnableURLScan       = false
	OutputFilePath      = ""
	NoRedirection       = false
	HtmlPath            = ""
	Logger              = false

	// search engine URLs
	googleurl     = "https://www.google.com/search?q="
	bingurl       = "https://www.bing.com/search?form=QBLH&q="
	duckduckgourl = "https://duckduckgo.com/?ia=web&q="
	yahoourl      = "https://search.yahoo.com/search?q="
	syndicatedurl = "https://syndicatedsearch.goog/afs/ads?adsafe=medium&adtest=off&adpage=1&channel=ch1&client=amg-informationvine&r=m&hl=en&ie=utf-8&adrep=5&oe=utf-8&type=0&format=p5%7Cn5&ad=n5p5&output=uds_ads_only&v=3&bsl=8&pac=0&u_his=5&uio=--&cont=text-ad-block-0%7Ctext-ad-block-1&rurl=https%3A%2F%2Fwww.ask.com%2Fweb%3F%26o%3D0%26an%3Dorganic%26ad%3DOther%2BSEO%26capLimitBypass%3Dfalse%26qo%3DserpSearchTopBox%26q&q="
	adsenseads    = "https://www.adsensecustomsearchads.com/afs/ads?adsafe=medium&adtest=off&adpage=1&channel=ch1&client=amg-informationvine&r=m&hl=en&ie=utf-8&adrep=5&oe=utf-8&type=0&format=p5%7Cn5&ad=n5p5&output=uds_ads_only&v=3&bsl=8&pac=0&u_his=5&uio=--&cont=text-ad-block-0%7Ctext-ad-block-1&rurl=https%3A%2F%2Fwww.ask.com%2Fweb%3F%26o%3D0%26an%3Dorganic%26ad%3DOther%2BSEO%26capLimitBypass%3Dfalse%26qo%3DserpSearchTopBox%26q&q="
	aolurl        = "https://search.aol.com/aol/search?q="

	// useful domains
	googledomain      = "google.com"
	bingdomain        = "bing.com"
	ddgdomain         = "duckduckgo.com"
	syndicateddomain  = "syndicatedsearch.goog"
	adsenseadsdomain  = "adsensecustomsearchads.com"
	doubleclickdomain = "ad.doubleclick.net"
	googleadsservices = "googleadservices.com"
	dadxio            = "d.adx.io"

	searchEngineURLs = map[string]string{
		"google":     googleurl,
		"bing":       bingurl,
		"duckduckgo": duckduckgourl,
		"yahoo":      yahoourl,
		"syndicated": syndicatedurl,
		"adsenseads": adsenseads,
		"aol":        aolurl,
	}

	// some search engines prefer specific User-Agent strings
	ChromeMacUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"

	// ad link selectors
	googleSelector     = "a.sVXRqc"
	bingSelector       = `li.b_adTop [role="link"]`
	ddgSelector        = `li[data-layout="ad"] a[data-testid="result-extras-url-link"]`
	yahooSelector      = `ol.searchCenterTopAds a[data-matarget="ad"]`
	syndicatedSelector = "a.si27"
	adsenseadsSelector = "a.si27"
	aolSelector        = `a[data-matarget="ad"]`
	adinfoSelector     = `a.si149`

	// other utils
	googleCookieBtn = "button#W0wltc"
	yahooCookieBtn  = `button[value="reject"]`
	yahooScrollBtn  = `button#scroll-down-btn`
	aolCookieBtn    = `button[value="reject"]`
	aolScrollBtn    = `button#scroll-down-btn`

	// search engine functions
	searchEnginesFunctions = []SearchEngineFunction{
		{EngineName: "google", SearchFunction: searchGoogleAds},
		{EngineName: "bing", SearchFunction: searchBingAds},
		{EngineName: "duckduckgo", SearchFunction: searchDuckDuckGoAds},
		{EngineName: "yahoo", SearchFunction: searchYahooAds},
		{EngineName: "syndicated", SearchFunction: searchSyndicatedAds},
		{EngineName: "adsenseads", SearchFunction: searchAdsenseAds},
		{EngineName: "aol", SearchFunction: searchAolAds},
	}

	// color variables
	green  = color.New(color.FgGreen)
	italic = color.New(color.Italic)
	red    = color.New(color.FgRed)

	// Laptop device to be used
	Laptop = devices.Device{
		Title:          "laptop",
		Capabilities:   []string{},
		UserAgent:      ChromeMacUA,
		AcceptLanguage: "en",
		Screen: devices.Screen{
			DevicePixelRatio: 1,
			Horizontal: devices.ScreenSize{
				Width:  1200,
				Height: 800,
			},
			Vertical: devices.ScreenSize{
				Width:  1200,
				Height: 800,
			},
		},
	}

	// Ascii art for the banner
	AsciiArt = `
███████╗███████╗ █████╗ ██████╗ ███████╗
██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝
███████╗█████╗  ███████║██║  ██║███████╗
╚════██║██╔══╝  ██╔══██║██║  ██║╚════██║
███████║███████╗██║  ██║██████╔╝███████║
╚══════╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝
Search Engine Ad Scanner - by andpalmier
`
)
