# Original from https://github.com/andpalmier/seads

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
- **Docker Support**: Install `seads` without affecting your local setup using Docker.
- **Export in JSON**: Export results of the execution in JSON format.
- **Multiple User-Agent support**: Provide your User-Agent string to click ads found in search engines.
- **Redirect chain detection**: Tracks URLs through redirects to detect and log chains. 
- **URLScan submission**: Submit to URLScan using API key.

## Known limitations:
- Due to the nature of search engine ads, a single search may not reveal all ads. Using concurrent headless browsers might slow down detection but ensures comprehensive ad gathering.
- Notifications on Slack and Telegram have character limits. Messages exceeding the limit won't be sent.

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
      path to config file (default "config.yaml")
  -concurrency int
    	number of concurrent headless browsers (default 4)
  -cleanlinks
    	print clear links in output (links will remain defanged in notifications)
  -notify
    	notify if unexpected domains are found
  -out
      path of JSON file containing links of gathered ads
  -screenshot string
    	path to store screenshots (if empty, the screenshot feature will be disabled)
  -redirect
      print redirection chain for ad links found
  -ua string
      User-Agent string to be used to click on ads
  -urlscan
      submit url to urlscan.io for analysis
  -html string
      path to store search engine result html page
  -noredirection
      do not follow redirection, if URLScan submit link to resolve by URLScan instead
```

Example:

```
seads -config config.yaml -notify
```

Example with urlscan with redirection handled chain by URLScan and save html page and screenshot on disk (assuming the folder screenshots and htmls already created):

```
seads -config config.yaml -urlscan -noredirection -screenshot ./screenshots -html ./htmls
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
  chatids: [CHATID#1,CHATID#2]

urlscan:
  token: "APIKEY"
  scanurl: "https://urlscan.io/api/v1/scan/"
  visibility: "unlisted"
  tags: "seads_ads_tracker"

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

Message creation date: 2024-06-07 18:52:51

* Search engine: Google
 Search term: ipad
 Domain: revendo[.]com
 Full link: https://revendo[.]com/de-ch/kaufen/k/ipad/?utm_source=google&utm_medium=cpc&utm_campaign=CH_ACQ_230710_GOOGLEADS_SEARCH_IPAD_WEB_Z00_P02_RSA_1%20_&utm_term=ipad&utm_content=653996737162&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYASAAEgIT9PD_BwE

* Search engine: Google
 Search term: ipad
 Domain: galaxus[.]ch
 Full link: https://www[.]galaxus[.]ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_DSA&campaignid=146803917&adgroupid=7678979517&adid=686471669496&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gclsrc=aw[.]ds

* Search engine: Google
 Search term: ipad
 Domain: digitec[.]ch
 Full link: https://www[.]digitec[.]ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_Dynamische+Suchanzeigen+Test&campaignid=12557018154&adgroupid=120228040035&adid=506763570197&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gclsrc=aw[.]ds

* Search engine: Google
 Search term: ipad
 Domain: interdiscount[.]ch
 Full link: https://www[.]interdiscount[.]ch/de/cms/apple-ipad?gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYBCAAEgKfZvD_BwE&gclsrc=aw[.]ds

* Search engine: DuckDuckGo
 Search term: ipad
 Domain: temu[.]com
 Full link: https://www[.]temu[.]com/ul/kuiper/un2[.]html?_p_rfs=1&subj=un-search-web&_p_jump_id=960&_x_vst_scene=adg&locale_override=192~en~CHF&search_key=%22ipad%22&_p_rfs=1&_x_ads_channel=bing&_x_ads_sub_channel=search&_x_ads_account=176283034&_x_ads_set=520457070&_x_ads_id=1321615336742536&_x_ads_creative_id=82601185328809&_x_ns_source=s&_x_ns_msclkid=5ca4051bf8991674f5c16378c2572438&_x_ns_match_type=e&_x_ns_bid_match_type=be&_x_ns_query=ipad&_x_ns_keyword=%22ipad%22&_x_ns_device=c&_x_ns_targetid=kwd-82601954111094:loc-175&_x_ns_extensionid=&msclkid=5ca4051bf8991674f5c16378c2572438&utm_source=bing&utm_medium=cpc&utm_campaign=0228_SEARCH_3038896834118172839&utm_term=%22ipad%22&utm_content=0228_SEARCH_1144514176907441598

* Search engine: DuckDuckGo
 Search term: ipad
 Domain: telekom[.]de
 Full link: https://www[.]telekom[.]de/shop/tablets/apple/apple-11-ipad-pro-2024-wifi-plus-5g/silber-256-gb?tariffId=MF_16213&msclkid=4def2dbdb77910e3fc7e4bde8db7216e&gclid=4def2dbdb77910e3fc7e4bde8db7216e&gclsrc=3p[.]ds&autoLogin=true&error=interaction_required&state=b7c95c0d-f271-4cba-873e-edff618f4e5b

* Search engine: DuckDuckGo
 Search term: ipad
 Domain: fust[.]ch
 Full link: https://www[.]fust[.]ch/de/r/pc-tablet-handy/tablet/apple-ipad-455[.]html?&msclkid=5bbafe3fa993173dfb6c985c935d4546&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand&gclid=5bbafe3fa993173dfb6c985c935d4546&gclsrc=3p[.]ds

This message was automatically sent by seads (github.com/andpalmier/seads)
```

output file example:

```json
[
  {
    "engine": "Google",
    "query": "ipad",
    "domain-final-url": "revendo.com",
    "final-url": "https://revendo.com/de-ch/kaufen/k/ipad/?utm_source=google&utm_medium=cpc&utm_campaign=CH_ACQ_230710_GOOGLEADS_SEARCH_IPAD_WEB_Z00_P02_RSA_1%20_&utm_term=ipad&utm_content=653996737162&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYASAAEgIT9PD_BwE",
    "redirect-chain": [
      "https://www.google.com/aclk?sa=L&ai=DChcSEwjT2Zb2-cmGAxWLlVAGHQDeB9MYABADGgJkZw&ase=2&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYASAAEgIT9PD_BwE&sig=AOD64_1H3DeWzBC8mcJbJqpLVSZCeAz0jA&q&nis=4&adurl",
      "https://monitor.clickcease.com/tracker/tracker.aspx?id=weAYNUhEaokBfK&adpos=&locphisical=9186940&locinterest=&adgrp=146223459485&kw=ipad&nw=g&url=https://revendo.com/de-ch/kaufen/k/ipad/%3Futm_source%3Dgoogle%26utm_medium%3Dcpc%26utm_campaign%3DCH_ACQ_230710_GOOGLEADS_SEARCH_IPAD_WEB_Z00_P02_RSA_1%2520_%26utm_term%3Dipad%26utm_content%3D653996737162%26gad_source%3D1&cpn=19934623633&device=c&ccpturl=revendo.ch&pl=&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYASAAEgIT9PD_BwE"
    ],
    "time": "2024-06-07T18:52:16.647178+02:00"
  },
  {
    "engine": "Google",
    "query": "ipad",
    "domain-final-url": "galaxus.ch",
    "final-url": "https://www.galaxus.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_DSA&campaignid=146803917&adgroupid=7678979517&adid=686471669496&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gclsrc=aw.ds",
    "redirect-chain": [
      "https://www.google.com/aclk?sa=L&ai=DChcSEwjT2Zb2-cmGAxWLlVAGHQDeB9MYABACGgJkZw&ase=2&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&sig=AOD64_0GTl02lMiPj9OkeIVO5z1ZZAW4Wg&q&nis=4&adurl",
      "https://ad.doubleclick.net/searchads/link/click?lid=39700007760435903&ds_s_kwgid=58700000631575915&ds_a_cid=106339197&ds_a_caid=146803917&ds_a_agid=7678979517&ds_a_fiid=&ds_a_lid=dsa-19959388920&ds_a_extid=&&ds_e_adid=686471669496&ds_e_matchtype=search&ds_e_device=c&ds_e_network=g&&ds_url_v=2&dc_eps=AHas8cDYKkEXNq_0BRZxpXAvmLD8QjvkIEB6p1pUFN2iYxABq8TRM6R4358dqROWC6z8ZG24cTDYspBb7oXS&acs_info=ZmluYWxfdXJsOiAiaHR0cHM6Ly93d3cuZ2FsYXh1cy5jaC9kZS9zMS9wcm9kdWN0dHlwZS90YWJsZXQtNDY5Igo&ds_dest_url=https://www.galaxus.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_DSA&campaignid=146803917&adgroupid=7678979517&adid=686471669496&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gclsrc=aw.ds&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&cls=1",
      "https://www.galaxus.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_DSA&campaignid=146803917&adgroupid=7678979517&adid=686471669496&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAiAAEgLG4_D_BwE&gclsrc=aw.ds"
    ],
    "time": "2024-06-07T18:52:16.647178+02:00"
  },
  {
    "engine": "Google",
    "query": "ipad",
    "domain-final-url": "digitec.ch",
    "final-url": "https://www.digitec.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_Dynamische+Suchanzeigen+Test&campaignid=12557018154&adgroupid=120228040035&adid=506763570197&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gclsrc=aw.ds",
    "redirect-chain": [
      "https://www.google.com/aclk?sa=L&ai=DChcSEwjT2Zb2-cmGAxWLlVAGHQDeB9MYABAAGgJkZw&ase=2&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&sig=AOD64_14NOo9WlciN55IEHg6h_gLW_LqlA&q&nis=4&adurl",
      "https://ad.doubleclick.net/searchads/link/click?lid=39700062263657243&ds_s_kwgid=58700006926909766&ds_a_cid=104953610&ds_a_caid=12557018154&ds_a_agid=120228040035&ds_a_fiid=&ds_a_lid=dsa-1199847883174&ds_a_extid=&&ds_e_adid=506763570197&ds_e_matchtype=search&ds_e_device=c&ds_e_network=g&&ds_url_v=2&dc_eps=AHas8cCXbR1ZBW0_9YXKT4muD0AjNmJVQY_fWxhgJJy3BtV1rujSqP1NJcSFnHPm4k8KApQbbyBuNYpyZnt9&acs_info=ZmluYWxfdXJsOiAiaHR0cHM6Ly93d3cuZGlnaXRlYy5jaC9kZS9zMS9wcm9kdWN0dHlwZS90YWJsZXQtNDY5Igo&ds_dest_url=https://www.digitec.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_Dynamische+Suchanzeigen+Test&campaignid=12557018154&adgroupid=120228040035&adid=506763570197&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gclsrc=aw.ds&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&cls=1",
      "https://www.digitec.ch/de/s1/producttype/tablet-469?filter=t_bra%3D47&utm_source=google&utm_medium=cpc&utm_campaign=SEA_DE_CH_Dynamische+Suchanzeigen+Test&campaignid=12557018154&adgroupid=120228040035&adid=506763570197&dgCidg=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYAyAAEgJFGfD_BwE&gclsrc=aw.ds"
    ],
    "time": "2024-06-07T18:52:16.647178+02:00"
  },
  {
    "engine": "Google",
    "query": "ipad",
    "domain-final-url": "interdiscount.ch",
    "final-url": "https://www.interdiscount.ch/de/cms/apple-ipad?gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYBCAAEgKfZvD_BwE&gclsrc=aw.ds",
    "redirect-chain": [
      "https://www.google.com/aclk?sa=L&ai=DChcSEwjT2Zb2-cmGAxWLlVAGHQDeB9MYABABGgJkZw&ase=2&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYBCAAEgKfZvD_BwE&sig=AOD64_1b1BaTE1d6dNkO7WGWh-_zpG4MBA&q&nis=4&adurl",
      "https://ad.doubleclick.net/searchads/link/click?lid=43700075248202559&ds_s_kwgid=58700008279170537&ds_a_cid=405395944&ds_a_caid=9221661016&ds_a_agid=145824926826&ds_a_fiid=&ds_a_lid=kwd-76826760&ds_a_extid=&&ds_e_adid=648877665266&ds_e_matchtype=search&ds_e_device=c&ds_e_network=g&&ds_url_v=2&dc_eps=AHas8cAVCSY3WjuQFtA8d1sJPHlZEHXC9K8gWM9Zq1RzhRVKVT71n6vjdTKXGhoLc8Pgv_tE2OoFp3oQjfAb&acs_info=ZmluYWxfdXJsOiAiaHR0cHM6Ly93d3cuaW50ZXJkaXNjb3VudC5jaC9kZS9jbXMvYXBwbGUtaXBhZCIK&ds_dest_url=https://www.interdiscount.ch/de/cms/apple-ipad?gclsrc=aw.ds&gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYBCAAEgKfZvD_BwE&cls=1",
      "https://www.interdiscount.ch/de/cms/apple-ipad?gad_source=1&gclid=EAIaIQobChMI09mW9vnJhgMVi5VQBh0A3gfTEAAYBCAAEgKfZvD_BwE&gclsrc=aw.ds"
    ],
    "time": "2024-06-07T18:52:16.647178+02:00"
  },
  {
    "engine": "DuckDuckGo",
    "query": "ipad",
    "domain-final-url": "temu.com",
    "final-url": "https://www.temu.com/ul/kuiper/un2.html?_p_rfs=1&subj=un-search-web&_p_jump_id=960&_x_vst_scene=adg&locale_override=192~en~CHF&search_key=%22ipad%22&_p_rfs=1&_x_ads_channel=bing&_x_ads_sub_channel=search&_x_ads_account=176283034&_x_ads_set=520457070&_x_ads_id=1321615336742536&_x_ads_creative_id=82601185328809&_x_ns_source=s&_x_ns_msclkid=5ca4051bf8991674f5c16378c2572438&_x_ns_match_type=e&_x_ns_bid_match_type=be&_x_ns_query=ipad&_x_ns_keyword=%22ipad%22&_x_ns_device=c&_x_ns_targetid=kwd-82601954111094:loc-175&_x_ns_extensionid=&msclkid=5ca4051bf8991674f5c16378c2572438&utm_source=bing&utm_medium=cpc&utm_campaign=0228_SEARCH_3038896834118172839&utm_term=%22ipad%22&utm_content=0228_SEARCH_1144514176907441598",
    "redirect-chain": [
      "https://duckduckgo.com/y.js?ad_domain=temu.com&ad_provider=bingv7aa&ad_type=txad&eddgt=8dwnGAKBaZb8InkWQQeEvQ%3D%3D&rut=1eb14dd44bc3f1c119b30beb8c7d9017a223674d069ce8777a84cbd910442aff&u3=https%3A%2F%2Fwww.bing.com%2Faclick%3Fld%3De8ZZchXS931psPvLB1vLOxLzVUCUwNZ685j%2DNXuholOZGG%2DZ7GMwQ8FASSW1LTHIVBle6jas5IX_JHqkv9Zo24muMIsOD0SL0JzPyYuyPK5F6UZuiZecA7DrK1%2DfcZQ9hsJiQi_9hk629L8wjJmQptDmNi8hGa2ilCAqawQCmEGmoN0VrxeoUuBhTaV3mO7dS4vyij2g%26u%3DaHR0cHMlM2ElMmYlMmZ3d3cudGVtdS5jb20lMmZ1bCUyZmt1aXBlciUyZnVuMi5odG1sJTNmX3BfcmZzJTNkMSUyNnN1YmolM2R1bi1zZWFyY2gtd2ViJTI2X3BfanVtcF9pZCUzZDk2MCUyNl94X3ZzdF9zY2VuZSUzZGFkZyUyNmxvY2FsZV9vdmVycmlkZSUzZDE5MiU3ZWVuJTdlQ0hGJTI2c2VhcmNoX2tleSUzZCUyNTIyaXBhZCUyNTIyJTI2X3BfcmZzJTNkMSUyNl94X2Fkc19jaGFubmVsJTNkYmluZyUyNl94X2Fkc19zdWJfY2hhbm5lbCUzZHNlYXJjaCUyNl94X2Fkc19hY2NvdW50JTNkMTc2MjgzMDM0JTI2X3hfYWRzX3NldCUzZDUyMDQ1NzA3MCUyNl94X2Fkc19pZCUzZDEzMjE2MTUzMzY3NDI1MzYlMjZfeF9hZHNfY3JlYXRpdmVfaWQlM2Q4MjYwMTE4NTMyODgwOSUyNl94X25zX3NvdXJjZSUzZHMlMjZfeF9uc19tc2Nsa2lkJTNkNWNhNDA1MWJmODk5MTY3NGY1YzE2Mzc4YzI1NzI0MzglMjZfeF9uc19tYXRjaF90eXBlJTNkZSUyNl94X25zX2JpZF9tYXRjaF90eXBlJTNkYmUlMjZfeF9uc19xdWVyeSUzZGlwYWQlMjZfeF9uc19rZXl3b3JkJTNkJTI1MjJpcGFkJTI1MjIlMjZfeF9uc19kZXZpY2UlM2RjJTI2X3hfbnNfdGFyZ2V0aWQlM2Rrd2QtODI2MDE5NTQxMTEwOTQlM2Fsb2MtMTc1JTI2X3hfbnNfZXh0ZW5zaW9uaWQlM2QlMjZtc2Nsa2lkJTNkNWNhNDA1MWJmODk5MTY3NGY1YzE2Mzc4YzI1NzI0MzglMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2QwMjI4X1NFQVJDSF8zMDM4ODk2ODM0MTE4MTcyODM5JTI2dXRtX3Rlcm0lM2QlMjUyMmlwYWQlMjUyMiUyNnV0bV9jb250ZW50JTNkMDIyOF9TRUFSQ0hfMTE0NDUxNDE3NjkwNzQ0MTU5OA%26rlid%3D5ca4051bf8991674f5c16378c2572438&vqd=4-166102188016361480819260538901484848581&iurl=%7B1%7DIG%3D216F730E9ADC4BB09EACBFEC50EE9BF8%26CID%3D071F188F7FBE601E01470C187E53619E%26ID%3DDevEx%2C5055.1",
      "https://www.bing.com/aclick?ld=e8ZZchXS931psPvLB1vLOxLzVUCUwNZ685j-NXuholOZGG-Z7GMwQ8FASSW1LTHIVBle6jas5IX_JHqkv9Zo24muMIsOD0SL0JzPyYuyPK5F6UZuiZecA7DrK1-fcZQ9hsJiQi_9hk629L8wjJmQptDmNi8hGa2ilCAqawQCmEGmoN0VrxeoUuBhTaV3mO7dS4vyij2g&u=aHR0cHMlM2ElMmYlMmZ3d3cudGVtdS5jb20lMmZ1bCUyZmt1aXBlciUyZnVuMi5odG1sJTNmX3BfcmZzJTNkMSUyNnN1YmolM2R1bi1zZWFyY2gtd2ViJTI2X3BfanVtcF9pZCUzZDk2MCUyNl94X3ZzdF9zY2VuZSUzZGFkZyUyNmxvY2FsZV9vdmVycmlkZSUzZDE5MiU3ZWVuJTdlQ0hGJTI2c2VhcmNoX2tleSUzZCUyNTIyaXBhZCUyNTIyJTI2X3BfcmZzJTNkMSUyNl94X2Fkc19jaGFubmVsJTNkYmluZyUyNl94X2Fkc19zdWJfY2hhbm5lbCUzZHNlYXJjaCUyNl94X2Fkc19hY2NvdW50JTNkMTc2MjgzMDM0JTI2X3hfYWRzX3NldCUzZDUyMDQ1NzA3MCUyNl94X2Fkc19pZCUzZDEzMjE2MTUzMzY3NDI1MzYlMjZfeF9hZHNfY3JlYXRpdmVfaWQlM2Q4MjYwMTE4NTMyODgwOSUyNl94X25zX3NvdXJjZSUzZHMlMjZfeF9uc19tc2Nsa2lkJTNkNWNhNDA1MWJmODk5MTY3NGY1YzE2Mzc4YzI1NzI0MzglMjZfeF9uc19tYXRjaF90eXBlJTNkZSUyNl94X25zX2JpZF9tYXRjaF90eXBlJTNkYmUlMjZfeF9uc19xdWVyeSUzZGlwYWQlMjZfeF9uc19rZXl3b3JkJTNkJTI1MjJpcGFkJTI1MjIlMjZfeF9uc19kZXZpY2UlM2RjJTI2X3hfbnNfdGFyZ2V0aWQlM2Rrd2QtODI2MDE5NTQxMTEwOTQlM2Fsb2MtMTc1JTI2X3hfbnNfZXh0ZW5zaW9uaWQlM2QlMjZtc2Nsa2lkJTNkNWNhNDA1MWJmODk5MTY3NGY1YzE2Mzc4YzI1NzI0MzglMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2QwMjI4X1NFQVJDSF8zMDM4ODk2ODM0MTE4MTcyODM5JTI2dXRtX3Rlcm0lM2QlMjUyMmlwYWQlMjUyMiUyNnV0bV9jb250ZW50JTNkMDIyOF9TRUFSQ0hfMTE0NDUxNDE3NjkwNzQ0MTU5OA&rlid=5ca4051bf8991674f5c16378c2572438",
      "https://www.temu.com/ul/kuiper/un2.html?_p_rfs=1&subj=un-search-web&_p_jump_id=960&_x_vst_scene=adg&locale_override=192~en~CHF&search_key=%22ipad%22&_p_rfs=1&_x_ads_channel=bing&_x_ads_sub_channel=search&_x_ads_account=176283034&_x_ads_set=520457070&_x_ads_id=1321615336742536&_x_ads_creative_id=82601185328809&_x_ns_source=s&_x_ns_msclkid=5ca4051bf8991674f5c16378c2572438&_x_ns_match_type=e&_x_ns_bid_match_type=be&_x_ns_query=ipad&_x_ns_keyword=%22ipad%22&_x_ns_device=c&_x_ns_targetid=kwd-82601954111094:loc-175&_x_ns_extensionid=&msclkid=5ca4051bf8991674f5c16378c2572438&utm_source=bing&utm_medium=cpc&utm_campaign=0228_SEARCH_3038896834118172839&utm_term=%22ipad%22&utm_content=0228_SEARCH_1144514176907441598"
    ],
    "time": "2024-06-07T18:52:35.825523+02:00"
  },
  {
    "engine": "DuckDuckGo",
    "query": "ipad",
    "domain-final-url": "telekom.de",
    "final-url": "https://www.telekom.de/shop/tablets/apple/apple-11-ipad-pro-2024-wifi-plus-5g/silber-256-gb?tariffId=MF_16213&msclkid=4def2dbdb77910e3fc7e4bde8db7216e&gclid=4def2dbdb77910e3fc7e4bde8db7216e&gclsrc=3p.ds&autoLogin=true&error=interaction_required&state=b7c95c0d-f271-4cba-873e-edff618f4e5b",
    "redirect-chain": [
      "https://duckduckgo.com/y.js?ad_domain=telekom.de&ad_provider=bingv7aa&ad_type=txad&eddgt=8dwnGAKBaZb8InkWQQeEvQ%3D%3D&rut=d9b404be3f9031cfab0986d8da4aab9635f46cd218894f7ddf12117aa657db55&u3=https%3A%2F%2Fwww.bing.com%2Faclick%3Fld%3De8gZyRkbqVLt2sM492q43GozVUCUzIIfsNdNKBOiAOBTlGBuWC7zaMJfhqwb7f0UkMAzqgal1KihWhWvMvCV%2DxAEFDf7Vy3zUnzD5%2Df70J6HpR4u6lXPc8jennbbCtFB8J8tQeYsbJmIpPe57knbxX%2DVTRpnF2IAjrh6W%2Dd7070PRvAsaJZaQ1ornuZqeE1VBVd9zRWA%26u%3DaHR0cHMlM2ElMmYlMmZhZC5kb3VibGVjbGljay5uZXQlMmZzZWFyY2hhZHMlMmZsaW5rJTJmY2xpY2slM2ZsaWQlM2Q0MzcwMDA4MDE0NjM0MjYyOCUyNmRzX3Nfa3dnaWQlM2Q1ODcwMDAwODcxNzE2OTQ3MiUyNmRzX2FfY2lkJTNkMTIyNDU4MTcyMSUyNmRzX2FfY2FpZCUzZDIxMjk1MzIxMTA5JTI2ZHNfYV9hZ2lkJTNkMTYyMDgyNTg1MTIzJTI2ZHNfYV9saWQlM2Rrd2QtNzY4MjY3NjAlMjYlMjZkc19lX2FkaWQlM2Q3MzE4NjY0OTU5ODI5NyUyNmRzX2VfdGFyZ2V0X2lkJTNka3dkLTczMTg3MDc5Nzg3ODg4JTI2JTI2ZHNfZV9uZXR3b3JrJTNkcyUyNmRzX3VybF92JTNkMiUyNmRzX2Rlc3RfdXJsJTNkaHR0cHMlM2ElMmYlMmZ3d3cudGVsZWtvbS5kZSUyZnNob3AlMmZ0YWJsZXRzJTJmYXBwbGUlMmZhcHBsZS0xMS1pcGFkLXByby0yMDI0LXdpZmktcGx1cy01ZyUyZnNpbGJlci0yNTYtZ2IlM2Z0YXJpZmZJZCUzZE1GXzE2MjEzJTI2Z2NsaWQlM2Q0ZGVmMmRiZGI3NzkxMGUzZmM3ZTRiZGU4ZGI3MjE2ZSUyNmdjbHNyYyUzZDNwLmRzJTI2JTI2bXNjbGtpZCUzZDRkZWYyZGJkYjc3OTEwZTNmYzdlNGJkZThkYjcyMTZl%26rlid%3D4def2dbdb77910e3fc7e4bde8db7216e&vqd=4-239786219880205881155441526817052983105&iurl=%7B1%7DIG%3D216F730E9ADC4BB09EACBFEC50EE9BF8%26CID%3D071F188F7FBE601E01470C187E53619E%26ID%3DDevEx%2C5078.1",
      "https://www.bing.com/aclick?ld=e8gZyRkbqVLt2sM492q43GozVUCUzIIfsNdNKBOiAOBTlGBuWC7zaMJfhqwb7f0UkMAzqgal1KihWhWvMvCV-xAEFDf7Vy3zUnzD5-f70J6HpR4u6lXPc8jennbbCtFB8J8tQeYsbJmIpPe57knbxX-VTRpnF2IAjrh6W-d7070PRvAsaJZaQ1ornuZqeE1VBVd9zRWA&u=aHR0cHMlM2ElMmYlMmZhZC5kb3VibGVjbGljay5uZXQlMmZzZWFyY2hhZHMlMmZsaW5rJTJmY2xpY2slM2ZsaWQlM2Q0MzcwMDA4MDE0NjM0MjYyOCUyNmRzX3Nfa3dnaWQlM2Q1ODcwMDAwODcxNzE2OTQ3MiUyNmRzX2FfY2lkJTNkMTIyNDU4MTcyMSUyNmRzX2FfY2FpZCUzZDIxMjk1MzIxMTA5JTI2ZHNfYV9hZ2lkJTNkMTYyMDgyNTg1MTIzJTI2ZHNfYV9saWQlM2Rrd2QtNzY4MjY3NjAlMjYlMjZkc19lX2FkaWQlM2Q3MzE4NjY0OTU5ODI5NyUyNmRzX2VfdGFyZ2V0X2lkJTNka3dkLTczMTg3MDc5Nzg3ODg4JTI2JTI2ZHNfZV9uZXR3b3JrJTNkcyUyNmRzX3VybF92JTNkMiUyNmRzX2Rlc3RfdXJsJTNkaHR0cHMlM2ElMmYlMmZ3d3cudGVsZWtvbS5kZSUyZnNob3AlMmZ0YWJsZXRzJTJmYXBwbGUlMmZhcHBsZS0xMS1pcGFkLXByby0yMDI0LXdpZmktcGx1cy01ZyUyZnNpbGJlci0yNTYtZ2IlM2Z0YXJpZmZJZCUzZE1GXzE2MjEzJTI2Z2NsaWQlM2Q0ZGVmMmRiZGI3NzkxMGUzZmM3ZTRiZGU4ZGI3MjE2ZSUyNmdjbHNyYyUzZDNwLmRzJTI2JTI2bXNjbGtpZCUzZDRkZWYyZGJkYjc3OTEwZTNmYzdlNGJkZThkYjcyMTZl&rlid=4def2dbdb77910e3fc7e4bde8db7216e",
      "https://ad.doubleclick.net/searchads/link/click?lid=43700080146342628&ds_s_kwgid=58700008717169472&ds_a_cid=1224581721&ds_a_caid=21295321109&ds_a_agid=162082585123&ds_a_lid=kwd-76826760&&ds_e_adid=73186649598297&ds_e_target_id=kwd-73187079787888&&ds_e_network=s&ds_url_v=2&ds_dest_url=https://www.telekom.de/shop/tablets/apple/apple-11-ipad-pro-2024-wifi-plus-5g/silber-256-gb?tariffId=MF_16213&gclid=4def2dbdb77910e3fc7e4bde8db7216e&gclsrc=3p.ds&&msclkid=4def2dbdb77910e3fc7e4bde8db7216e",
      "https://www.telekom.de/shop/tablets/apple/apple-11-ipad-pro-2024-wifi-plus-5g/silber-256-gb?tariffId=MF_16213&&msclkid=4def2dbdb77910e3fc7e4bde8db7216e&gclid=4def2dbdb77910e3fc7e4bde8db7216e&gclsrc=3p.ds"
    ],
    "time": "2024-06-07T18:52:35.825523+02:00"
  },
  {
    "engine": "DuckDuckGo",
    "query": "ipad",
    "domain-final-url": "fust.ch",
    "final-url": "https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?&msclkid=5bbafe3fa993173dfb6c985c935d4546&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand&gclid=5bbafe3fa993173dfb6c985c935d4546&gclsrc=3p.ds",
    "redirect-chain": [
      "https://duckduckgo.com/y.js?ad_domain=fust.ch&ad_provider=bingv7aa&ad_type=txad&eddgt=8dwnGAKBaZb8InkWQQeEvQ%3D%3D&rut=ed86c66e71fea51bd2487367820d13e5f31b4d108da29d3ec7cc8d4a2fbe4dc7&u3=https%3A%2F%2Fwww.bing.com%2Faclick%3Fld%3De8JJlg7u7kxQwlqfN9OFBZ0jVUCUzC3aWBbElI961o%2D0KcnqmSyQuyPx_l%2DaNQOnu0X7EAG71aEd7E4AKzeTtJvHI1eX6laBBUsx0CULwgcni0QRG%2DuMA1F185t7UXzxXO0aIKv1bYNv6epFFCXiw4IBB9mgUQV_XKvMRYH2yMKjfx1dSuwPPaM9AABm4bcYpRXb2dVg%26u%3DaHR0cHMlM2ElMmYlMmZhZC5kb3VibGVjbGljay5uZXQlMmZzZWFyY2hhZHMlMmZsaW5rJTJmY2xpY2slM2ZsaWQlM2Q0MzcwMDA3NjE5MTQ0NDQzMCUyNmRzX3Nfa3dnaWQlM2Q1ODcwMDAwODI3MDU2NjkzNCUyNmRzX2FfY2lkJTNkNDExMjA1Mzk1JTI2ZHNfYV9jYWlkJTNkMTk2NjA4ODM0MTYlMjZkc19hX2FnaWQlM2QxNDg4MzgwODcyMzElMjZkc19hX2xpZCUzZGt3ZC0yNTk5Mzc5NiUyNiUyNmRzX2VfYWRpZCUzZDgyOTQ0NzcyODI3MjU0JTI2ZHNfZV90YXJnZXRfaWQlM2Rrd2QtODI5NDU0NDM0MzQwMTAlM2Fsb2MtMTc1JTI2JTI2ZHNfZV9uZXR3b3JrJTNkcyUyNmRzX3VybF92JTNkMiUyNmRzX2Rlc3RfdXJsJTNkaHR0cHMlM2ElMmYlMmZ3d3cuZnVzdC5jaCUyZmRlJTJmciUyZnBjLXRhYmxldC1oYW5keSUyZnRhYmxldCUyZmFwcGxlLWlwYWQtNDU1Lmh0bWwlM2ZnY2xpZCUzZDViYmFmZTNmYTk5MzE3M2RmYjZjOTg1YzkzNWQ0NTQ2JTI2Z2Nsc3JjJTNkM3AuZHMlMjYlMjZtc2Nsa2lkJTNkNWJiYWZlM2ZhOTkzMTczZGZiNmM5ODVjOTM1ZDQ1NDYlMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2RCaW5nJTI1MjAtJTI1MjBOQnJhbmQlMjUyMC0lMjUyMFMlMjUyMC0lMjUyMEQlMjUyMC0lMjUyME1NJTI1MjBQQyUyNTIwTWFya2UlMjUyMEFwcGxlJTI2dXRtX3Rlcm0lM2RpcGFkJTI2dXRtX2NvbnRlbnQlM2QxX0FwcGxlJTI1M0QyX0FwcGxlJTI1MjBpUGFkJWMyJWE2M19OYnJhbmQ%26rlid%3D5bbafe3fa993173dfb6c985c935d4546&vqd=4-150458699415843129215768299163044858549&iurl=%7B1%7DIG%3D216F730E9ADC4BB09EACBFEC50EE9BF8%26CID%3D071F188F7FBE601E01470C187E53619E%26ID%3DDevEx%2C5099.1",
      "https://www.bing.com/aclick?ld=e8JJlg7u7kxQwlqfN9OFBZ0jVUCUzC3aWBbElI961o-0KcnqmSyQuyPx_l-aNQOnu0X7EAG71aEd7E4AKzeTtJvHI1eX6laBBUsx0CULwgcni0QRG-uMA1F185t7UXzxXO0aIKv1bYNv6epFFCXiw4IBB9mgUQV_XKvMRYH2yMKjfx1dSuwPPaM9AABm4bcYpRXb2dVg&u=aHR0cHMlM2ElMmYlMmZhZC5kb3VibGVjbGljay5uZXQlMmZzZWFyY2hhZHMlMmZsaW5rJTJmY2xpY2slM2ZsaWQlM2Q0MzcwMDA3NjE5MTQ0NDQzMCUyNmRzX3Nfa3dnaWQlM2Q1ODcwMDAwODI3MDU2NjkzNCUyNmRzX2FfY2lkJTNkNDExMjA1Mzk1JTI2ZHNfYV9jYWlkJTNkMTk2NjA4ODM0MTYlMjZkc19hX2FnaWQlM2QxNDg4MzgwODcyMzElMjZkc19hX2xpZCUzZGt3ZC0yNTk5Mzc5NiUyNiUyNmRzX2VfYWRpZCUzZDgyOTQ0NzcyODI3MjU0JTI2ZHNfZV90YXJnZXRfaWQlM2Rrd2QtODI5NDU0NDM0MzQwMTAlM2Fsb2MtMTc1JTI2JTI2ZHNfZV9uZXR3b3JrJTNkcyUyNmRzX3VybF92JTNkMiUyNmRzX2Rlc3RfdXJsJTNkaHR0cHMlM2ElMmYlMmZ3d3cuZnVzdC5jaCUyZmRlJTJmciUyZnBjLXRhYmxldC1oYW5keSUyZnRhYmxldCUyZmFwcGxlLWlwYWQtNDU1Lmh0bWwlM2ZnY2xpZCUzZDViYmFmZTNmYTk5MzE3M2RmYjZjOTg1YzkzNWQ0NTQ2JTI2Z2Nsc3JjJTNkM3AuZHMlMjYlMjZtc2Nsa2lkJTNkNWJiYWZlM2ZhOTkzMTczZGZiNmM5ODVjOTM1ZDQ1NDYlMjZ1dG1fc291cmNlJTNkYmluZyUyNnV0bV9tZWRpdW0lM2RjcGMlMjZ1dG1fY2FtcGFpZ24lM2RCaW5nJTI1MjAtJTI1MjBOQnJhbmQlMjUyMC0lMjUyMFMlMjUyMC0lMjUyMEQlMjUyMC0lMjUyME1NJTI1MjBQQyUyNTIwTWFya2UlMjUyMEFwcGxlJTI2dXRtX3Rlcm0lM2RpcGFkJTI2dXRtX2NvbnRlbnQlM2QxX0FwcGxlJTI1M0QyX0FwcGxlJTI1MjBpUGFkJWMyJWE2M19OYnJhbmQ&rlid=5bbafe3fa993173dfb6c985c935d4546",
      "https://ad.doubleclick.net/searchads/link/click?lid=43700076191444430&ds_s_kwgid=58700008270566934&ds_a_cid=411205395&ds_a_caid=19660883416&ds_a_agid=148838087231&ds_a_lid=kwd-25993796&&ds_e_adid=82944772827254&ds_e_target_id=kwd-82945443434010:loc-175&&ds_e_network=s&ds_url_v=2&ds_dest_url=https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?gclid=5bbafe3fa993173dfb6c985c935d4546&gclsrc=3p.ds&&msclkid=5bbafe3fa993173dfb6c985c935d4546&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand",
      "https://www.fust.ch/de/r/pc-tablet-handy/tablet/apple-ipad-455.html?&msclkid=5bbafe3fa993173dfb6c985c935d4546&utm_source=bing&utm_medium=cpc&utm_campaign=Bing%20-%20NBrand%20-%20S%20-%20D%20-%20MM%20PC%20Marke%20Apple&utm_term=ipad&utm_content=1_Apple%3D2_Apple%20iPad%C2%A63_Nbrand&gclid=5bbafe3fa993173dfb6c985c935d4546&gclsrc=3p.ds"
    ],
    "time": "2024-06-07T18:52:35.825523+02:00"
  }
]
```

## 3rd party libraries

- Rod: [GitHub repo](https://github.com/go-rod/rod), [documentation](https://go-rod.github.io/)
- Shoutrrr: [GitHub repo](https://github.com/containrrr/shoutrrr), [documentation](https://containrrr.dev/shoutrrr/v0.8/)
- Fatih/color: [GitHub repo](https://github.com/fatih/color), [Go reference](https://pkg.go.dev/github.com/fatih/color)
- Carlmjohnson/requests [GitHub repo](https://github.com/carlmjohnson/requests), [Go reference](https://pkg.go.dev/github.com/carlmjohnson/requests)
