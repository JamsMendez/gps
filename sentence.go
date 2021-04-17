package gps

import "strings"

func (gps *GPS) parseNmeaSetence(sentence string) {
	cksum := strings.Split(sentence, "*")
	size := len(cksum)
	if size >= 2 {
		if cksum[1] != getNmeaChecksum(cksum[0][1:]) {
			return
		}

		segments := strings.Split(cksum[0], ",")
		size := len(segments)
		if size == 0 {
			return
		}

		if segments[0] == sGPGGA {
			// Time, position and fix related data
			time := toFloat64(segments[1])
			latitude := degToDec(segments[2], segments[3], 2, 8)
			longitude := degToDec(segments[4], segments[5], 3, 8)
			altitude := toFloat64(segments[9])

			gps.latitude = latitude
			gps.longitude = longitude
			gps.altitude = altitude
			gps.time = time
		}

		if segments[0] == sGPGSA {
			// Operating details
			satellites := segments[3:15]
			pdop := segments[15]
			hdop := segments[16]
			vdop := segments[17]

			gps.satellites = satellites
			gps.pdop = pdop
			gps.hdop = hdop
			gps.vdop = vdop
		}

		if segments[0] == sGPRMC {
			// GPS & Transit data
			time := toFloat64(segments[1])
			latitude := degToDec(segments[3], segments[4], 2, 8)
			longitude := degToDec(segments[5], segments[6], 3, 8)
			course := toFloat64(segments[8])
			speed := toSpeed(segments[7], 2)

			gps.longitude = longitude
			gps.latitude = latitude
			gps.course = course
			gps.speed = speed
			gps.time = time
		}

		if segments[0] == sGPVTG {
			// Track Mode Good and Ground Speed
			course := toFloat64(segments[1])
			speed := toFloat64(segments[5])

			gps.course = course
			gps.speed = speed
		}
	}
}
