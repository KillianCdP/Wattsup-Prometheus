package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.bug.st/serial"
	"log"
	"net/http"
)

const interval int = 3

const namespace = "WUPExporter"

var serialMode = &serial.Mode{
	BaudRate: 115200,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.OneStopBit,
}

var deviceName string
var serialPort string

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	flag.StringVar(&serialPort, "serialPort", ports[0], "Serial device path")
	flag.StringVar(&deviceName, "deviceName", "unknown", "Monitored device name")
	flag.Parse()
	println("Serial path :", serialPort)
	http.Handle("/metrics", promhttp.Handler())
	exporter := NewExporter(interval)
	prometheus.MustRegister(exporter)
	log.Fatal(http.ListenAndServe(":9091", nil))
}
