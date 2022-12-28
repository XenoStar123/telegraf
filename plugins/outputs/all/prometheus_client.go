//go:build !custom || outputs || outputs.prometheus_client

package all

import _ "github.com/XenoStar123/telegraf/plugins/outputs/prometheus_client" // register plugin
