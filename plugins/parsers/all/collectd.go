//go:build !custom || parsers || parsers.collectd

package all

import _ "github.com/XenoStar123/telegraf/plugins/parsers/collectd" // register plugin
