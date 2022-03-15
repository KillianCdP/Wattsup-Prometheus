package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.bug.st/serial"
)

const interval int = 3

const namespace = "WUPExporter"

var serialMode = &serial.Mode{
	BaudRate: 115200,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.OneStopBit,
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	exporter := NewExporter(interval)
	prometheus.MustRegister(exporter)
	log.Fatal(http.ListenAndServe(":9091", nil))
}
