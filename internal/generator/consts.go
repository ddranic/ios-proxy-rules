package generator

type Platform string

const (
	PlatformShadowrocket Platform = "shadowrocket"
	PlatformLoon         Platform = "loon"
	PlatformSurge        Platform = "surge"
	PlatformClash        Platform = "clash"
)

var Platforms = []Platform{
	PlatformShadowrocket,
	PlatformLoon,
	PlatformSurge,
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
	PlatformClash:        ExtensionYaml,
}
