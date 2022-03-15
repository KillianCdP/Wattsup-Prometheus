package main

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {

	serialMode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}

	for _, p := range ports {
		fmt.Printf("Found port: %s\n", p.Name)
		if p.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", p.VID, p.PID)
			fmt.Printf("   USB serial %s\n", p.SerialNumber)
		}
	}

	port, err := serial.Open("/dev/ttyUSB0", serialMode)
	if err != nil {
		log.Fatal("Could not open serial port: ", err)
	}
	wu := WUPPort{S: port}

	for {
		d, err := wu.Read()
		if err != nil {
			log.Println(err)
			wu.Reset()
		}
		fmt.Println(d)
		time.Sleep(2 * time.Second)
	}
	// Send the string "10,20,30\n\r" to the serial port
	//n, err := port.Write([]byte("#L,W,3,E,-,15;\n\r"))
	// Log all fields
	// n, err := port.Write([]byte("#C,W,18,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1;\n\r"))
	//n, err := port.Write([]byte("#C,R,18,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1;\n\r"))

}
