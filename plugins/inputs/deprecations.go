package inputs

import "github.com/XenoStar123/telegraf"

// Deprecations lists the deprecated plugins
var Deprecations = map[string]telegraf.DeprecationInfo{
	"cassandra": {
		Since:  "1.7.0",
		Notice: "use 'inputs.jolokia2' with the 'cassandra.conf' example configuration instead",
	},
	"cisco_telemetry_gnmi": {
		Since:  "1.15.0",
		Notice: "has been renamed to 'gnmi'",
	},
	"http_listener": {
		Since:  "1.9.0",
		Notice: "has been renamed to 'influxdb_listener', use 'inputs.influxdb_listener' or 'inputs.http_listener_v2' instead",
	},
	"httpjson": {
		Since:  "1.6.0",
		Notice: "use 'inputs.http' instead",
	},
	"io": {
		Since:     "0.10.0",
		RemovalIn: "2.0.0",
		Notice:    "use 'inputs.diskio' instead",
	},
	"jolokia": {
		Since:  "1.5.0",
		Notice: "use 'inputs.jolokia2' instead",
	},
	"kafka_consumer_legacy": {
		Since:  "1.4.0",
		Notice: "use 'inputs.kafka_consumer' instead, NOTE: 'kafka_consumer' only supports Kafka v0.8+",
	},
	"KNXListener": {
		Since:  "1.20.1",
		Notice: "has been renamed to 'knx_listener'",
	},
	"logparser": {
		Since:  "1.15.0",
		Notice: "use 'inputs.tail' with 'grok' data format instead",
	},
	"snmp_legacy": {
		Since:  "1.0.0",
		Notice: "use 'inputs.snmp' instead",
	},
	"tcp_listener": {
		Since:  "1.3.0",
		Notice: "use 'inputs.socket_listener' instead",
	},
	"udp_listener": {
		Since:  "1.3.0",
		Notice: "use 'inputs.socket_listener' instead",
	},
}
