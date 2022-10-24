package gps

import (
	"bytes"
	"testing"
	"time"
)

type serialPortFake struct {
	bytes.Buffer
}

func (s *serialPortFake) Close() error {
	return nil
}

// connectFake allows to simulate a serial connection
func (gps *GPS) connectFake() (err error) {
	var serialPortFake serialPortFake

	gps.serialPort = &serialPortFake
	gps.data = make(chan string)
	gps.isConnected = true

	return err
}

// writeFake allows to simulate a writes serial
func (gps *GPS) writeFake(lines []string) (err error) {
	for _, line := range lines {
		buffer := []byte(line)
		buffer = append(buffer, '\r', '\n')

		_, err = gps.serialPort.Write(buffer)
	}

	return
}

func TestGPS(t *testing.T) {
	t.Run("gps parsing success", func(t *testing.T) {
		lines := []string{
			"$GPGGA,202530.00,5109.0262,N,11401.8407,W,5,40,0.5,1097.36,M,-17.00,M,18,TSTR*61",
			"$GPGSA,A,3,04,27,09,16,08,03,07,21,,,,,1.62,1.02,1.25*02",
			"$GPRMC,203522.00,A,5109.0262308,N,11401.8407342,W,0.004,133.4,130522,0.0,E,D*2B",
			"$GPVTG,269.49,T,,M,0.02,N,0.04,K,D*3E",
		}

		want := Position{
			Latitude:  51.15043718,
			Longitude: -114.0306789,
			Altitude:  1097.36,
			Speed:     0.02,
			Course:    269.49,
		}

		mGPS := GPS{
			Port:     "/dev/ttyUSB0",
			BaudRate: 9600,
		}

		err := mGPS.connectFake()
		if err != nil {
			t.Errorf("expected error %v, want nil", err)
		}

		go mGPS.Reading()

		defer mGPS.Disconnect()

		mGPS.writeFake(lines)

		time.Sleep(250 * time.Millisecond)

		p, err := mGPS.FetchPosition()
		if err != nil {
			t.Errorf("expected %v, want nil", err)
		}

		if p.Latitude != want.Latitude {
			t.Errorf("expected %v, latitude want %v", p.Latitude, want.Latitude)
		}

		if p.Longitude != want.Longitude {
			t.Errorf("expected %v, longitude want %v", p.Longitude, want.Longitude)
		}

		if p.Altitude != want.Altitude {
			t.Errorf("expected %v, altitude want %v", p.Altitude, want.Altitude)
		}

		if p.Speed != want.Speed {
			t.Errorf("expected %v, speed want %v", p.Speed, want.Speed)
		}

		if p.Course != want.Course {
			t.Errorf("expected %v, course want %v", p.Course, want.Course)
		}
	})

	t.Run("gps parsing empty", func(t *testing.T) {
		lines := []string{
			"$GPGGA,161715.000,1858.3654,N,09334.2837,W,1,08,1.11,20.9,M,-8.8,M,,*5E",
			"$GPGSA,A,3,01,08,04,21,07,09,27,17,,,,,1.97,1.11,1.63*0B",
			"$GPRMC,161715.000,A,1858.3654,N,09334.2837,W,0.01,147.61,140421,,,A*75",
			"$GPVTG,147.61,T,,M,0.01,N,0.02,K,A*3B",
		}

		want := Position{
			Latitude:  0,
			Longitude: 0,
			Altitude:  0,
		}

		mGPS := GPS{
			Port:     "/dev/ttyUSB0",
			BaudRate: 9600,
		}

		err := mGPS.connectFake()
		if err != nil {
			t.Errorf("expected error %v, want nil", err)
		}

		go mGPS.Reading()

		defer mGPS.Disconnect()

		mGPS.writeFake(lines)

		time.Sleep(250 * time.Millisecond)

		p, err := mGPS.FetchPosition()
		if err != nil {
			t.Errorf("expected %v, want nil", err)
		}

		if p.Latitude != want.Latitude {
			t.Errorf("expected %v, latitude want %v", p.Latitude, want.Latitude)
		}

		if p.Longitude != want.Longitude {
			t.Errorf("expected %v, longitude want %v", p.Longitude, want.Longitude)
		}

		if p.Altitude != want.Altitude {
			t.Errorf("expected %v, altitude want %v", p.Altitude, want.Altitude)
		}
	})
}
