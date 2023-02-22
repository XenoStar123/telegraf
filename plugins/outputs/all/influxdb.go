//go:build !custom || outputs || outputs.influxdb

package all

import _ "github.com/XenoStar123/telegraf/plugins/outputs/influxdb" // register plugin
