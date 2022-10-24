GPS 
========
A Go package to allow the Adafruit Ultimate GPS module to be read from the serial port.

	$GPGGA,161715.000,1858.3654,N,09334.2837,W,1,08,1.11,20.9,M,-8.8,M,,*5E
	$GPGSA,A,3,01,08,04,21,07,09,27,17,,,,,1.97,1.11,1.63*0B
	$GPRMC,161715.000,A,1858.3654,N,09334.2837,W,0.01,147.61,140421,,,A*75
	$GPVTG,147.61,T,,M,0.01,N,0.02,K,A*3B

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jamsMendez/gps"
)

func main() {
	mGps := gps.GPS{
		Port:     "COM5",
		BaudRate: 9600,
	}

	err := mGps.Connect()
	if err != nil {
		log.Fatal("GPS Connect ERROR: ", err)
	}

	go mGPs.Reading()

	defer mGps.Disconnect()

	for {
		p, err := mGps.FetchPosition()
		if err != nil {
			return
		}

		fmt.Println(p.Latitude, p.Longitude, p.Altitude, p.Speed, p.Course)

		time.Sleep(time.Second * 2)
	}
}
```

# Contact

Github: [https://github.com/jamsmendez](https://github.com/jamsmendez/)

Twitter: [https://twitter.com/jamsmendez](https://twitter.com/jamsmendez)