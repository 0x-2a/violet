package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/y3sh/violet/rplidar"
)

const (
	usbAddr = "/dev/tty.usbserial-0001"
)

func main() {
	lidar := rplidar.NewRPLidar(usbAddr, 115200)
	err := lidar.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := lidar.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()

	points := make(chan *rplidar.RPLidarPoint)
	var wg sync.WaitGroup
	wg.Add(2)

	go listen(lidar, points)
	go handlePoint(points)

	wg.Wait()
}

func DebugTimestamp(ts int64) string {
	return fmt.Sprintf("%d (%s)", ts, time.Unix(ts, 0).UTC().Format("2006-01-02 15:04:05 MST"))
}

func handlePoint(points chan *rplidar.RPLidarPoint){
	count := 0
	iter := 10
	var pointArr = [361]string{}

	for p := range points {
		count++
		_ = p

		present := "____"
		if int(p.Distance) < 2200 {
			present = "****"
		}

		pointArr[360-int(p.Angle)] = present

		if count%iter == 0 {
			fmt.Println(pointArr[310:])
		}

		// if p.Angle < 1 {`
		// 	fmt.Printf("Quality: %v\tAngle: %.2f\tDistance: %.2f\n", p.Quality, p.Angle, p.Distance)
		// }
	}
}

func listen(lidar *rplidar.RPLidar, points chan *rplidar.RPLidarPoint) {
	err := lidar.StartScanFn(func(p *rplidar.RPLidarPoint) bool {
		points <- p

		return false
	})
	if err != nil {
		log.Fatal(err)
	}
}
