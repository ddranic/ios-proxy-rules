# ios-proxy-rules

Auto rule generator that converts `v2fly/domain-list-community` source lists into iOS/cross-client rule formats.

## Platforms

- Shadowrocket / Loon / Surge (`.list`)
- Clash / Stash (`.yaml`)

## How to Use

You don't need to run the generator yourself. Use these auto-updated direct links for your proxy client:

### Shadowrocket / Surge / Loon
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **TikTok:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/tiktok.list`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/youtube.list`

### Clash / Stash
Add these URLs to your `rule-providers`:

*   **TikTok:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/tiktok.yaml`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/youtube.yaml`

> [!TIP]
> Find more rules for other services in [/rules](./rules) folder.

## Run

Use defaults (`-input data`, `-output rules`):

```bash
go run ./cmd/rulegen
```
