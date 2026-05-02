package generator

type Platform string

const (
	PlatformShadowrocket Platform = "shadowrocket"
	PlatformLoon         Platform = "loon"
	PlatformSurge        Platform = "surge"
	PlatformStash        Platform = "stash"
	PlatformClash        Platform = "clash"
)

var Platforms = []Platform{
	PlatformShadowrocket,
	PlatformLoon,
	PlatformSurge,
	PlatformStash,
	PlatformClash,
}

type Extension string

const (
	ExtensionList = ".list"
	ExtensionYaml = ".yaml"
)

var PlatformsExtensions = map[Platform]string{
	PlatformShadowrocket: ExtensionList,
	PlatformLoon:         ExtensionList,
	PlatformSurge:        ExtensionList,
	PlatformStash:        ExtensionList,
	PlatformClash:        ExtensionYaml,
}
