# ios-proxy-rules

Auto rule generator that converts `v2fly` source lists into iOS rule formats.

## Platforms

- Shadowrocket (`.list`)
- Loon (`.list`)
- Surge (`.list`)
- Clash / Stash (`.yaml`)

## How to Use

> [!IMPORTANT]  
> You don't need to run the generator yourself. Use these auto-updated direct links for your proxy client.

#### Shadowrocket
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geosite/shadowrocket/youtube.list`
*   **US IPs:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geoip/shadowrocket/us.list`

#### Loon
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geosite/loon/youtube.list`
*   **US IPs:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geoip/loon/us.list`

#### Surge
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geosite/surge/youtube.list`
*   **US IPs:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geoip/surge/us.list`

#### Clash / Stash
Add these URLs to your `rule-providers`:

*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geosite/clash/youtube.yaml`
*   **US IPs:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/geoip/clash/us.yaml`

> [!TIP]
> Find more rules in [/geosite](./geosite) (domains) and [/geoip](./geoip) (IP ranges) folders.

## Run

```bash
go run ./cmd/rulegen -geosite dlc.dat -geoip geoip.dat
```
