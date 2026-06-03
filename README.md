# ios-proxy-rules

[![Update](https://github.com/ddranic/ios-proxy-rules/actions/workflows/update.yml/badge.svg)](https://github.com/ddranic/ios-proxy-rules/actions/workflows/update.yml)
[![Last Commit](https://img.shields.io/github/last-commit/ddranic/ios-proxy-rules)](https://github.com/ddranic/ios-proxy-rules/commits/main)

Auto rule generator that converts `v2fly` source lists into iOS rule formats.

## Platforms

- Shadowrocket (`.list`)
- Loon (`.list`)
- Surge (`.list`)
- Clash / Stash (`.yaml`)
- sing-box (`.srs`)

## How to Use

> [!IMPORTANT]  
> You don't need to run the generator yourself. Use these auto-updated direct links for your proxy platform.

Example for Shadowrocket:

*   **YouTube GEOSITE** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geosite/shadowrocket/youtube.list`
*   **China GEOIP:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geoip/shadowrocket/cn.list`

> [!TIP]
> Browse all rules: [/geosite](./geosite) (domains) · [/geoip](./geoip) (IP ranges)

## Run

```bash
go run ./cmd/rulegen -geosite dlc.dat -geoip geoip.dat
```
