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
		nil, nil,
	)
	volts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "volts"),
		"Volts",
		nil, nil,
	)
	amps = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amps"),
		"Amps",
		nil, nil,
	)
	wh = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wh"),
		"W/h",
		nil, nil,
	)
	cost = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cost"),
		"Cost",
		nil, nil,
	)
	whpmo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "whpmo"),
		"Wh / month",
		nil, nil,
	)
	costpmo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "costpmo"),
		"Cost per month",
		nil, nil,
	)
	wmax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wmax"),
		"Max watts",
		nil, nil,
	)
	vmax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "vmax"),
		"Max volts",
		nil, nil,
	)
	amax = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amax"),
		"Max amps",
		nil, nil,
	)
	wmin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "wmin"),
		"Min watts",
		nil, nil,
	)
	vmin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "vmin"),
		"Min volts",
		nil, nil,
	)
	amin = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "amin"),
		"Min amps",
		nil, nil,
	)
	pf = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "pf"),
		"Power factor",
		nil, nil,
	)
	dc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dc"),
		"Duty Cycle",
		nil, nil,
	)
	pc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "pc"),
		"Power cycle",
		nil, nil,
	)
	hz = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "hz"),
		"Line frequency",
		nil, nil,
	)
	va = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "va"),
		"Volt amps",
		nil, nil,
	)
)

type Exporter struct {
	port WUPPort
}

func NewExporter(loggingInterval int) *Exporter {
	port, err := serial.Open("/dev/ttyUSB0", serialMode)
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
	ch <- prometheus.MustNewConstMetric(watts, prometheus.GaugeValue, float64(data.Watts))
	ch <- prometheus.MustNewConstMetric(volts, prometheus.GaugeValue, float64(data.Volts))
	ch <- prometheus.MustNewConstMetric(amps, prometheus.GaugeValue, float64(data.Amps))
	ch <- prometheus.MustNewConstMetric(wh, prometheus.GaugeValue, float64(data.WH))
	ch <- prometheus.MustNewConstMetric(cost, prometheus.GaugeValue, float64(data.Cost))
	ch <- prometheus.MustNewConstMetric(whpmo, prometheus.GaugeValue, float64(data.WHPMo))
	ch <- prometheus.MustNewConstMetric(costpmo, prometheus.GaugeValue, float64(data.CostPMo))
	ch <- prometheus.MustNewConstMetric(wmax, prometheus.GaugeValue, float64(data.Wmax))
	ch <- prometheus.MustNewConstMetric(vmax, prometheus.GaugeValue, float64(data.Vmax))
	ch <- prometheus.MustNewConstMetric(amax, prometheus.GaugeValue, float64(data.Amax))
	ch <- prometheus.MustNewConstMetric(wmin, prometheus.GaugeValue, float64(data.Wmin))
	ch <- prometheus.MustNewConstMetric(vmin, prometheus.GaugeValue, float64(data.Vmin))
	ch <- prometheus.MustNewConstMetric(amin, prometheus.GaugeValue, float64(data.Amin))
	ch <- prometheus.MustNewConstMetric(pf, prometheus.GaugeValue, float64(data.PF))
	ch <- prometheus.MustNewConstMetric(dc, prometheus.GaugeValue, float64(data.DC))
	ch <- prometheus.MustNewConstMetric(pc, prometheus.GaugeValue, float64(data.PC))
	ch <- prometheus.MustNewConstMetric(hz, prometheus.GaugeValue, float64(data.Hz))
	ch <- prometheus.MustNewConstMetric(va, prometheus.GaugeValue, float64(data.VA))
}
