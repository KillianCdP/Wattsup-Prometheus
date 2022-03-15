package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"go.bug.st/serial"
)

type Request struct {
	C1                byte
	C2                byte
	NumberOfArguments uint8
	Arguments         []byte
}

func NewRequest(c1 byte, c2 byte, n uint8, args ...byte) Request {
	return Request{
		C1:                c1,
		C2:                c2,
		NumberOfArguments: n,
		Arguments:         args,
	}
}

type Data struct {
	W       int // Watts
	V       int // Volts
	A       int // Amps
	WH      int // W/h
	Cost    int // Cost
	WHPMo   int // W/h per month
	CostPMo int // Cost per month
	Wmax    int // Max watts
	Vmax    int // max volts
	Amax    int // max amps
	Wmin    int // min watts
	Vmin    int // min volts
	Amin    int // min amps
	PF      int // power factor
	DC      int // Duty cycle
	PC      int // Power cycle
	Hz      int // Line frequency
	VA      int // volt amps
}

// #d,-,18,972,2277,433,5,0,70035,5602,973,2280,435,972,2274,432,95,0,0,500,1020;
func fromString(s string) Data {
	offset := 3
	fields := strings.Split(s, ",")
	var d [18]int
	var err error
	for i := offset; i < 18; i++ {
		d[i-offset], err = strconv.Atoi(fields[i])
		if err != nil {
			log.Printf("Could not parse field correctly : \n\t%s", fields[i])
		}
	}

	return Data{
		W:       d[0] / 10,
		V:       d[1] / 10,
		A:       d[2] / 1000,
		WH:      d[3],
		Cost:    d[4],
		WHPMo:   d[5],
		CostPMo: d[6],
		Wmax:    d[7],
		Vmax:    d[8],
		Amax:    d[9],
		Wmin:    d[10],
		Vmin:    d[11],
		Amin:    d[12],
		PF:      d[13],
		DC:      d[14],
		PC:      d[15],
		Hz:      d[16] / 10,
		VA:      d[17],
	}
}

type WUPPort struct {
	S serial.Port
}

func (p WUPPort) Write(request []byte) (int, error) {
	return p.S.Write(request)
}

func (p WUPPort) Read() (Data, error) {
	// Read and print the response
	buff := make([]byte, 1000)
	var d Data
	// Reads up to 100 bytes
	n, err := p.S.Read(buff)
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		fmt.Println("\nEOF")
		return Data{}, errors.New("no data read")
	}
	d = fromString(string(buff[:n]))
	return d, nil
}

func (p *WUPPort) Reset() {
	_, err := p.Write([]byte("#C,R,0;\n\r"))
	if err != nil {
		log.Fatal("Could not reset serial port")
		return
	}
	log.Println("Succesfully reset the port")
}
