# seads - Search Engine ADs Scanner

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc Card](https://godoc.org/github.com/andpalmier/seads?status.svg)](https://godoc.org/github.com/andpalmier/seads)
[![Go Report Card](https://goreportcard.com/badge/github.com/andpalmier/apkingo)](https://goreportcard.com/report/github.com/andpalmier/seads)
[![follow on X](https://img.shields.io/twitter/follow/andpalmier?style=social&logo=x)](https://x.com/intent/follow?screen_name=andpalmier)

`seads` (Search Engine ADs Scanner) is a utility designed to automatically detect advertisements displayed on most popular search engines when searching for a user-submitted keywords.

For a comprehensive guide on how to use `seads`, please refer to [this blog post](https://andpalmier.com/posts/seads/).

![seads](https://github.com/andpalmier/seads/blob/main/img/seads.gif?raw=true)

## Features:
- **Automated reporting**: Easily send reports of findings via email, Slack, or Telegram.
- **Concurrent Search**: Specify multiple headless instances to gather as many ads as possible concurrently.
- **Screenshot Support**: Capture screenshots of ads found in search engines for evidence.
- **Docker Support**: Install seads without affecting your local setup using Docker.

## Known limitations:
- Due to the nature of search engine ads, a single search may not reveal all ads. Using concurrent headless browsers might slow down detection but ensures comprehensive ad gathering.
- Notifications via Slack and Telegram have character limits. Messages exceeding the limit won't be sent.

## Installation

### Download binary

You can download `seads` from the [releases section](https://github.com/andpalmier/seads/releases).

### Using Go install

You can compile it from source by running:

```
go install github.com/andpalmier/seads/cmd/seads@latest
```

### Using Docker:

Clone the GitHub repo and run Docker:

```
docker build -t seads .
docker run -v "$(pwd)":/mnt seads -h
```

## Usage

You can run `seads` with the following flags:

```
  -config string (REQUIRED)
    	path to config file (default "config.yaml").
  -concurrency int
    	number of concurrent headless browsers (default 4).
  -cleanlinks
    	print clear links in output (links will remain defanged in notifications).
  -notify
    	notify if unexpected domains are found.
  -screenshot string
    	path to store screenshots (if empty, the screenshot feature will be disabled).
```

Example:

```
seads -config config.yaml -notify
```

Docker example:

```
docker run -it -v "$(pwd)":/mnt seads -config /mnt/config.yaml -notify
```

## How to use

After installing `seads`, create a `config.yaml` file with the following structure:

```yaml
mail:
  host: MAILHOST
  port: 587
  username: USERNAME
  password: PASSWORD
  from: FROMADDRESS
  recipients: [RECIPIENTADDRESS#1,RECIPIENTADDRESS#2]

slack:
  token: SLACKTOKEN
  channels: [CHANNEL#1,CHANNEL#2]

telegram:
  token: TELEGRAMTOKEN
  chatid: [CHATID#1,CHATID#2]

queries:
  - query: "ipad"
    expected-domains: [apple.com, amazon.com]

  - query: "as roma"
    expected-domains: []
```

The field `expected-domains` is used to specify domains we are expecting to appear in the ads of search engines while searching for the specified keywords.
Domains in `expected-domains` will still appear in the output of `seads`, but won’t be sent in the notification.

Run `seads` with the following command:

```
seads -config config.yaml -screenshot scr -notify
```

output example:

![seads](https://github.com/andpalmier/seads/blob/main/img/seads.gif?raw=true)

screenshots example:

![seads_yahoo_apple](https://github.com/andpalmier/seads/blob/main/img/example-yahoo-apple.png?raw=true)

notification example:

```
Here are the "unexpected domains" found during the last execution of seads:

Message creation date: 2024-03-12 22:28:14

* Search engine: Yahoo
 Search term: apple
 Domain: reparaturpc[.]ch
 Full link: www[.]https://reparaturpc[.]ch/de/?msclkid=75c3ce8f8942156ac179ab7f41a03704

* Search engine: Yahoo
 Search term: apple
 Domain: fust[.]ch
 Full link: https://www[.]fust[.]ch/de/marken/apple[.]html?&msclkid=a836011a07061ba4052864eacfe7d0fd&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=apple&utm_content=1_Apple%3D2_undefined%C2%A63_Nbrand&gclid=a836011a07061ba4052864eacfe7d0fd&gclsrc=3p[.]ds

* Search engine: Yahoo
 Search term: apple
 Domain: jobs[.]ch
 Full link: https://www[.]jobs[.]ch/en/vacancies/?term=apple&utm_source=bing&utm_medium=search&utm_campaign=wb:jobs|tg:b2c|cn:ww|lg:en|ct:search,nonbrand,company|cd:company|mg:job-application|pd:y|tt:cpc|gt:keyword,nonbrand,company|gd:company&msclkid=17ac4d7d0b0616628f40288dc3e79a46&utm_term=apple&utm_content=gt%3Akeyword,nonbrand,company%7Cgd%3Acompany

* Search engine: Yahoo
 Search term: apple
 Domain: amazon[.]com
 Full link: https://www[.]amazon[.]com/s?k=applwe&adgrpid=1344703557775981&hvadid=84044278562817&hvbmt=be&hvdev=c&hvlocphy=3322&hvnetw=o&hvqmt=e&hvtargid=kwd-84044521042995%3Aloc-175&hydadcr=29387_14610683&tag=mh0b-20&ref=pd_sl_7xha1yy51_e


This message was automatically sent by seads (www.github.com/andpalmier/seads)
```


## 3rd party libraries

- Rod: [GitHub repo](https://github.com/go-rod/rod), [documentation](https://go-rod.github.io/)
- Shoutrrr: [GitHub repo](https://github.com/containrrr/shoutrrr), [documentation](https://containrrr.dev/shoutrrr/v0.8/)
- Fatih/color: [GitHub repo](https://github.com/fatih/color), [Go reference](https://pkg.go.dev/github.com/fatih/color)

## Next steps

- [ ] Add flag to allow submission of User-Agent string to be used in headless browsers
- [ ] Add more search engines
