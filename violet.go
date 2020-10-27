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
	lidar.Reset()
	time.Sleep(10 * time.Second)
	status, errcode, err := lidar.Health()
	if err != nil {
		log.Fatal(err)
	} else if status == "Warning" {
		log.Printf("Lidar status: %v Error Code: %v\n", status, errcode)
	} else if status == "Error" {
		log.Fatalf("Lidar status: %v Error Code: %v\n", status, errcode)
	}
	// lidar.StartMotor()
	scanResults, err := lidar.StartScan(10)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range scanResults {
		fmt.Printf("Quality: %v\tAngle: %.2f\tDistance: %.2f\n", p.Quality, p.Angle, p.Distance)
	}

	// lidar.Disconnect()
}
