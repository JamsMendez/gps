package gps

import (
	"bufio"
	"errors"
	"log"

	"github.com/tarm/serial"
)

type serialPort interface {
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Close() error
}

type GPS struct {
	latitude   float64
	longitude  float64
	altitude   float64
	course     float64
	speed      float64
	time       float64
	satellites []string
	pdop       string
	hdop       string
	vdop       string

	serialPort  serialPort
	data        chan string
	isConnected bool

	Port     string
	BaudRate int
	Debug    bool
}

type Position struct {
	Latitude   float64  `json:"latitude"`
	Longitude  float64  `json:"longitude"`
	Altitude   float64  `json:"altitude"`
	Course     float64  `json:"course"`
	Speed      float64  `json:"speed"`
	Time       float64  `json:"time"`
	Satellites []string `json:"satellites"`
	Pdop       string   `json:"pdop"`
	Hdop       string   `json:"hdop"`
	Vdop       string   `json:"vdop"`
}

// Reading reads from gps
func (gps *GPS) Reading() {
	scanner := bufio.NewScanner(gps.serialPort)
	for scanner.Scan() {
		gps.parseNmeaSetence(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if gps.Debug {
			log.Println("GPS.SerialPort.Reading.ERROR: ", err)
		}

		gps.Disconnect()
	}
}

func (gps *GPS) isClosed(ch chan string) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

// Connect open serial connection with gps
func (gps *GPS) Connect() error {
	config := &serial.Config{Name: gps.Port, Baud: gps.BaudRate}
	var err error
	gps.serialPort, err = serial.OpenPort(config)
	if err != nil {
		if gps.Debug {
			log.Println("GPS.SerialPort.OpenPort.ERROR: ", err)
		}

		return err
	}

	gps.data = make(chan string)
	gps.isConnected = true

	return err
}

// Disconnect ... close serial connection with gps
func (gps *GPS) Disconnect() error {
	if gps.isConnected {
		gps.isConnected = false
	}

	if !gps.isClosed(gps.data) {
		if gps.data != nil {
			close(gps.data)
		}
	}

	var err error

	if gps.serialPort != nil {
		err = gps.serialPort.Close()
		if err != nil {
			if gps.Debug {
				log.Println("GPS.SerialPort.Close.ERROR: ", err)
			}

			return err
		}
	}

	return err
}

// IsConnected return state of the connection with gps
func (gps *GPS) IsConnected() bool {
	return gps.isConnected
}

// FetchPosition return current gps position
func (gps *GPS) FetchPosition() (Position, error) {
	var err error

	position := Position{
		Latitude:   gps.latitude,
		Longitude:  gps.longitude,
		Altitude:   gps.altitude,
		Course:     gps.course,
		Speed:      gps.speed,
		Time:       gps.time,
		Satellites: gps.satellites,
		Pdop:       gps.pdop,
		Hdop:       gps.hdop,
		Vdop:       gps.vdop,
	}

	if !gps.isConnected {
		err = errors.New("GPS is disconnected")

		return position, err
	}

	return position, err
}
