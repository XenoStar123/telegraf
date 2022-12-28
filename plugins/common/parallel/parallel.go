package parallel

import "github.com/XenoStar123/telegraf"

type Parallel interface {
	Enqueue(telegraf.Metric)
	Stop()
}
