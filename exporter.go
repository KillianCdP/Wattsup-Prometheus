package main

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"go.bug.st/serial"
)

var (
	watts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "watts"),
		"Watts",
		[]string{"device"}, nil,
	)
	volts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volts"),
		"Volts",
		[]string{"device"}, nil,
	)
	amps = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amps"),
		"Amps",
		[]string{"device"}, nil,
	)
	wh = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wh"),
		"W/h",
		[]string{"device"}, nil,
	)
	cost = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cost"),
		"Cost",
		[]string{"device"}, nil,
	)
	whpmo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "whpmo"),
		"Wh / month",
		[]string{"device"}, nil,
	)
	costpmo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "costpmo"),
		"Cost per month",
		[]string{"device"}, nil,
	)
	wmax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wmax"),
		"Max watts",
		[]string{"device"}, nil,
	)
	vmax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "vmax"),
		"Max volts",
		[]string{"device"}, nil,
	)
	amax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amax"),
		"Max amps",
		[]string{"device"}, nil,
	)
	wmin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wmin"),
		"Min watts",
		[]string{"device"}, nil,
	)
	vmin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "vmin"),
		"Min volts",
		[]string{"device"}, nil,
	)
	amin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amin"),
		"Min amps",
		[]string{"device"}, nil,
	)
	pf = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "pf"),
		"Power factor",
		[]string{"device"}, nil,
	)
	dc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dc"),
		"Duty Cycle",
		[]string{"device"}, nil,
	)
	pc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "pc"),
		"Power cycle",
		[]string{"device"}, nil,
	)
	hz = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "hz"),
		"Line frequency",
		[]string{"device"}, nil,
	)
	va = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "va"),
		"Volt amps",
		[]string{"device"}, nil,
	)
)

type Exporter struct {
	port WUPPort
}

func NewExporter(loggingInterval int) *Exporter {
	port, err := serial.Open(serialPort, serialMode)
	if err != nil {
		log.Fatal("Could not open serial port: ", err)
	}
	wu := WUPPort{S: port}
	wu.LoggingStop()
	wu.LoggingStart("internal", loggingInterval)
	return &Exporter{
		port: wu,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- watts
	ch <- volts
	ch <- amps
	ch <- wh
	ch <- cost
	ch <- whpmo
	ch <- costpmo
	ch <- wmax
	ch <- vmax
	ch <- amax
	ch <- wmin
	ch <- vmin
	ch <- amin
	ch <- pf
	ch <- dc
	ch <- pc
	ch <- hz
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	data, err := e.port.ReadData()
	if err != nil {
		log.Println("Error while collecting data", err)
		e.port.Reset()
		return
	}
	ch <- prometheus.MustNewConstMetric(watts, prometheus.GaugeValue, float64(data.Watts), deviceName)
	ch <- prometheus.MustNewConstMetric(volts, prometheus.GaugeValue, float64(data.Volts), deviceName)
	ch <- prometheus.MustNewConstMetric(amps, prometheus.GaugeValue, float64(data.Amps), deviceName)
	ch <- prometheus.MustNewConstMetric(wh, prometheus.GaugeValue, float64(data.WH), deviceName)
	ch <- prometheus.MustNewConstMetric(cost, prometheus.GaugeValue, float64(data.Cost), deviceName)
	ch <- prometheus.MustNewConstMetric(whpmo, prometheus.GaugeValue, float64(data.WHPMo), deviceName)
	ch <- prometheus.MustNewConstMetric(costpmo, prometheus.GaugeValue, float64(data.CostPMo), deviceName)
	ch <- prometheus.MustNewConstMetric(wmax, prometheus.GaugeValue, float64(data.Wmax), deviceName)
	ch <- prometheus.MustNewConstMetric(vmax, prometheus.GaugeValue, float64(data.Vmax), deviceName)
	ch <- prometheus.MustNewConstMetric(amax, prometheus.GaugeValue, float64(data.Amax), deviceName)
	ch <- prometheus.MustNewConstMetric(wmin, prometheus.GaugeValue, float64(data.Wmin), deviceName)
	ch <- prometheus.MustNewConstMetric(vmin, prometheus.GaugeValue, float64(data.Vmin), deviceName)
	ch <- prometheus.MustNewConstMetric(amin, prometheus.GaugeValue, float64(data.Amin), deviceName)
	ch <- prometheus.MustNewConstMetric(pf, prometheus.GaugeValue, float64(data.PF), deviceName)
	ch <- prometheus.MustNewConstMetric(dc, prometheus.GaugeValue, float64(data.DC), deviceName)
	ch <- prometheus.MustNewConstMetric(pc, prometheus.GaugeValue, float64(data.PC), deviceName)
	ch <- prometheus.MustNewConstMetric(hz, prometheus.GaugeValue, float64(data.Hz), deviceName)
	ch <- prometheus.MustNewConstMetric(va, prometheus.GaugeValue, float64(data.VA), deviceName)
}
