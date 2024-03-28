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
- **Export in JSON**: Export results of the execution in JSON format.
- **Multiple User-Agent support**: Provide your User-Agent string to click ads found in search engines. 

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
  -out
        path of JSON file containing links of gathered ads
  -screenshot string
    	path to store screenshots (if empty, the screenshot feature will be disabled).
  -ua string
        User-Agent string to be used to click on ads.
```

Example:

```
seads -config config.yaml -notify
```

Docker example:

```
docker run -v "$(pwd)":/mnt seads -config /mnt/config.yaml -notify
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
seads -config config.yaml -screenshot scr -out results.json -notify
```

or

```
docker run -v "$(pwd)":/mnt seads -config /mnt/config.yaml -screenshot /mnt/scr -out /mnt/results.json -notify
```

output example:

![seads](https://github.com/andpalmier/seads/blob/main/img/seads.gif?raw=true)

screenshot example:

![seads_bing_ipad](https://github.com/andpalmier/seads/blob/main/img/example-bing-ipad.png?raw=true)

notification example:

```
Here are the "unexpected domains" found during the last execution of seads:

Message creation date: 2024-03-28 16:19:59

* Search engine: Bing
 Search term: ipad
 Domain: fust[.]ch
 Full link: https://www[.]fust[.]ch/de/r/pc-tablet-handy/tablet/apple-ipad-455[.]html?gclid=104a09e8ed7b13b08aec6fed67fe1784&gclsrc=3p[.]ds&msclkid=104a09e8ed7b13b08aec6fed67fe1784&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand

* Search engine: Bing
 Search term: ipad
 Domain: interdiscount[.]ch
 Full link: https://www[.]interdiscount[.]ch/de/tablet--c512000?gclsrc=aw[.]ds&gclsrc=3p[.]ds&msclkid=af4155bd2b4d1d3c297a33b8a1c87a49

* Search engine: Bing
 Search term: ipad
 Domain: amazon[.]de
 Full link: https://www[.]amazon[.]de/s?k=ipad&adgrpid=1189672236393152&hvadid=74354631353614&hvbmt=be&hvdev=c&hvlocphy=3322&hvnetw=o&hvqmt=e&hvtargid=kwd-74354720760700%3Aloc-175&hydadcr=29225_2368270&tag=hyddemsn-21&ref=pd_sl_876by72h0s_e

* Search engine: Bing
 Search term: ipad
 Domain: online-preisvergleich[.]de
 Full link: https://online-preisvergleich[.]de/search?q=iPad+Air&em_src=kw&em_cmp=microsoft/Computer&amp;Software&mc=cGiYaMHb0lOd&mscampaign=422647285&msadgroup=1163284371005528&msdevice=c&mstargetid=kwd-72705931935350:loc-175&msmatchtype=p&msnetwork=o&msclid=bb959cce59fc1e843813facbd3797f11&msclkid=bb959cce59fc1e843813facbd3797f11

* Search engine: Yahoo
 Search term: ipad
 Domain: fust[.]ch
 Full link: https://www[.]fust[.]ch/de/r/pc-tablet-handy/tablet/apple-ipad-455[.]html?&msclkid=b0b52ff226f111e6258d49f7ac6044da&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand&gclid=b0b52ff226f111e6258d49f7ac6044da&gclsrc=3p[.]ds

* Search engine: DuckDuckGo
 Search term: ipad
 Domain: fust[.]ch
 Full link: https://www[.]fust[.]ch/de/r/pc-tablet-handy/tablet/apple-ipad-455[.]html?&msclkid=7a0a9f42347e10eee57ff947bee681a9&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand&gclid=7a0a9f42347e10eee57ff947bee681a9&gclsrc=3p[.]ds

* Search engine: DuckDuckGo
 Search term: ipad
 Domain: manor[.]ch
 Full link: https://www[.]manor[.]ch/de/apple/ipad/b/apple/apple-ipad-ipad?msclkid=f58b08acfac4136d7cc96ea9705cdffb&utm_source=bing&utm_medium=cpc&utm_campaign=DE+-+GEN+-+MULTIMEDIA&utm_term=ipad+tablets&utm_content=Brand+-+Apple+-+iPad


This message was automatically sent by seads (www.github.com/andpalmier/seads)
```

output file example:

```json
[
  {
    "engine": "Google",
    "query": "ipad",
    "domain": "apple.com",
    "link": "https://www.apple.com/chde/ipad/",
    "time": "2024-03-28T16:19:32.815712+01:00"
  },
  {
    "engine": "Bing",
    "query": "ipad",
    "domain": "fust.ch",
    "link": "https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?gclid=104a09e8ed7b13b08aec6fed67fe1784\\u0026gclsrc=3p.ds\\u0026msclkid=104a09e8ed7b13b08aec6fed67fe1784\\u0026utm_source=bing\\u0026utm_medium=cpc\\u0026utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple\\u0026utm_term=ipad\\u0026utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand",
    "time": "2024-03-28T16:19:44.295591+01:00"
  },
  {
    "engine": "Bing",
    "query": "ipad",
    "domain": "interdiscount.ch",
    "link": "https://www.interdiscount.ch/de/tablet--c512000?gclsrc=aw.ds\\u0026gclsrc=3p.ds\\u0026msclkid=af4155bd2b4d1d3c297a33b8a1c87a49",
    "time": "2024-03-28T16:19:44.295591+01:00"
  },
  {
    "engine": "Bing",
    "query": "ipad",
    "domain": "amazon.de",
    "link": "https://www.amazon.de/s?k=ipad\\u0026adgrpid=1189672236393152\\u0026hvadid=74354631353614\\u0026hvbmt=be\\u0026hvdev=c\\u0026hvlocphy=3322\\u0026hvnetw=o\\u0026hvqmt=e\\u0026hvtargid=kwd-74354720760700%3Aloc-175\\u0026hydadcr=29225_2368270\\u0026tag=hyddemsn-21\\u0026ref=pd_sl_876by72h0s_e",
    "time": "2024-03-28T16:19:44.295591+01:00"
  },
  {
    "engine": "Bing",
    "query": "ipad",
    "domain": "online-preisvergleich.de",
    "link": "https://online-preisvergleich.de/search?q=iPad+Air\\u0026em_src=kw\\u0026em_cmp=microsoft/Computer\\u0026amp;Software\\u0026mc=cGiYaMHb0lOd\\u0026mscampaign=422647285\\u0026msadgroup=1163284371005528\\u0026msdevice=c\\u0026mstargetid=kwd-72705931935350:loc-175\\u0026msmatchtype=p\\u0026msnetwork=o\\u0026msclid=bb959cce59fc1e843813facbd3797f11\\u0026msclkid=bb959cce59fc1e843813facbd3797f11",
    "time": "2024-03-28T16:19:44.295591+01:00"
  },
  {
    "engine": "Yahoo",
    "query": "ipad",
    "domain": "amazon.com",
    "link": "https://www.amazon.com/s?k=%E2%80%9Cipad%E2%80%9D\\u0026i=electronics\\u0026adgrpid=1340306317325785\\u0026hvadid=83769406602961\\u0026hvbmt=be\\u0026hvdev=c\\u0026hvlocphy=3322\\u0026hvnetw=o\\u0026hvqmt=e\\u0026hvtargid=kwd-83770149440260%3Aloc-175\\u0026hydadcr=8776_10900136\\u0026tag=mh0b-20\\u0026ref=pd_sl_43gcbo4bh9_e",
    "time": "2024-03-28T16:19:52.044445+01:00"
  },
  {
    "engine": "Yahoo",
    "query": "ipad",
    "domain": "fust.ch",
    "link": "https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?\\u0026msclkid=b0b52ff226f111e6258d49f7ac6044da\\u0026utm_source=bing\\u0026utm_medium=cpc\\u0026utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple\\u0026utm_term=ipad\\u0026utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand\\u0026gclid=b0b52ff226f111e6258d49f7ac6044da\\u0026gclsrc=3p.ds",
    "time": "2024-03-28T16:19:52.044445+01:00"
  },
  {
    "engine": "DuckDuckGo",
    "query": "ipad",
    "domain": "fust.ch",
    "link": "https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?\\u0026msclkid=7a0a9f42347e10eee57ff947bee681a9\\u0026utm_source=bing\\u0026utm_medium=cpc\\u0026utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple\\u0026utm_term=ipad\\u0026utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand\\u0026gclid=7a0a9f42347e10eee57ff947bee681a9\\u0026gclsrc=3p.ds",
    "time": "2024-03-28T16:19:59.847625+01:00"
  },
  {
    "engine": "DuckDuckGo",
    "query": "ipad",
    "domain": "manor.ch",
    "link": "https://www.manor.ch/de/apple/ipad/b/apple/apple-ipad-ipad?msclkid=f58b08acfac4136d7cc96ea9705cdffb\\u0026utm_source=bing\\u0026utm_medium=cpc\\u0026utm_campaign=DE+-+GEN+-+MULTIMEDIA\\u0026utm_term=ipad+tablets\\u0026utm_content=Brand+-+Apple+-+iPad",
    "time": "2024-03-28T16:19:59.847625+01:00"
  }
]
```

## 3rd party libraries

- Rod: [GitHub repo](https://github.com/go-rod/rod), [documentation](https://go-rod.github.io/)
- Shoutrrr: [GitHub repo](https://github.com/containrrr/shoutrrr), [documentation](https://containrrr.dev/shoutrrr/v0.8/)
- Fatih/color: [GitHub repo](https://github.com/fatih/color), [Go reference](https://pkg.go.dev/github.com/fatih/color)
