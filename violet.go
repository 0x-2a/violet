package main

import (
	"fmt"
	"log"
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

	info, err := lidar.DeviceInfo()
	log.Print(info)

	status, errcode, err := lidar.Health()
	if err != nil {
		log.Fatal(err)
	} else if status == "Warning" {
		log.Printf("Lidar status: %v Error Code: %v\n", status, errcode)
	} else if status == "Error" {
		log.Fatalf("Lidar status: %v Error Code: %v\n", status, errcode)
	}
	// lidar.StartMotor()

	// time.Sleep(10 * time.Second)

	count := 0
	iter := 1000

	lastTime := time.Now()
	err = lidar.StartScanFn(func(p *rplidar.RPLidarPoint) {
		count++
		if count % iter == 0 {
			now := time.Now()
			timeDiff := now.Sub(lastTime).Seconds()
			lastTime = now
			fmt.Printf("%s %d %f per second\n", DebugTimestamp(now.Unix()), count,  float32(iter)/float32(timeDiff))
		}

		// if p.Angle < 1 {
		// 	fmt.Printf("Quality: %v\tAngle: %.2f\tDistance: %.2f\n", p.Quality, p.Angle, p.Distance)
		// }
	})
	if err != nil {
		log.Fatal(err)
	}
	// for _, p := range scanResults {
	// 	fmt.Printf("Quality: %v\tAngle: %.2f\tDistance: %.2f\n", p.Quality, p.Angle, p.Distance)
	// }

	// lidar.StopMotor()
	lidar.Disconnect()
}

func DebugTimestamp(ts int64) string {
	return fmt.Sprintf("%d (%s)", ts, time.Unix(ts, 0).UTC().Format("2006-01-02 15:04:05 MST"))
}
