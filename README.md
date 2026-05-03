# ios-proxy-rules

Auto rule generator that converts `v2fly/domain-list-community` source lists into iOS rule formats.

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

*   **Facebook:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/shadowrocket/facebook.list`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/shadowrocket/youtube.list`

#### Loon
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **Facebook:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/loon/facebook.list`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/loon/youtube.list`

#### Surge
Copy the URL and add it as a **Remote Rule** (Rule List):

*   **Facebook:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/surge/facebook.list`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/surge/youtube.list`

#### Clash / Stash
Add these URLs to your `rule-providers`:

*   **Facebook:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/clash/facebook.yaml`
*   **YouTube:** `https://raw.githubusercontent.com/ddranic/ios-proxy-rules/main/rules/clash/youtube.yaml`

> [!TIP]
> Find more rules for other services in [/rules](./rules) folder.

## Run

Use defaults (`-input data`, `-output rules`):

```bash
go run ./cmd/rulegen
```
