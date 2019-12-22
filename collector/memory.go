package collector

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/mem"
)

const (
	subsystem = "memory"
)

type memoryCollector struct{}

func init() {
	registCollector(subsystem, NewMemoryCollector)
}

func NewMemoryCollector() (Collector, error) {
	return &memoryCollector{}, nil
}

func (c *memoryCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType

	v, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("could not get memory info: %s", err)
	}

	s := reflect.Indirect(reflect.ValueOf(v))
	t := s.Type()

	for i := 0; i < t.NumField(); i++ {
		f := s.Field(i)
		ff := f.Interface()
		if strings.Contains(t.Field(i).Name, "Total") == true {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		var f64 float64
		if _, ok := ff.(float64); ok {
			f64 = ff.(float64)
		} else {
			f64 = float64(ff.(uint64))
		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, subsystem, t.Field(i).Name),
				fmt.Sprintf("Memory information filed %s", t.Field(i).Name),
				nil, nil,
			),
			metricType, f64,
		)
	}
	return nil
}
