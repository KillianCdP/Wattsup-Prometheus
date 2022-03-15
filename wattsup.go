package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
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
	Watts   int // Watts
	Volts   int // Volts
	Amps    int // Amps
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

func fromString(s string) (Data, error) {
	offset := 3
	s = strings.TrimRight(s, "\n\r")
	if len(s) == 0 {
		return Data{}, errors.New("empty string")
	}
	if s[0] != '#' {
		return Data{}, errors.New("wrong start token")
	}
	if s[len(s)-1] != ';' {
		return Data{}, fmt.Errorf("wrong end token %v", s[len(s)-1:])
	}
	fields := strings.Split(s, ",")
	if len(fields) < 18+offset {
		return Data{}, errors.New("data too short")
	}
	if fields[2] != "18" {
		return Data{}, errors.New("not a valid data message")
	}
	// Trim ;
	fields[17+offset] = strings.TrimRight(fields[17+offset], ";")

	var d [18]int
	var err error
	for i := offset; i < 18; i++ {
		d[i-offset], err = strconv.Atoi(fields[i])
		if err != nil {
			return Data{}, fmt.Errorf("error while parsing data field %v", fields[i])
		}
	}
	fmt.Println("d", d)

	return Data{
		Watts:   d[0] / 10,
		Volts:   d[1] / 10,
		Amps:    d[2] / 1000,
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
	}, nil
}

type WUPPort struct {
	S serial.Port
}

func (p WUPPort) Write(request []byte) (int, error) {
	return p.S.Write(request)
}

func (p WUPPort) Read() (string, error) {
	var err error
	buff := make([]byte, 256)
	var line []byte

	for {
		n, err := p.S.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			return "", errors.New("no data read")
		}
		line = append(line, buff[:n]...)
		if strings.Contains(string(buff[:n]), "\n") {
			break
		}
	}
	return string(line), err
}

var lineRegex = regexp.MustCompile(".*(#[^;]+;).*")

func extractLine(data string) (string, error) {
	r := lineRegex.FindStringSubmatch(data)
	if len(r) > 0 {
		// Return last matching line
		return r[len(r)-1], nil
	}
	return "", errors.New("could not find a valid line")
}

func (p WUPPort) ReadData() (Data, error) {
	s, err := p.Read()
	if err != nil {
		fmt.Println("Could not read data:", err)
	}
	s, err = extractLine(s)
	if err != nil {
		return Data{}, err
	}
	d, err := fromString(s)
	return d, err
}

func (p WUPPort) ReadLastData() (Data, error) {
	_, err := p.Write([]byte("#D,R,0"))
	if err != nil {
		return Data{}, errors.New("could not request data")
	}
	d, err := p.ReadData()
	return d, err
}

func (p *WUPPort) Reset() {
	_, err := p.Write([]byte("#C,R,0;\n\r"))
	if err != nil {
		log.Fatal("Could not reset serial port")
		return
	}
	// Setup internal logging
	p.Write([]byte("#L,W,3,I,-,5;"))
	log.Println("Succesfully reset the port")
}

func (p *WUPPort) LoggingStart(mode string, interval int) {
	var err error
	switch mode {
	case "internal":
		_, err = p.Write([]byte(fmt.Sprintf("#L,W,3,I,-,%d;\n\r", interval)))
	case "external":
		_, err = p.Write([]byte(fmt.Sprintf("#L,W,3,E,-,%d;\n\r", interval)))
	default:
		log.Fatal("Wrong logging mode", mode)
		return
	}
	if err != nil {
		log.Println("Could not start logging", err)
		return
	}
	log.Println("Logging requested")
}

func (p *WUPPort) LoggingStop() {
	_, err := p.Write([]byte("#L,R,0;\n\r"))
	if err != nil {
		log.Println("Could not stop logging", err)
		return
	}
	log.Println("Logging stopped")
}
