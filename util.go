package gps

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

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

func toSpeed(value string, fixed int) float64 {
	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}

	f64 = f64 * 0.514444 // knots to m/s

	sFixed := fmt.Sprintf("%d", fixed)
	format := "%." + sFixed + "f"
	s := fmt.Sprintf(format, f64)
	f64, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}

	return f64
}

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
