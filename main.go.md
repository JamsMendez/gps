package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type Geolocation struct {
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

func main() {
	config := &serial.Config{Name: "COM5", Baud: 9600}
	serialPort, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{}

	go func() {
		line := []byte{}
		for {
			buffer := make([]byte, 64)
			n, err := serialPort.Read(buffer)
			if err != nil {
				fmt.Println("SerialPort.Read.ERROR: ", err)
			}

			chuck := buffer[:n]
			size := len(chuck)
			for j := 0; j < size; j++ {
				line = append(line, chuck[j])

				if chuck[j] == 10 {
					s := string(line)
					parts := strings.Split(s, "\r\n")
					if len(parts) > 0 {
						first := parts[0]

						lines = append(lines, first)
						line = []byte{}
					}
				}
			}
		}

	}()

	sentences := []string{}

	for {
		size := len(lines)
		for i := 0; i < size; i++ {
			line := lines[i]

			sentences = append(sentences, line)

			if strings.Contains(line, "$GPVTG") {
				next := i + 1
				if next < size {
					nLines := []string{}
					for j := next; j < size; j++ {
						nLines = append(nLines, lines[j])
					}
					//fmt.Println("PRE: ", len(lines))
					lines = nLines
					//fmt.Println("PRO: ", len(lines))
				}

				break
			}
		}

		size = len(sentences)
		if size > 0 {
			geo := parseNmeaSetence(sentences)
			fmt.Println(geo.Latitude, geo.Longitude, geo.Altitude, geo.Time, geo.Speed)

			sentences = []string{}
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func parseNmeaSetence(sentences []string) Geolocation {
	geolocation := Geolocation{}

	for _, sentence := range sentences {
		cksum := strings.Split(sentence, "*")
		size := len(cksum)
		if size >= 2 {
			if cksum[1] != getNmeaChecksum(cksum[0][1:]) {
				return geolocation
			}

			segments := strings.Split(cksum[0], ",")
			size := len(segments)
			if size == 0 {
				return geolocation
			}

			if segments[0] == "$GPGGA" {
				// Time, position and fix related data
				time := toFloat64(segments[1])
				latitude := degToDec(segments[2], segments[3], 2, 8)
				longitude := degToDec(segments[4], segments[5], 3, 8)
				altitude := toFloat64(segments[9])

				geolocation.Latitude = latitude
				geolocation.Longitude = longitude
				geolocation.Altitude = altitude
				geolocation.Time = time
			}

			if segments[0] == "$GPGSA" {
				// Operating details
				satellites := segments[3:15]
				pdop := segments[15]
				hdop := segments[16]
				vdop := segments[17]

				geolocation.Satellites = satellites
				geolocation.Pdop = pdop
				geolocation.Hdop = hdop
				geolocation.Vdop = vdop
			}

			if segments[0] == "$GPRMC" {
				// GPS & Transit data
				time := toFloat64(segments[1])
				latitude := degToDec(segments[3], segments[4], 2, 8)
				longitude := degToDec(segments[5], segments[6], 3, 8)
				course := toFloat64(segments[8])
				speed := toFloat64(segments[7])

				geolocation.Longitude = longitude
				geolocation.Latitude = latitude
				geolocation.Course = course
				geolocation.Speed = speed
				geolocation.Time = time
			}

			if segments[0] == "$GPVTG" {
				// Track Mode Good and Ground Speed
				course := toFloat64(segments[1])
				speed := toFloat64(segments[5])

				geolocation.Course = course
				geolocation.Speed = speed
			}
		}
	}

	return geolocation
}

func getNmeaChecksum(value string) string {
	var cksum uint8
	size := len(value)
	for i := 0; i < size; i++ {
		cksum ^= byte(value[i])
	}

	buffer := []byte{cksum}
	s := hex.EncodeToString(buffer)

	nCksum := strings.ToUpper(s)
	size = len(nCksum)
	if size < 2 {
		nCksum = fmt.Sprintf("00%s", s)
		nCksum = nCksum[len(nCksum)-2:]
	}

	return nCksum
}

func toFloat64(value string) float64 {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}

	return f64
}

/*func toSpeed(value string, fixed int) float64 {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}

	f64 = f64 * 0.514444

	sFixed := fmt.Sprintf("%d", fixed)
	format := "%." + sFixed + "f"
	s := fmt.Sprintf(format, f64)
	f64, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return f64
}*/

func degToDec(degrees, cardinal string, intDigitsLength, fixed int) float64 {
	if degrees != "" {
		first := degrees[0:intDigitsLength]
		last := degrees[intDigitsLength:]

		f64, err := strconv.ParseFloat(first, 64)
		if err != nil {
			return 0.0
		}

		l64, err := strconv.ParseFloat(last, 64)
		if err != nil {
			return 0.0
		}

		decimal := f64 + (l64 / 60)

		if cardinal == "S" || cardinal == "W" {
			decimal = decimal * -1
		}

		sFixed := fmt.Sprintf("%d", fixed)
		format := "%." + sFixed + "f"
		s := fmt.Sprintf(format, decimal)
		f64, err = strconv.ParseFloat(s, 64)
		if err != nil {
			return 0.0
		}

		return f64
	}

	return 0.0
}
