package generator

type Platform string

const (
	PlatformShadowrocket Platform = "shadowrocket"
	PlatformLoon         Platform = "loon"
	PlatformSurge        Platform = "surge"
	PlatformClash        Platform = "clash"
	PlatformSingbox      Platform = "singbox"
)

var Platforms = []Platform{
	PlatformShadowrocket,
	PlatformLoon,
	PlatformSurge,
	PlatformClash,
	PlatformSingbox,
}

type Extension string

const (
	ExtensionList = ".list"
	ExtensionYaml = ".yaml"
	ExtensionJSON = ".json"
)

var PlatformsExtensions = map[Platform]string{
	PlatformShadowrocket: ExtensionList,
	PlatformLoon:         ExtensionList,
	PlatformSurge:        ExtensionList,
	PlatformClash:        ExtensionYaml,
	PlatformSingbox:      ExtensionJSON,
}
