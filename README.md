# ios-proxy-rules

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

*   **YouTube:** `https://cdn.jsdelivr.net/gh/ddranic/ios-proxy-rules@main/geosite/shadowrocket/youtube.list`
*   **China IPs:** `https://cdn.jsdelivr.net/gh/ddranic/ios-proxy-rules@main/geoip/shadowrocket/cn.list`

> [!TIP]
> Browse all rules: [/geosite](./geosite) (domains) · [/geoip](./geoip) (IP ranges)

## Run

```bash
go run ./cmd/rulegen -geosite dlc.dat -geoip geoip.dat
```
