package splunkmetric

import (
	"encoding/json"
	"log"

	"github.com/XenoStar123/telegraf"
)

type serializer struct {
	HecRouting              bool
	SplunkmetricMultiMetric bool
	OmitEventTag            bool
}

type CommonTags struct {
	Time   float64
	Host   string
	Index  string
	Source string
	Fields map[string]interface{}
}

type HECTimeSeries struct {
	Time   float64                `json:"time"`
	Event  string                 `json:"event,omitempty"`
	Host   string                 `json:"host,omitempty"`
	Index  string                 `json:"index,omitempty"`
	Source string                 `json:"source,omitempty"`
	Fields map[string]interface{} `json:"fields"`
}

// NewSerializer Setup our new serializer
func NewSerializer(splunkmetricHecRouting bool, splunkmetricMultimetric bool, splunkmetricOmitEventTag bool) *serializer {
	/*	Define output params */
	s := &serializer{
		HecRouting:              splunkmetricHecRouting,
		SplunkmetricMultiMetric: splunkmetricMultimetric,
		OmitEventTag:            splunkmetricOmitEventTag,
	}
	return s
}

func (s *serializer) Serialize(metric telegraf.Metric) ([]byte, error) {
	return s.createObject(metric), nil
}

func (s *serializer) SerializeBatch(metrics []telegraf.Metric) ([]byte, error) {
	var serialized []byte

	for _, metric := range metrics {
		m := s.createObject(metric)
		if m != nil {
			serialized = append(serialized, m...)
		}
	}

	return serialized, nil
}

func (s *serializer) createMulti(metric telegraf.Metric, dataGroup HECTimeSeries, commonTags CommonTags) (metricGroup []byte, err error) {
	/* When splunkmetric_multimetric is true, then we can write out multiple name=value pairs as part of the same
	** event payload. This only works when the time, host, and dimensions are the same for every name=value pair
	** in the timeseries data.
	**
	** The format for multimetric data is 'metric_name:nameOfMetric = valueOfMetric'
	 */
	var metricJSON []byte

	// Set the event data from the commonTags above.
	if !s.OmitEventTag {
		dataGroup.Event = "metric"
	}
	dataGroup.Time = commonTags.Time
	dataGroup.Host = commonTags.Host
	dataGroup.Index = commonTags.Index
	dataGroup.Source = commonTags.Source
	dataGroup.Fields = commonTags.Fields

	// Stuff the metric data into the structure.
	for _, field := range metric.FieldList() {
		value, valid := verifyValue(field.Value)

		if !valid {
			log.Printf("D! Can not parse value: %v for key: %v", field.Value, field.Key)
			continue
		}

		dataGroup.Fields["metric_name:"+metric.Name()+"."+field.Key] = value
	}

	// Manage the rest of the event details based upon HEC routing rules
	switch s.HecRouting {
	case true:
		// Output the data as a fields array and host,index,time,source overrides for the HEC.
		metricJSON, err = json.Marshal(dataGroup)
	default:
		// Just output the data and the time, useful for file based outputs
		dataGroup.Fields["time"] = dataGroup.Time
		metricJSON, err = json.Marshal(dataGroup.Fields)
	}
	if err != nil {
		return nil, err
	}
	// Let the JSON fall through to the return below
	metricGroup = metricJSON

	return metricGroup, nil
}

func (s *serializer) createSingle(metric telegraf.Metric, dataGroup HECTimeSeries, commonTags CommonTags) (metricGroup []byte, err error) {
	/* The default mode is to generate one JSON entity per metric (required for pre-8.0 Splunks)
	**
	** The format for single metric is 'nameOfMetric = valueOfMetric'
	 */

	var metricJSON []byte

	for _, field := range metric.FieldList() {
		value, valid := verifyValue(field.Value)

		if !valid {
			log.Printf("D! Can not parse value: %v for key: %v", field.Value, field.Key)
			continue
		}

		dataGroup.Time = commonTags.Time

		// Apply the common tags from above to every record.
		if !s.OmitEventTag {
			dataGroup.Event = "metric"
		}
		dataGroup.Host = commonTags.Host
		dataGroup.Index = commonTags.Index
		dataGroup.Source = commonTags.Source
		dataGroup.Fields = commonTags.Fields

		dataGroup.Fields["metric_name"] = metric.Name() + "." + field.Key
		dataGroup.Fields["_value"] = value

		switch s.HecRouting {
		case true:
			// Output the data as a fields array and host,index,time,source overrides for the HEC.
			metricJSON, err = json.Marshal(dataGroup)
		default:
			// Just output the data and the time, useful for file based outputs
			dataGroup.Fields["time"] = dataGroup.Time
			metricJSON, err = json.Marshal(dataGroup.Fields)
		}

		metricGroup = append(metricGroup, metricJSON...)

		if err != nil {
			return nil, err
		}
	}

	return metricGroup, nil
}

func (s *serializer) createObject(metric telegraf.Metric) (metricGroup []byte) {
	/*  Splunk supports one metric json object, and does _not_ support an array of JSON objects.
	     ** Splunk has the following required names for the metric store:
		 ** metric_name: The name of the metric
		 ** _value:      The value for the metric
		 ** time:       The timestamp for the metric
		 ** All other index fields become dimensions.
	*/

	dataGroup := HECTimeSeries{}

	// The tags are common to all events in this timeseries
	commonTags := CommonTags{}

	commonTags.Fields = map[string]interface{}{}

	// Break tags out into key(n)=value(t) pairs
	for n, t := range metric.Tags() {
		if n == "host" {
			commonTags.Host = t
		} else if n == "index" {
			commonTags.Index = t
		} else if n == "source" {
			commonTags.Source = t
		} else {
			commonTags.Fields[n] = t
		}
	}
	commonTags.Time = float64(metric.Time().UnixNano()) / float64(1000000000)
	switch s.SplunkmetricMultiMetric {
	case true:
		metricGroup, _ = s.createMulti(metric, dataGroup, commonTags)
	default:
		metricGroup, _ = s.createSingle(metric, dataGroup, commonTags)
	}

	// Return the metric group regardless of if it's multimetric or single metric.
	return metricGroup
}

func verifyValue(v interface{}) (value interface{}, valid bool) {
	switch v.(type) {
	case string:
		valid = false
		value = v
	case bool:
		if v == bool(true) {
			// Store 1 for a "true" value
			valid = true
			value = 1
		} else {
			// Otherwise store 0
			valid = true
			value = 0
		}
	default:
		valid = true
		value = v
	}
	return value, valid
}
