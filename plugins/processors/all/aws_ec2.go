//go:build !custom || processors || processors.aws_ec2

package all

import _ "github.com/XenoStar123/telegraf/plugins/processors/aws/ec2" // register plugin
