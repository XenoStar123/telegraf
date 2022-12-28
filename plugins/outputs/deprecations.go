package outputs

import "github.com/XenoStar123/telegraf"

// Deprecations lists the deprecated plugins
var Deprecations = map[string]telegraf.DeprecationInfo{
	"riemann_legacy": {
		Since:  "1.3.0",
		Notice: "use 'outputs.riemann' instead (see https://github.com/XenoStar123/telegraf/issues/1878)",
	},
}
