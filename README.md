# seads - Search Engine ADs Scanner

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc Card](https://godoc.org/github.com/andpalmier/seads?status.svg)](https://godoc.org/github.com/andpalmier/seads)
[![Go Report Card](https://goreportcard.com/badge/github.com/andpalmier/seads)](https://goreportcard.com/report/github.com/andpalmier/seads)
[![follow on X](https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x)](https://x.com/intent/follow?screen_name=andpalmier)

`seads` (Search Engine ADs Scanner) is a utility designed to automatically detect advertisements displayed on most popular search engines when searching for a user-submitted keywords.

![demo](https://github.com/andpalmier/seads/blob/main/img/seads.gif?raw=true)

## Why? 🤔

Cybercriminals are increasingly using search engines ads to drive traffic to phishing sites, malware downloads, and other harmful content. `seads` aims to help security researchers and incident response teams identify these ads quickly and efficiently.

### Features ⚡️

- **Multiple search engines support**: Currently supports Google*, Bing, DuckDuckGo, Yahoo, Aol, Syndicated* and AdSense*.
- **Automated reporting**: Send reports of findings via email, Slack, or Telegram.
- **Concurrent search**: Specify multiple headless instances to gather as many ads as possible concurrently.
- **Screenshots**: Capture screenshots of ads found in search engines for evidence.
- **Docker support**: Run `seads` using Docker.
- **Export in JSON**: Export results of the execution in JSON format.
- **Custom User-Agent**: Provide your User-Agent string to be used to click on ads found.
- **Redirect chain detection**: Tracks URLs through redirects to detect and log chains. 
- **URLScan submission**: Submit the link to [URLScan](https://urlscan.io) using your API key.
- **Advertiser name detection**: Print the advertiser name and location (only for Google, Syndicated and AdSense).
- **Defanging**: Defang URLs in notifications to prevent accidental clicks.
- **HTML page saving**: Save the HTML page of the search engine result for later analysis.
- **Expected domains**: Specify expected domains to filter out known advertisers from notifications.

\* **NB**: Currently, Google Search detection doesn't always work, and the automated browser is often prompted by a CAPTCHA. As a workaround, Syndicated and AdSense are used to gather ads from Google ([see here](https://support.google.com/adsense/answer/14201307)). This may not be 100% accurate, but it is the best available option as of now.

### Known limitations 😕
- Due to the nature of search engine ads, a single search may not reveal all ads. Using concurrent headless browsers might slow down detection but ensures comprehensive ad gathering.
- Notifications on Slack and Telegram have character limits. Messages exceeding the limit won't be sent.

## Getting started 🛠️

> [!CAUTION]  
> `seads` will click on dangerous links, and potentially expose your IP address to malicious sites. Consider using a VPN/proxy or running `seads` in a VPS to avoid this risk.

### Install

You can download `seads` from the [releases section](https://github.com/andpalmier/seads/releases).
Or using `go`:

```bash
go install github.com/andpalmier/seads@latest
seads -h
```

For Docker, you need to clone the GitHub repo and run:

```bash
docker build -t seads .
docker run -v "$(pwd)":/mnt seads -h
```

### Before running

Create a `config.yaml` file with the following structure (be sure to check out the full example in the repo!):

```yaml
urlscan:
  token: "APIKEY"
  scanurl: "https://urlscan.io/api/v1/scan/"
  visibility: "unlisted"
  tags: "seads_ads_tracker"

global-domain-exclusion:
  exclusion-list: [ebay.com, amazon.com]
  
queries:
  - query: "ipad"
    expected-domains: [apple.com]

  - query: "as roma"
    expected-domains: []
```

This config will search for ads related to `ipad` and `as roma`.
The field `expected-domains` is used to specify domains we are expecting to appear in the ads of search engines while searching for the specified keywords.
Domains in `expected-domains` and in `global-domain-exclusion` will still appear in the output of `seads`, but won’t be sent in the notification.
The `urlscan` section is used to specify the API key and the URLScan API endpoint. The `visibility` field can be set to `public`, `unlisted`, or `private`. The `tags` field is used to specify tags for the scan.

## Examples 📖

Running `seads` with the provided config, storing the results in `results.json`, and sending notifications via the channels configured in the config:

```bash
seads -config config.yaml -out results.json -notify
```

Same as above, but in Docker:

```bash
docker run -v "$(pwd)":/mnt seads -config /mnt/config.yaml -out /mnt/results.json -notify
```

Running with redirection chain handled by URLScan, and save HTML page and screenshot:

```bash
seads -config config.yaml -urlscan -noredirect -screenshot ./screenshots -html ./htmls
```

Same as above, but in Docker:

```bash
docker run -v "$(pwd)":/mnt seads -config /mnt/config.yaml -urlscan -noredirect -screenshot /mnt/screenshots -html /mnt/htmls
```

You can leverage the notification feature by automating the execution of `seads` using a cron job or a task scheduler.
For example, in a Linux machine you can set up a cron job to run `seads` every day at 9 AM by adding the following line to your crontab:

```bash
0 9 * * * /path/to/seads -config /path/to/config.yaml -screenshot /path/to/screenshots -notify
```

Be sure to update the command to reflect the correct paths for the `seads` binary and the configuration file.
If using Docker, adjust the command accordingly for Docker execution.
On Windows and macOS, you can achieve similar scheduling using Task Scheduler or `launchd`.

Screenshot example:
![screenshot](https://github.com/andpalmier/seads/blob/main/img/example-bing-ipad.png?raw=true)

### Available flags

```
  -config (string) [REQUIRED]
    path to config file (default "config.yaml")
  -concurrency (int)
    number of concurrent headless browsers (default 4)
  -cleanlinks
    print clear links in output (links will remain defanged in notifications)
  -html (string)
    path to store search engine result html page (if empty, the htmlPath feature will be disabled)
  -log
    enable detailed logging, VERY VERBOSE!
  -noredirect
    do not follow redirection; if "urlscan" is enabled, submit link to resolve by URLScan instead
  -notify
    notify if unexpected domains are found (requires notifications fields in config.yaml)
  -out (string)
    path of JSON file containing links of gathered ads
  -screenshot (string)
    path to store screenshots (if empty, the screenshot feature will be disabled)
  -printredirect
    print redirection chain for ad links found
  -ua (string)
    User-Agent string to be used to click on ads
  -urlscan
    submit url to urlscan.io for analysis
```

## 3rd party libraries 📚

- Rod: [GitHub repo](https://github.com/go-rod/rod), [documentation](https://go-rod.github.io/)
- Shoutrrr: [GitHub repo](https://github.com/containrrr/shoutrrr), [documentation](https://containrrr.dev/shoutrrr/v0.8/)
- Fatih/color: [GitHub repo](https://github.com/fatih/color), [Go reference](https://pkg.go.dev/github.com/fatih/color)
- Carlmjohnson/requests [GitHub repo](https://github.com/carlmjohnson/requests), [Go reference](https://pkg.go.dev/github.com/carlmjohnson/requests)

## Thank you 🙏🏻

I'd like to thank [@raufridzuan](https://github.com/raufridzuan) for his help and contribution to this project!
