package generator

type Platform string

const (
	PlatformShadowrocket Platform = "shadowrocket"
	PlatformLoon         Platform = "loon"
	PlatformSurge        Platform = "surge"
	PlatformStash        Platform = "stash"
	PlatformClash        Platform = "clash"
)

var AllPlatforms = []Platform{
	PlatformShadowrocket,
	PlatformLoon,
	PlatformSurge,
	PlatformStash,
	PlatformClash,
}

const (
	extList = ".list"
	extYAML = ".yaml"
)
